package subscriptions

import (
	"context"
	"encoding/json"
	"gazur/pkg/common"
	"gazur/pkg/subscriptions/ignores"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"log"
)

func merge(cs ...<-chan error) <-chan error {
	out := make(chan error)

	for _, c := range cs {
		go func() {
			for v := range c {
				out <- v
			}
		}()
	}
	return out
}

func PipelineStart(identity *common.Identity, configFile *string) error {
	ctx := context.Background()

	client, err := armsubscriptions.NewClient(identity, nil)
	if err != nil {
		log.Fatalf("Unable to get subscriptions client: %v", err)
	}
	pager := GetPager(client)

	// settings filter

	ignoresLoader := ignores.LoaderFromFile{Path: configFile}
	filterConfig, _ := ignores.GetSubIgnores(&ignoresLoader)

	filters := Filter{}
	filters.AddFilter(func(s *armsubscriptions.Subscription) bool {
		_, exists := filterConfig.SubscriptionsToIgnore[*s.SubscriptionID]
		return exists
	})

	// Getting data
	fetching, fetchErr := ListSubsFromPager(ctx, pager)
	filtering, filterErr := filters.FilterSubs(ctx, fetching)

	for {
		select {
		case out, ok := <-filtering:
			if !ok {
				return nil
			}

			ser, _ := json.MarshalIndent(out, "", "    ")
			log.Printf("Got result from pipeline: %s", string(ser))
		case err := <-merge(fetchErr, filterErr):
			log.Printf("Got error in the pipeline: %v", err)
			return err
		}
	}
}
