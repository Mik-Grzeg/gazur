package subscriptions

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"log"
)

type Filter struct {
	check func(subscription *armsubscriptions.Subscription) bool
	next  *Filter
}

func (f *Filter) AddFilter(check func(subscription *armsubscriptions.Subscription) bool) *Filter {
	if f.check == nil {
		f.check = check
		return f
	} else {
		newFilter := Filter{
			check: check,
			next:  nil,
		}
		f.next = &newFilter
		return &newFilter
	}
}

func (f *Filter) run(subscription *armsubscriptions.Subscription) bool {
	filter := f
	for filter != nil {
		if filter.check(subscription) {
			return false
		}
		filter = filter.next
	}
	return true
}

func (f *Filter) FilterSubs(c context.Context, in <-chan *armsubscriptions.Subscription) (<-chan *armsubscriptions.Subscription, <-chan error) {
	out := make(chan *armsubscriptions.Subscription)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for {
			select {
			case sub, ok := <-in:
				if !ok {
					log.Printf("[Filterer] Got nil from input channel")
					return
				}

				if f.run(sub) {
					out <- sub
				}
			case <-c.Done():
				log.Printf("Filterer closing channel")
				return
			}
		}
	}()

	return out, errc
}
