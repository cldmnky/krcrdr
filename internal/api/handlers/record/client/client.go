package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptrace"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/go-retryablehttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewApiClient(endpoint, bearerToken string, insecure bool) (*api.Client, error) {
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(bearerToken)
	if err != nil {
		return nil, err
	}
	// setup a retryable client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
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
