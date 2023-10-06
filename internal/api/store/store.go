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
		CreateTenant(ctx context.Context, tenantId string, tenant []byte) ([]byte, error)
		GetTenant(ctx context.Context, tenantId string) ([]byte, error)
		ListTenants(ctx context.Context) ([]string, error)
	}
)

func NewStore(streamService StreamService, kvService KVService) Store {
	return &store{
		stream: streamService,
		kv:     kvService,
	}
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

// WriteStream writes a record to the stream.
func (s *store) WriteStream(ctx context.Context, tenant string, record *api.Record) error {
	return s.stream.Write(ctx, tenant, record)
}
