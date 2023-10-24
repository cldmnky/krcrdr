package frostdb

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/run"
	"go.opentelemetry.io/otel/trace"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	natsprovider "github.com/cldmnky/krcrdr/internal/api/store/providers/nats"
)

var logger = logf.Log.WithName("apiHandler")

type FrostDBIndex struct {
	runGroup  *run.Group
	natsStore *natsprovider.NatsStore
	tracer    trace.Tracer
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
// This is a no-op for now.
func (i *FrostDBIndex) Start(ctx context.Context) error {
	var stream jetstream.Stream
	// Get all tenants
	tenants, err := i.natsStore.ListTenants(ctx)
	if err != nil {
		return err
	}
	// For each tenant get the stream
	for _, tenant := range tenants {
		// Get the stream
		stream, err = i.natsStore.GetOrCreateStream(ctx, tenant)
		if err != nil {
			return err
		}
		info, _ := stream.Info(ctx)
		logger.Info("got stream", "stream", info)
		// Run a consumer for the stream
		err = i.runConsumer(ctx, stream)
		if err != nil {
			return err
		}
	}
	if err := i.natsStore.WatchTenants(ctx); err != nil {
		return err
	}
	return nil
}

// runConsumer will run a consumer for the given stream.
func (i *FrostDBIndex) runConsumer(ctx context.Context, stream jetstream.Stream) error {
	// Create a consumer for the stream.
	streamName := stream.CachedInfo().Config.Name
	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       stream.CachedInfo().Config.Name,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: fmt.Sprintf("%s.*", streamName),
	})
	if err != nil {
		return err
	}
	consCtx, err := cons.Consume(func(msg jetstream.Msg) {
		logger.Info("got message", "subject", msg.Subject, "data", string(msg.Data()))
	})
	if err != nil {
		return err
	}
	defer consCtx.Stop()
	return fmt.Errorf("consumer stopped")
}
