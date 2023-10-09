package options

import (
	"github.com/spf13/cobra"
)

type WebhookOptions struct {
	// Addr is the address the webhook server binds to.
	Addr string
	// CertDir is the directory containing the TLS certs.
	CertDir string
	// MetricsAddr is the address the metric endpoint binds to.
	MetricsAddr string
	// ProbeAddr is the address the probe endpoint binds to.
	ProbeAddr string
	// Debug enables debug logging.
	Debug bool
	// ApiRemoteAddr is the address of the API endpoint.
	ApiRemoteAddr string
	// ApiToken is the token used to authenticate to the API.
	ApiToken string
	// GenerateCerts generates the TLS certs.
	GenerateCerts bool
	// CertName is the name of the cert.
	CertName string
	// KeyName is the name of the key.
	KeyName string
}

var _ Interface = (*WebhookOptions)(nil)

func (o *WebhookOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Addr, "addr", ":9443", "The address the webhook server binds to.")
	cmd.Flags().StringVar(&o.CertDir, "cert-dir", "/tmp/k8s-webhook-server/serving-certs", "The directory containing the TLS certs.")
	cmd.Flags().StringVar(&o.MetricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	cmd.Flags().StringVar(&o.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	cmd.Flags().StringVar(&o.ApiRemoteAddr, "api-remote-address", "http://localhost:8082", "The address of the API endpoint. Env var: WEBHOOK_API_REMOTE_ADDRESS")
	cmd.Flags().StringVar(&o.ApiToken, "api-token", "", "The token used to authenticate to the API. Env var: WEBHOOK_API_TOKEN")
	cmd.Flags().BoolVar(&o.GenerateCerts, "generate-certs", false, "Generate the TLS certs.")
	cmd.Flags().StringVar(&o.CertName, "cert-name", "tls.crt", "The name of the cert.")
	cmd.Flags().StringVar(&o.KeyName, "key-name", "tls.key", "The name of the key.")
	cmd.Flags().BoolVar(&o.Debug, "debug", false, "Enable debug logging")
	// Set apitoken to required
	cmd.MarkFlagRequired("api-token")
}
