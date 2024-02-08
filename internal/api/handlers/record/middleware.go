package record

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/mitchellh/mapstructure"
	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/cldmnky/krcrdr/internal/api/auth"
	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
)

func CreateApiMiddleware(v auth.JWSValidator) ([]gin.HandlerFunc, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}
	options := middleware.Options{
		ErrorHandler: func(c *gin.Context, message string, statusCode int) {
			c.String(statusCode, "error: "+message)
		},
		MultiErrorHandler: func(me openapi3.MultiError) error {
			return fmt.Errorf("errors: %s", me.Error())
		},
		Options: openapi3filter.Options{
			MultiError:         true,
			AuthenticationFunc: auth.NewAuthenticator(v),
		},
	}
	return []gin.HandlerFunc{middleware.OapiRequestValidatorWithOptions(spec, &options), TenantMiddleware()}, nil
}

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant, err := getTenantClaimsFromContext(c)
		if err != nil {
			c.IndentedJSON(400, err)
			return
		}
		c.Set("tenant", tenant)
		c.Next()
	}
}

func getTenantClaimsFromContext(c *gin.Context) (*auth.Tenant, error) {
	ck, ok := c.Get(auth.JWTClaimsContextKey)
	if !ok {
		return nil, auth.ErrClaimsInvalid
	}
	// Get privateClaims from ck
	pc := ck.(jwt.Token).PrivateClaims()
	// Get the tenant claim from the private claims
	t, ok := pc[auth.TenantClaim]
	if !ok {
		return nil, fmt.Errorf("tenant claim not found")
	}
	var tenant auth.Tenant
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &tenant,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	if err := decoder.Decode(t); err != nil {
		return nil, fmt.Errorf("failed to decode tenant claim: %w", err)
	}
	return &tenant, nil
}
