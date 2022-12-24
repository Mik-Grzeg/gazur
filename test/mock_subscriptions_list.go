package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/brianvoe/gofakeit/v6"
	"io"
	"log"
	"os"
)

const ValsPerPage = 3

type MockedSubPager struct {
	nextPage int
	subs     [][]*armsubscriptions.Subscription
}

func NewMockedSubPagerFromFixtures(pages int, filePath string) MockedSubPager {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteResult, _ := io.ReadAll(file)
	var data []*armsubscriptions.Subscription
	json.Unmarshal(byteResult, &data)

	if pages > len(data) {
		log.Fatalf("Length of fixture subscription list is shorter than pages. Unable to create pager")
	}

	dataSlices := make([][]*armsubscriptions.Subscription, pages)
	quotient := len(data) / pages
	remainder := len(data) % pages
	log.Printf("Len of fixtures: %d, Pages: %d\nQuotient: %d, Remainder: %d", len(data), pages, quotient, remainder)

	for i := 0; i < pages-1; i++ {
		log.Printf("data[%d*%d : %d+%d]", i, quotient, i, quotient)
		dataSlices[i] = data[i*quotient : i+quotient]
	}
	log.Printf("last remainder: data[%d-%d]=%v", len(data), remainder, data[len(data)-remainder:])
	dataSlices[pages-1] = data[len(data)-1-remainder:]

	serialized, err := json.Marshal(dataSlices)
	log.Printf("Num of pages: %d\nPages: %v", len(dataSlices), string(serialized))

	return MockedSubPager{
		nextPage: 0,
		subs:     dataSlices,
	}
}

func NewMockedSubPager(pages int, itemsPerPage int) MockedSubPager {
	data := make([][]*armsubscriptions.Subscription, pages)

	for i := 0; i < pages; i++ {
		data[i] = make([]*armsubscriptions.Subscription, itemsPerPage)
		for j := 0; j < itemsPerPage; j++ {
			tenantId := gofakeit.UUID()
			subId := gofakeit.UUID()
			color := gofakeit.Color()
			name := gofakeit.AppName()
			state := armsubscriptions.SubscriptionStateEnabled

			id := fmt.Sprintf("/subscriptions/%s", subId)
			data[i][j] = &armsubscriptions.Subscription{
				ID:                   &id,
				SubscriptionID:       &subId,
				TenantID:             &tenantId,
				SubscriptionPolicies: nil,
				ManagedByTenants:     nil,
				AuthorizationSource:  nil,
				Tags: map[string]*string{
					"color": &color,
				},
				DisplayName: &name,
				State:       &state,
			}
		}
	}

	return MockedSubPager{
		nextPage: 0,
		subs:     data,
	}
}

func (m *MockedSubPager) More(c armsubscriptions.ClientListResponse) bool {
	return m.nextPage < len(m.subs)
}

func (m *MockedSubPager) Fetcher(ctx context.Context, c *armsubscriptions.ClientListResponse) (armsubscriptions.ClientListResponse, error) {
	data := armsubscriptions.SubscriptionListResult{
		NextLink: nil,
		Value:    m.subs[m.nextPage],
	}
	m.nextPage += 1
	return armsubscriptions.ClientListResponse{SubscriptionListResult: data}, nil
}

func (m *MockedSubPager) NewListPager(options *armsubscriptions.ClientListOptions) *runtime.Pager[armsubscriptions.ClientListResponse] {
	handler := runtime.PagingHandler[armsubscriptions.ClientListResponse]{
		More:    m.More,
		Fetcher: m.Fetcher,
	}
	return runtime.NewPager(handler)
}
