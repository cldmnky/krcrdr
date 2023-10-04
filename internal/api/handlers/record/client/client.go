package client

import (
	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/go-retryablehttp"
)

func NewApiClient(endpoint, bearerToken string) (*api.Client, error) {
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(bearerToken)
	if err != nil {
		return nil, err
	}
	// setup a retryable client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	return api.NewClient(endpoint, api.WithRequestEditorFn(bearerTokenProvider.Intercept), api.WithHTTPClient(retryClient.StandardClient()))
}
