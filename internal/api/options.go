package api

import (
	"crypto/tls"
	"os"
	"path/filepath"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record"
	"github.com/cldmnky/krcrdr/internal/api/store"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/trace"
)

var (
	DefaultPort     = 8443
	DefaultCertDir  = filepath.Join(os.TempDir(), "api/serving-certs")
	DefaultCertName = "tls.crt"
	DefaultKeyName  = "tls.key"
)

type Options struct {
	// Host is the address that the server will listen on
	// Defaults to "" which means all interfaces
	Host string
	// Port is the port that the server will listen on
	// Defaults to 8443
	Port int
	// CertDir is the directory that contains the TLS certificate and key
	// Defaults to <temp-dir>/krcrdr/serving-certs
	CertDir string
	// CertName is the name of the TLS certificate
	// Defaults to tls.crt
	CertName string
	// KeyName is the name of the TLS key
	// Defaults to tls.key
	KeyName string
	// ClientCAName is the name of the CA certificate used to verify client certificates
	// Defaults to "" which means no client certificate verification
	ClientCAName string
	// TLSOpts is used to allow configuring the TLS config used for the server
	TLSOpts []func(*tls.Config)
	// ApiLogger is the logger used for the api server
	ApiLogger logr.Logger
	// Authenticator is the authenticator used to verify incoming requests
	Authenticator record.JWSValidator
	// Store is the store used to persist records
	Store store.Store
	// Tracer is the tracer used for tracing
	Tracer trace.Tracer
	// Debug enables debug logging
	Debug bool
}

func (o *Options) setDefaults() {
	if o.Port == 0 {
		o.Port = DefaultPort
	}
	if o.CertDir == "" {
		o.CertDir = DefaultCertDir
	}
	if o.CertName == "" {
		o.CertName = DefaultCertName
	}
	if o.KeyName == "" {
		o.KeyName = DefaultKeyName
	}
}
