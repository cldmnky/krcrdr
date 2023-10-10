package tracer

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
)

func TestNewExporter(t *testing.T) {
	tests := []struct {
		name        string
		exType      ExporterType
		otlpAddress string
		wantErr     bool
	}{
		{
			name:        "OTLP exporter",
			exType:      ExporterTypeOTLP,
			otlpAddress: "localhost:4317",
			wantErr:     false,
		},
		{
			name:        "Stdio exporter",
			exType:      ExporterTypeConsole,
			otlpAddress: "",
			wantErr:     false,
		},
		{
			name:        "Noop exporter",
			exType:      ExporterTypeNoop,
			otlpAddress: "",
			wantErr:     false,
		},
		{
			name:        "Unknown exporter",
			exType:      "unknown",
			otlpAddress: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewExporter(string(tt.exType), tt.otlpAddress, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExporter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name     string
		exporter Exporter
		wantErr  bool
	}{
		{
			name:     "OTLP exporter",
			exporter: &otlptrace.Exporter{},
			wantErr:  false,
		},
		{
			name:     "Console exporter",
			exporter: &consoleExporter{},
		},
		{
			name:     "Noop exporter",
			exporter: &NoopExporter{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProvider(context.Background(), "version", tt.exporter)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
