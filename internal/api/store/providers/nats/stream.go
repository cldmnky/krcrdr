package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
)

// StreamService is the interface for a stream store.
func NewStream(url string, tracer trace.Tracer, options ...nats.Option) (*NatsStore, error) {
	nc, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(nc, jetstream.WithPublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}
	return &NatsStore{
		nc:     nc,
		js:     js,
		tracer: tracer,
	}, nil
}

func (s *NatsStore) Write(ctx context.Context, tenantId string, record *api.Record) (uint64, error) {
	ctx, span := s.tracer.Start(ctx, "nats/Write")
	defer span.End()
	// marshal the record to JSON.
	recordJSON, err := json.Marshal(record)
	if err != nil {
		span.RecordError(err)
		return 0, err
	}
	subject := CreateSubjectFromRecord(strings.ToUpper(tenantId), record)
	span.SetAttributes(attribute.String("subject", subject))
	ack, err := s.js.Publish(ctx, subject, recordJSON)
	if err != nil {
		span.RecordError(err)
		return 0, err
	}
	span.SetAttributes(attribute.Int64("sequence", int64(ack.Sequence)))
	return ack.Sequence, nil
}

// Consume will consume messages from the stream.
func (s *NatsStore) Consume(ctx context.Context, stream string) error {
	_, err := s.js.CreateOrUpdateConsumer(ctx, stream, jetstream.ConsumerConfig{
		Durable:   stream,
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}
	return nil
}

// GetOrCreateStream will get or create a stream with the given name.
func (s *NatsStore) GetOrCreateStream(ctx context.Context, stream string) (jetstream.Stream, error) {
	// Try to get the stream first.
	str, err := s.js.Stream(ctx, stream)
	if err == nil {
		// If the stream already exists, return it.
		return str, nil
	}
	// If the stream does not exist, create it.
	str, err = s.js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     stream,
		Subjects: []string{fmt.Sprintf("%s.>", strings.ToUpper(stream))},
	})
	if err != nil {
		return nil, err
	}
	return str, nil
}

func CreateSubjectFromRecord(tenantId string, record *api.Record) string {
	if record.Namespace == "" {
		return fmt.Sprintf("%s.%s.cluster.%s.%s", tenantId, record.Cluster, record.Kind.Kind, record.Name)
	}
	return fmt.Sprintf("%s.%s.namespace.%s.%s.%s", tenantId, record.Cluster, record.Namespace, record.Kind.Kind, record.Name)
}
