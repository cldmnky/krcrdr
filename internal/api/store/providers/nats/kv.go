package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
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

// WatchTenants will watch for changes to the tenants kv
func (s *NatsStore) WatchTenants(ctx context.Context) error {
	w, err := s.kv.WatchAll(ctx)
	if err != nil {
		return err
	}
	for {
		select {
		case kve := <-w.Updates():
			logger.Info("got update", "key", kve.Key, "value", kve.Value())
		case <-ctx.Done():
			return nil
		}
	}
}
