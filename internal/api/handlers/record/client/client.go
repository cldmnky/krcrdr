package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/go-retryablehttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
)

var logger = logf.Log.WithName("retryableClient")

// NewApiClient creates a new instance of the api.Client struct with the provided endpoint, bearerToken, and insecure flag.
// If insecure is true, TLS errors will be ignored.
// It returns a pointer to the api.Client struct and an error if any.
func NewApiClient(endpoint, bearerToken string, insecure bool) (*api.Client, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}
	if bearerToken == "" {
		return nil, fmt.Errorf("bearerToken cannot be empty")
	}
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(bearerToken)
	if err != nil {
		return nil, err
	}
	// setup a retryable client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Logger = logger
	// setup a transport that will ignore TLS errors if insecure is true
	if insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		// add the transport to the http client
		retryClient.HTTPClient.Transport = tr
	}

	tr := otelhttp.NewTransport(retryClient.HTTPClient.Transport,
		otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
			return otelhttptrace.NewClientTrace(ctx)
		}),
	)
	retryClient.HTTPClient.Transport = tr

	return api.NewClient(endpoint, api.WithRequestEditorFn(bearerTokenProvider.Intercept), api.WithHTTPClient(retryClient.HTTPClient))
}
