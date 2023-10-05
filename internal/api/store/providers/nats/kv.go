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

type (
	kvOptions struct{}
	KVOption  func(*kvOptions)
)

func NewKV(addr string, options ...kvOptions) (*NatsStore, error) {
	nc, err := nats.Connect(addr)
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

// WithTenantKey sets the key for the tenant.
func WithTenantKey(key string) KVOption {
	return func(o *kvOptions) {
		TennantKey = key
	}
}

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
