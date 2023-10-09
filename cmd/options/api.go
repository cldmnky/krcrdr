package options

import (
	apiopts "github.com/cldmnky/krcrdr/internal/api"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var _ Interface = (*ControllerOptions)(nil)

type ApiOptions struct {
	// Address to listen on for the API server. Default: ":8443"
	Addr string
	// CertDir is the directory containing the TLS certs. Default: "/tmp/krcrdr/serving-certss"
	CertDir string
	// CertName is the name of the TLS cert. Default: "tls.crt"
	CertName string
	// KeyName is the name of the TLS key. Default: "tls.key"
	KeyName string
	// ClientCAName is the name of the CA certificate used to verify client certificates. Default: ""
	ClientCAName string
	// Debug enables debug logging. Default: false
	Debug bool
	// GenerateSelfSignedCert generates a self-signed certificate and key. Default: false
	GenerateSelfSignedCert bool
	// RunNatsServer runs an embedded NATS server, should only be done for testing. Default: false
	RunNatsServer bool
	// NatsUrl is the NATS server URL. Default: "nats://localhost:4222"
	NatsUrl []string
	// NatsUserCredentials is the NATS user credentials. Default: ""
	NatsUserCredentials string
}

func (o *ApiOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Addr, "addr", ":8443", "Address to listen on for the API server")
	cmd.Flags().StringVar(&o.CertDir, "cert-dir", apiopts.DefaultCertDir, "Directory containing the TLS certs")
	cmd.Flags().StringVar(&o.CertName, "cert-name", apiopts.DefaultCertName, "Name of the TLS cert")
	cmd.Flags().StringVar(&o.KeyName, "key-name", apiopts.DefaultKeyName, "Name of the TLS key")
	cmd.Flags().StringVar(&o.ClientCAName, "client-ca-name", "", "Name of the CA certificate used to verify client certificates")
	cmd.Flags().BoolVar(&o.GenerateSelfSignedCert, "generate-self-signed-cert", false, "Generate a self-signed certificate and key")
	cmd.Flags().BoolVar(&o.RunNatsServer, "run-nats-server", false, "Run an embedded NATS server, should only be done for testing")
	cmd.Flags().StringSliceVar(&o.NatsUrl, "nats-url", []string{nats.DefaultURL}, "NATS server URL(s)")
	cmd.Flags().StringVar(&o.NatsUserCredentials, "nats-user-credentials", "", "NATS user credentials")
	cmd.Flags().BoolVar(&o.Debug, "debug", false, "Enable debug logging")
}
