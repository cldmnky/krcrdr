package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// StreamService is the interface for a stream store.
func NewStream(url string, options ...nats.Option) (*NatsStore, error) {
	nc, err := nats.Connect(url, options...)
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
