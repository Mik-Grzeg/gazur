package server

import (
	"gazur/pkg/common"
	"gazur/pkg/subscriptions/ignores"
	"gazur/pkg/subscriptions"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

type Gazur struct {
	identity *common.Identity
	configFile *string
}

func New(identity *common.Identity, configFile *string) Gazur {
	return Gazur{
		identity,
		configFile,
	}
}

func (g *Gazur) Run() error {

	ignoresLoader := ignores.LoaderFromFile{Path: g.configFile}
	filterConfig, _ := ignores.GetSubIgnores(&ignoresLoader)

	filters := Filter{}
	filters.AddFilter(func(s *armsubscriptions.Subscription) bool {
		_, exists := filterConfig.SubscriptionsToIgnore[*s.SubscriptionID]
		return exists
	})

	if err := subscriptions.PipelineStart(g.identity, filters); err != nil {
		return err
	}

	return nil
}
