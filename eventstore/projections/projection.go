package projections

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	es "github.com/novabankapp/common.data/eventstore"
)

type Worker func(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error

type Projection interface {
	Subscribe(ctx context.Context, prefixes []string, poolSize int, worker Worker) error
	ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error
	When(ctx context.Context, evt es.Event) error
}
