package store

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/novabankapp/common.data/eventstore"
)

func NewEventStoreDB(cfg eventstore.EventStoreConfig) (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	return esdb.NewClient(settings)
}
