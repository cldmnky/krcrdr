package tracer

import (
	"context"
	"fmt"
	"io"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type (
	ExporterType string
	SamplerType  string
)

const (
	ExporterTypeOTEL    ExporterType = "otel"
	ExporterTypeConsole ExporterType = "console"
	ExporterTypeNoop    ExporterType = "noop"

	SamplerTypeAlways     SamplerType = "always"
	SamplerTypeNever      SamplerType = "never"
	SamplerTypeRatioBased SamplerType = "ratio_based"
)

type Exporter interface {
	sdktrace.SpanExporter

	Start(context.Context) error
}

// NewExporter creates a new Exporter based on the provided exporter type and OTLP address.
// If the exporter type is "otlp", a new GRPCExporter is created with the provided OTLP address.
// If the exporter type is "stdio", a new ConsoleExporter is created with os.Stdout.
// If the exporter type is unknown, a new NoopExporter is created and an error is returned.
func NewExporter(exType, otlpAddress string, writer io.Writer) (Exporter, error) {
	switch strings.ToLower(exType) {
	case string(ExporterTypeOTEL):
		return NewOTELExporter(otlpAddress)
	case string(ExporterTypeConsole):
		return NewConsoleExporter(writer)
	case string(ExporterTypeNoop):
		return NewNoopExporter(), nil
	default:
		return NewNoopExporter(), fmt.Errorf("unknown exporter type: %s", exType)
	}
}

type NoopExporter struct{}

func NewNoopExporter() *NoopExporter {
	return &NoopExporter{}
}
func (n NoopExporter) ExportSpans(_ context.Context, _ []sdktrace.ReadOnlySpan) error {
	return nil
}

func (n NoopExporter) MarshalLog() interface{} {
	return nil
}

func (n NoopExporter) Shutdown(_ context.Context) error {
	return nil
}

func (n NoopExporter) Start(_ context.Context) error {
	return nil
}

type consoleExporter struct {
	*stdouttrace.Exporter
}

func (c *consoleExporter) Start(_ context.Context) error {
	return nil
}

func NewOTELExporter(otlpAddress string) (Exporter, error) {
	return otlptracegrpc.NewUnstarted(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otlpAddress),
	), nil
}

func NewConsoleExporter(w io.Writer) (Exporter, error) {
	exp, err := stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
	if err != nil {
		return nil, err
	}
	return &consoleExporter{exp}, nil
}

func NewProvider(ctx context.Context, version string, exporter sdktrace.SpanExporter, opts ...sdktrace.TracerProviderOption) (trace.TracerProvider, error) {
	if exporter == nil {
		return trace.NewNoopTracerProvider(), nil
	}

	res, err := resources(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	provider := sdktrace.NewTracerProvider(
		// Sampler options:
		// - sdktrace.NeverSample()
		// - sdktrace.TraceIDRatioBased(0.01)
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
	)
	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(provider)

	return provider, nil
}

func resources(ctx context.Context, version string) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("parca"),
			semconv.ServiceVersionKey.String(version),
		),
		resource.WithFromEnv(),   // pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables
		resource.WithProcess(),   // This option configures a set of Detectors that discover process information
		resource.WithOS(),        // This option configures a set of Detectors that discover OS information
		resource.WithContainer(), // This option configures a set of Detectors that discover container information
		resource.WithHost(),      // This option configures a set of Detectors that discover host information
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	return res, nil
}
