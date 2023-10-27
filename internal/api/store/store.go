package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/trace"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
)

const (
	TennantKey = "tenants"
)

var (
	// Add custom errors here.
	ErrTenantAlreadyExists = errors.New("tenant already exists")
)

// Package store provides an interface and implementation for storing and retrieving data for tenants.
// The Store interface defines methods for creating tenants, getting tenants, listing tenants, writing records, and starting the indexer.
// The store struct implements the Store interface and uses StreamService, KVService, and IndexService to store and retrieve data.
type (
	Store interface {
		// KVService
		CreateTenant(ctx context.Context, tenant *Tenant) (*Tenant, error)
		GetTenant(ctx context.Context, tenantId string) (*Tenant, error)
		ListTenants(ctx context.Context) ([]string, error)
		WatchTenants(ctx context.Context) (<-chan KVEntry, <-chan struct{})
		// StreamService
		Write(ctx context.Context, tenantId string, record *api.Record) (uint64, error)
		StartIndexer(ctx context.Context) error
	}

	store struct {
		stream StreamService
		kv     KVService
		index  IndexService
		tracer trace.Tracer
		log    logr.Logger
	}

	StreamService interface {
		Write(ctx context.Context, tenant string, record *api.Record) (uint64, error)
	}

	KVService interface {
		CreateTenant(ctx context.Context, tenantId string, tenant []byte) ([]byte, error)
		GetTenant(ctx context.Context, tenantId string) ([]byte, error)
		ListTenants(ctx context.Context) ([]string, error)
		Watch(ctx context.Context) (<-chan KVEntry, <-chan struct{})
	}

	IndexService interface {
		Write(ctx context.Context, seq uint64, record *api.Record) error
		Start(ctx context.Context) error
	}

	KVEntry interface {
		Bucket() string
		Key() string
		Value() []byte
		Revision() uint64
		Created() time.Time
		Operation() KVWatchOp
	}
)

func NewStore(streamService StreamService, kvService KVService, indexService IndexService) Store {
	return &store{
		stream: streamService,
		kv:     kvService,
		index:  indexService,
	}
}

func (s *store) StartIndexer(ctx context.Context) error {
	s.log.Info("Starting indexer")
	return s.index.Start(ctx)
}

func (s *store) CreateTenant(ctx context.Context, tenant *Tenant) (*Tenant, error) {
	t, err := tenant.ToJSON()
	if err != nil {
		return nil, err
	}

	_, err = s.kv.CreateTenant(ctx, tenant.Id, t)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func (s *store) GetTenant(ctx context.Context, tenantId string) (*Tenant, error) {
	_, err := s.kv.GetTenant(ctx, tenantId)
	if err != nil {
		return nil, err
	}
	var tenant Tenant

	t, err := s.kv.GetTenant(ctx, tenantId)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(t, &tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (s *store) ListTenants(ctx context.Context) ([]string, error) {
	return s.kv.ListTenants(ctx)
}

// Write writes a record to the stream.
func (s *store) Write(ctx context.Context, tenant string, record *api.Record) (uint64, error) {
	return s.stream.Write(ctx, tenant, record)
}

// WatchTenant watches the tenant for changes.
func (s *store) WatchTenants(ctx context.Context) (<-chan KVEntry, <-chan struct{}) {
	return s.kv.Watch(ctx)
}
