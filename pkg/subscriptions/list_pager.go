package subscriptions

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"log"
)

type SubClient struct {
	*armsubscriptions.Client
	Ctx context.Context
}

func ListSubsFromPager(c context.Context, pager *runtime.Pager[armsubscriptions.ClientListResponse]) (<-chan *armsubscriptions.Subscription, <-chan error) {
	out := make(chan *armsubscriptions.Subscription)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for pager.More() {
			nextResult, err := pager.NextPage(c)
			if err != nil {
				log.Printf("Source got error: %v", err)
				errc <- err
			}

			for _, v := range nextResult.Value {
				select {
				case out <- v:
					log.Printf("Source sending to channel: %v", v)
				case <-c.Done():
					log.Printf("Source closing channel")
					return
				}
			}
		}
	}()

	return out, errc
}

type SubPager interface {
	NewListPager(options *armsubscriptions.ClientListOptions) *runtime.Pager[armsubscriptions.ClientListResponse]
}

func GetPager(c SubPager) *runtime.Pager[armsubscriptions.ClientListResponse] {
	return c.NewListPager(nil)
}
