package store

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
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
		CreateTenant(ctx context.Context, tenant *Tenant) (*Tenant, error)
		GetTenant(ctx context.Context, tenantId string) (*Tenant, error)
		ListTenants(ctx context.Context) ([]string, error)
		Write(ctx context.Context, tenantId string, record *api.Record) (uint64, error)
		StartIndexer(ctx context.Context) error
	}

	store struct {
		stream StreamService
		kv     KVService
		index  IndexService
	}

	StreamService interface {
		Write(ctx context.Context, tenant string, record *api.Record) (uint64, error)
	}

	KVService interface {
		CreateTenant(ctx context.Context, tenantId string, tenant []byte) ([]byte, error)
		GetTenant(ctx context.Context, tenantId string) ([]byte, error)
		ListTenants(ctx context.Context) ([]string, error)
	}

	IndexService interface {
		Write(ctx context.Context, seq uint64, record *api.Record) error
		Start(ctx context.Context) error
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
