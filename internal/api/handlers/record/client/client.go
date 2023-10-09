package client

import (
	"crypto/tls"
	"net/http"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/go-retryablehttp"
)

func NewApiClient(endpoint, bearerToken string, insecure bool) (*api.Client, error) {
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(bearerToken)
	if err != nil {
		return nil, err
	}
	// setup a retryable client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	httpClient := retryClient.StandardClient()
	// setup a transport that will ignore TLS errors if insecure is true
	if insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient.Transport = tr
	}
	return api.NewClient(endpoint, api.WithRequestEditorFn(bearerTokenProvider.Intercept), api.WithHTTPClient(httpClient))
}
