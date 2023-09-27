package record

import (
	"fmt"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
)

func CreateApiMiddleware(v JWSValidator) ([]gin.HandlerFunc, error) {
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
			AuthenticationFunc: NewAuthenticator(v),
		},
	}
	return []gin.HandlerFunc{middleware.OapiRequestValidatorWithOptions(spec, &options)}, nil
}
