package test

import (
	"context"
	"gazur/pkg/subscriptions"
	"gazur/pkg/subscriptions/ignores"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionFetch(t *testing.T) {
	pages := 3
	itemsPerPage := 5

	client := NewMockedSubPager(pages, itemsPerPage)
	pager := subscriptions.GetPager(&client)

	for i := 0; i < pages; i++ {
		assert.True(t, pager.More())

		nextResult, err := pager.NextPage(nil)
		assert.NoError(t, err)

		assert.Equal(t, client.subs[i], nextResult.Value)
	}
	assert.False(t, pager.More())
}

func PassSubsToChannel(in chan<- *armsubscriptions.Subscription, client *MockedSubPager) {
	defer close(in)
	for _, batchSubs := range client.subs {
		for _, s := range batchSubs {
			in <- s
		}
	}
}

func TestFilteringOutSubProperly(t *testing.T) {
	pages := 3
	itemsPerPage := 5
	ctx := context.Background()

	client := NewMockedSubPager(pages, itemsPerPage)
	toFilterOut := client.subs[0][0]

	filterConfig := ignores.Config{SubscriptionsToIgnore: map[string][]string{*toFilterOut.SubscriptionID: nil}}

	filters := subscriptions.Filter{}
	filters.AddFilter(func(s *armsubscriptions.Subscription) bool {
		_, exists := filterConfig.SubscriptionsToIgnore[*s.SubscriptionID]
		return exists
	})

	in := make(chan *armsubscriptions.Subscription)
	go PassSubsToChannel(in, &client)

	filtered, filterErr := filters.FilterSubs(ctx, in)
	for value := range filtered {
		// ensure subscription that was set to ignore
		// is not present in the solutions list
		assert.NotEqual(t, value, toFilterOut)
	}

	// ensure there was no error while filtering
	assert.Empty(t, filterErr)
}

func TestFilteringOutRestOfSubsPresentInResult(t *testing.T) {
	pages := 1
	itemsPerPage := 5
	ctx := context.Background()

	client := NewMockedSubPager(pages, itemsPerPage)
	toFilterOut := client.subs[0][0]

	filterConfig := ignores.Config{SubscriptionsToIgnore: map[string][]string{*toFilterOut.SubscriptionID: nil}}

	filters := subscriptions.Filter{}
	filters.AddFilter(func(s *armsubscriptions.Subscription) bool {
		_, exists := filterConfig.SubscriptionsToIgnore[*s.SubscriptionID]
		return exists
	})

	in := make(chan *armsubscriptions.Subscription)
	go PassSubsToChannel(in, &client)

	filtered, filterErr := filters.FilterSubs(ctx, in)

	i := 0
	for value := range filtered {
		assert.Contains(t, client.subs[0][1:], value)
		i++
	}

	assert.Equal(t, i, 4)

	// ensure there was no error while filtering
	assert.Len(t, filterErr, 0)
}

func TestAllSubsIgnoredResultsInEmptyList(t *testing.T) {
	pages := 1
	itemsPerPage := 5
	ctx := context.Background()

	client := NewMockedSubPager(pages, itemsPerPage)
	toFilterOut := make(map[string][]string)
	for _, sub := range client.subs[0] {
		toFilterOut[*sub.SubscriptionID] = nil
	}

	filterConfig := ignores.Config{SubscriptionsToIgnore: toFilterOut}

	filters := subscriptions.Filter{}
	filters.AddFilter(func(s *armsubscriptions.Subscription) bool {
		_, exists := filterConfig.SubscriptionsToIgnore[*s.SubscriptionID]
		return exists
	})

	in := make(chan *armsubscriptions.Subscription)
	go PassSubsToChannel(in, &client)

	filtered, filterErr := filters.FilterSubs(ctx, in)

	select {
	case valShouldBeNil := <-filtered:
		assert.Nil(t, valShouldBeNil)
	case <-time.After(10 * time.Second):
		t.Fatal("Unable to receive from filtered channel ")
	}

	// ensure there was no error while filtering
	assert.Empty(t, filterErr)

}
