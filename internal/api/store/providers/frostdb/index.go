package frostdb

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	natsprovider "github.com/cldmnky/krcrdr/internal/api/store/providers/nats"
)

var logger = logf.Log.WithName("apiHandler")

type FrostDBIndex struct {
	//runGroup  *run.Group
	natsStore *natsprovider.NatsStore
	tracer    trace.Tracer
	tenants   []string
}

func NewIndex(natsStore *natsprovider.NatsStore, tracer trace.Tracer) (*FrostDBIndex, error) {
	return &FrostDBIndex{
		natsStore: natsStore,
		tracer:    tracer,
	}, nil
}

// Implement Write
func (i *FrostDBIndex) Write(ctx context.Context, seq uint64, record *api.Record) error {
	return nil
}

// Implement IndexService Start
func (i *FrostDBIndex) Start(ctx context.Context) error {
	// Get all tenants
	tenants, err := i.natsStore.ListTenants(ctx)
	if err != nil {
		return err
	}
	i.tenants = tenants
	return nil
}
