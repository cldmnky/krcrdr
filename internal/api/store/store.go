package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"errors"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	TennantKey = "tenants"
)

var (
	// Add custom errors here.
	ErrTenantAlreadyExists = errors.New("tenant already exists")
)

type (
	Store interface {
		CreateTenant(ctx context.Context, tenantId string, tenant *Tenant) (*Tenant, error)
		GetTenant(ctx context.Context, tenantId string) (*Tenant, error)
		WriteStream(ctx context.Context, tenantId string, record *api.Record) error
	}

	store struct {
		stream StreamService
		kv     KVService
	}

	StreamService interface {
		Write(ctx context.Context, tenant string, record *api.Record) error
	}

	KVService interface {
		CreateTenant(ctx context.Context, tenantId string, tenant *Tenant) (*Tenant, error)
		GetTenant(ctx context.Context, tenantId string) (*Tenant, error)
	}

	NatsStore struct {
		nc *nats.Conn
		js jetstream.JetStream
		kv jetstream.KeyValue
	}
)

func NewStore(streamService StreamService, kvService KVService) Store {
	return &store{
		stream: streamService,
		kv:     kvService,
	}
}

func (s *store) CreateTenant(ctx context.Context, tenantId string, tenant *Tenant) (*Tenant, error) {
	return s.kv.CreateTenant(ctx, tenantId, tenant)
}

func (s *store) GetTenant(ctx context.Context, tenantId string) (*Tenant, error) {
	return s.kv.GetTenant(ctx, tenantId)
}

// WriteStream writes a record to the stream.
func (s *store) WriteStream(ctx context.Context, tenant string, record *api.Record) error {
	return s.stream.Write(ctx, tenant, record)
}

// StreamService is the interface for a stream store.
func NewNatsStream(addr string) (*NatsStore, error) {
	nc, err := nats.Connect(addr)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(nc, jetstream.WithPublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}
	return &NatsStore{
		nc: nc,
		js: js,
	}, nil
}

func (s *NatsStore) Write(ctx context.Context, tenantId string, record *api.Record) error {
	// marshal the record to JSON.
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}
	subject := s.createSubjectFromRecord(record)
	_, err = s.js.Publish(ctx, fmt.Sprintf("%s.%s", strings.ToUpper(tenantId), subject), recordJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *NatsStore) createSubjectFromRecord(record *api.Record) string {
	if record.Namespace == "" {
		return fmt.Sprintf("%s.cluster.%s.%s", record.Cluster, record.Kind.Kind, record.Name)
	}
	return fmt.Sprintf("%s.namespace.%s.%s.%s", record.Cluster, record.Namespace, record.Kind.Kind, record.Name)
}

// KVService is the interface for a key-value store.
func NewNatsKV(addr string) (*NatsStore, error) {
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

type Tenant struct {
	Name    string   `json:"name"`
	ApiKeys []string `json:"apiKeys"`
}

func (t *Tenant) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (s *NatsStore) CreateTenant(ctx context.Context, tenantId string, tenant *Tenant) (*Tenant, error) {
	tenantJSON, err := tenant.ToJSON()
	if err != nil {
		return nil, err
	}

	// Try to get the tenant first.
	_, err = s.kv.Get(ctx, tenantId)
	if err == nil {
		// If the tenant already exists, return an error.
		return nil, ErrTenantAlreadyExists
	}
	_, err = s.kv.Put(ctx, tenantId, tenantJSON)
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

func (s *NatsStore) GetTenant(ctx context.Context, tenantId string) (*Tenant, error) {
	v, err := s.kv.Get(ctx, tenantId)
	if err != nil {
		return nil, err
	}
	// Unmarshal the JSON into a Tenant object.
	var tenant Tenant
	err = json.Unmarshal(v.Value(), &tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}
