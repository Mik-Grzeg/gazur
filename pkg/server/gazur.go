package server

import (
	"gazur/pkg/common"
	"gazur/pkg/subscriptions"
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
	if err := subscriptions.PipelineStart(g.identity, g.configFile); err != nil {
		return err
	}

	return nil
}
