package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
)

type NatsStore struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	kv     jetstream.KeyValue
	tracer trace.Tracer
}
