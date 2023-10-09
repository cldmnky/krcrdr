package nats

import (
	"context"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var (
	TennantKey = "tenants"
)

func NewKV(url string, options ...nats.Option) (*NatsStore, error) {
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
		nc: nc,
		js: js,
		kv: kv,
	}, nil
}

// CreateTenant creates a new tenant in the NatsStore with the given tenantId and tenant data.
// If the tenant already exists, it returns an error.
// It also creates a stream for the tenant with the given tenantId and subjects.
// If the stream already exists, it ignores the error.
// It returns the tenant data as []byte and an error if any occurred.
func (s *NatsStore) CreateTenant(ctx context.Context, tenantId string, tenant []byte) ([]byte, error) {
	// Try to get the tenant first.
	_, err := s.kv.Get(ctx, tenantId)
	if err == nil {
		// If the tenant already exists, return an error.
		return nil, fmt.Errorf("tenant already exists")
	}
	_, err = s.kv.Put(ctx, tenantId, tenant)
	if err != nil {
		return nil, err
	}
	// Create the stream for the tenant.
	_, err = s.js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     tenantId,
		Subjects: []string{fmt.Sprintf("%s.>", strings.ToUpper(tenantId))},
	})
	if err != nil {
		// If the stream already exists, ignore the error.
		if err != jetstream.ErrStreamNameAlreadyInUse {
			return nil, err
		}
	}

	return s.GetTenant(ctx, tenantId)
}

func (s *NatsStore) GetTenant(ctx context.Context, tenantId string) ([]byte, error) {
	v, err := s.kv.Get(ctx, tenantId)
	if err != nil {
		return nil, err
	}
	return v.Value(), nil
}

func (s *NatsStore) ListTenants(ctx context.Context) ([]string, error) {
	return s.kv.Keys(ctx)
}
