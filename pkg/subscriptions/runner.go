package subscriptions

import (
	"context"
	"encoding/json"
	"gazur/pkg/common"
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

func PipelineStart(identity *common.Identity, filters *Filter) error {
	ctx := context.Background()

	client, err := armsubscriptions.NewClient(identity, nil)
	if err != nil {
		log.Fatalf("Unable to get subscriptions client: %v", err)
	}
	pager := GetPager(client)

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

func check(identity *common.Identity, filters *Filter) error {
	ctx := context.Background()

	client, err := armsubscriptions.NewClient(identity, nil)
	if err != nil {
		log.Fatalf("Unable to get subscriptions client: %v", err)
	}
	pager := GetPager(client)

	fetching, _ := ListSubsFromPager(ctx, pager)
	taskManager := NewTaskManager().AddSource(&filters.FilterSub, fetching)

	outC := taskManager.Run()

	for elem := range outC {
		log.Printf("Elem: %v", elem)
	}

}
