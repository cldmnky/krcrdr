package options

import (
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"
)

type RootOptions struct {
	// Version is the version of the application.
	Version string
	// OTLPAddr is the address of the OTLP collector.
	OTLPAddr string
	// Exporter is the name of the exporter to use. Defaults to noop.
	OTLPExporter string
	// Debug enables debug logging.
	Debug bool
	// Tracer
	Tracer trace.Tracer
}

var _ Interface = (*RootOptions)(nil)

func (o *RootOptions) AddFlags(cmd *cobra.Command) {
	// Global flags
	cmd.PersistentFlags().StringVar(&o.OTLPAddr, "otlp-addr", "localhost:4317", "The address of the OTLP tracing collector.")
	cmd.PersistentFlags().StringVar(&o.OTLPExporter, "otlp-exporter", "noop", "The name of the tracing exporter to use. Valid options are otel, console. Defaults to noop.")
	cmd.PersistentFlags().BoolVar(&o.Debug, "debug", false, "Enable debug logging")
	// Set version
	cmd.Version = o.Version
}
