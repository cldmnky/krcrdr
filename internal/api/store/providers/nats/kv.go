package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cldmnky/krcrdr/internal/api/store"
)

var (
	TennantKey = "tenants"
	logger     = logf.Log.WithName("apiHandler")
)

func NewKV(url string, tracer trace.Tracer, options ...nats.Option) (*NatsStore, error) {
	nc, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(nc, jetstream.WithPublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}
	kv, err := js.CreateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket: TennantKey,
	})
	if err != nil {
		return nil, err
	}
	return &NatsStore{
		nc:     nc,
		js:     js,
		kv:     kv,
		tracer: tracer,
	}, nil
}

// CreateTenant creates a new tenant in the NatsStore with the given tenantId and tenant data.
// If the tenant already exists, it returns an error.
// It also creates a stream for the tenant with the given tenantId and subjects.
// If the stream already exists, it ignores the error.
// It returns the tenant data as []byte and an error if any occurred.
func (s *NatsStore) CreateTenant(ctx context.Context, tenantId string, tenant []byte) ([]byte, error) {
	// trace
	ctx, span := s.tracer.Start(ctx, "nats/CreateTenant")
	defer span.End()

	// Try to get the tenant first.
	span.AddEvent("Get tenant")
	_, err := s.kv.Get(ctx, tenantId)
	if err == nil {
		// If the tenant already exists, return an error.
		span.RecordError(fmt.Errorf("tenant already exists"))
		return nil, fmt.Errorf("tenant already exists")
	}
	span.AddEvent("Put tenant")
	_, err = s.kv.Put(ctx, tenantId, tenant)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	// Create the stream for the tenant.
	span.AddEvent("Create stream", trace.WithAttributes(
		attribute.String("tenant", tenantId),
	))
	_, err = s.GetOrCreateStream(ctx, tenantId)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return s.GetTenant(ctx, tenantId)
}

func (s *NatsStore) GetTenant(ctx context.Context, tenantId string) ([]byte, error) {
	ctx, span := s.tracer.Start(ctx, "nats/GetTenant")
	defer span.End()

	span.AddEvent("Get tenant")
	v, err := s.kv.Get(ctx, tenantId)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return v.Value(), nil
}

func (s *NatsStore) ListTenants(ctx context.Context) ([]string, error) {
	ctx, span := s.tracer.Start(ctx, "nats/ListTenants")
	defer span.End()
	tenants, err := s.kv.Keys(ctx)
	if err != nil {
		if err != jetstream.ErrNoKeysFound {
			span.RecordError(err)
			return nil, err
		}
		return nil, nil
	}
	span.SetAttributes(attribute.Int("tenants", len(tenants)))
	return tenants, nil
}

// Watch returns a channel of KVEntry and a done channel.
func (s *NatsStore) Watch(ctx context.Context) (<-chan store.KVEntry, <-chan struct{}) {
	var (
		err     error
		updates = make(chan store.KVEntry)
		done    = make(chan struct{})
	)
	w, err := s.newWatcher(ctx)
	if err != nil {
		logger.Error(err, "Error setting up watcher")
		close(updates)
		close(done)
		return nil, nil
	}
	// Run the watcher in a goroutine.
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("Context done, closing watcher")
				w.Stop()
				close(updates)
				close(done)
				return
			case <-done:
				logger.Info("Done channel closed, closing watcher")
				w.Stop()
				close(updates)
				return
			case entry, ok := <-w.Updates():
				if !ok {
					logger.Info("Updates channel closed, closing watcher")
					return
				}
				if entry == nil {
					logger.Info("Loaded initial state")
					continue
				}
				updates <- kvEntry{
					bucket:    entry.Bucket(),
					key:       entry.Key(),
					value:     entry.Value(),
					revision:  entry.Revision(),
					created:   entry.Created(),
					operation: store.KVWatchOp(entry.Operation()),
				}

			}
		}
	}()

	return updates, done
}

// newWatcher sets up a nats kv WatchAll and return a KeyWatcher.
func (s *NatsStore) newWatcher(ctx context.Context) (jetstream.KeyWatcher, error) {
	watcher, err := s.kv.WatchAll(nats.Context(ctx))
	for err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			logger.Error(err, "Error setting up watcher, retrying in 100ms")
			watcher, err = s.kv.WatchAll(nats.Context(ctx))
			time.Sleep(100 * time.Millisecond)
		}
	}

	return watcher, nil
}
