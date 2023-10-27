package nats

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"

	"github.com/cldmnky/krcrdr/internal/api/store"
)

type NatsStore struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	kv     jetstream.KeyValue
	tracer trace.Tracer
}

// kvEntry is each key-value entry in the store and the operation associated with the kv pair.
type kvEntry struct {
	bucket    string
	key       string
	value     []byte
	revision  uint64
	created   time.Time
	operation store.KVWatchOp
}

// Bucket returns the bucket.
func (k kvEntry) Bucket() string {
	return k.bucket
}

// Key returns the key
func (k kvEntry) Key() string {
	return k.key
}

// Value returns the value.
func (k kvEntry) Value() []byte {
	return k.value
}

// Revision returns the unique revision of the key-value pair.
func (k kvEntry) Revision() uint64 {
	return k.revision
}

// Created returns the time the key-value pair was created.
func (k kvEntry) Created() time.Time {
	return k.created
}

// Operation returns the operation on that key-value pair.
func (k kvEntry) Operation() store.KVWatchOp {
	return k.operation
}
