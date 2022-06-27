package projections

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	domainbase "github.com/novabankapp/common.data/domain/base"
	"github.com/novabankapp/common.infrastructure/logger"
	"golang.org/x/sync/errgroup"
)

const (
	EsAll          = "$all"
	CassProjection = "(CassandraDB Projection)"
)

type CassandraProjection struct {
	Log logger.Logger
	Db  *esdb.Client
	Cfg *Subscriptions
}

func (o *CassandraProjection) runWorker(ctx context.Context, worker Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}
func (c *CassandraProjection) Subscribe(ctx context.Context, prefixes []string, poolSize int, worker Worker) error {
	//reflect.ValueOf(E).Type().Name()
	c.Log.Infof("(starting order subscription) prefixes: {%+v}", prefixes)

	err := c.Db.CreatePersistentSubscriptionAll(ctx, c.Cfg.CassandraProjectionGroupName, esdb.PersistentAllSubscriptionOptions{
		Filter: &esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: prefixes},
	})
	if err != nil {
		if subscriptionError, ok := err.(*esdb.PersistentSubscriptionError); !ok || ok && (subscriptionError.Code != 6) {
			c.Log.Errorf("(CreatePersistentSubscriptionAll) err: {%v}", subscriptionError.Error())
		}
	}

	stream, err := c.Db.ConnectToPersistentSubscription(
		ctx,
		EsAll,
		c.Cfg.CassandraProjectionGroupName,
		esdb.ConnectToPersistentSubscriptionOptions{},
	)
	if err != nil {
		return err
	}
	defer stream.Close()

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i <= poolSize; i++ {
		g.Go(c.runWorker(ctx, worker, stream, i))
	}
	return g.Wait()
}

/*func (c *CassandraProjection[E]) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {

	for {
		event := stream.Recv()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if event.SubscriptionDropped != nil {
			c.Log.Errorf("(SubscriptionDropped) err: {%v}", event.SubscriptionDropped.Error)
			return errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
		}

		if event.EventAppeared != nil {
			c.Log.ProjectionEvent(CassProjection, c.Cfg.CassandraProjectionGroupName, event.EventAppeared, workerID)

			err := c.When(ctx, es.NewEventFromRecorded(event.EventAppeared.Event))
			if err != nil {
				c.Log.Errorf("(CassProjection.when) err: {%v}", err)

				if err := stream.Nack(err.Error(), esdb.Nack_Retry, event.EventAppeared); err != nil {
					c.Log.Errorf("(stream.Nack) err: {%v}", err)
					return errors.Wrap(err, "stream.Nack")
				}
			}

			err = stream.Ack(event.EventAppeared)
			if err != nil {
				c.Log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}
			c.Log.Infof("(ACK) event commit: {%v}", *event.EventAppeared.Commit)
		}
	}
}

func (c *CassandraProjection[E]) When(ctx context.Context, evt es.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "CassandraProjection.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	switch evt.GetEventType() {

	case v1.WalletCreated:
		return c.onWalletCreated(ctx, evt)
	case v1.WalletCredited:
		return c.onWalletCredited(ctx, evt)
	case v1.WalletDebited:
		return c.onWalletDebited(ctx, evt)
	case v1.WalletCreditReserved:
		return c.onWalletCreditReserved(ctx, evt)
	case v1.WalletBlacklisted:
		return c.onWalletBlacklisted(ctx, evt)
	case v1.WalletLocked:
		return c.onWalletLocked(ctx, evt)
	case v1.WalletCreditReleased:
		return c.onWalletCreditReleased(ctx, evt)

	default:
		c.Log.Warnf("(CassandraProjection) [When unknown EventType] eventType: {%s}", evt.EventType)
		return es.ErrInvalidEventType
	}
}*/

func NewCassandraProjection[E domainbase.NoSqlEntity]() *CassandraProjection[E] {
	return &CassandraProjection[E]{}
}
