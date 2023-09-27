package record

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/lestrrat-go/jwx/jwt"
	middleware "github.com/oapi-codegen/gin-middleware"
)

type JWSValidator interface {
	ValidateJWS(jws string) (jwt.Token, error)
}

const (
	JWTClaimsContextKey = "jwt_claims"
	PermissonClaim      = "perms"
)

var (
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("authorization header is malformed")
	ErrClaimsInvalid     = errors.New("provided claims do not match expected scopes")
)

func GetJWSFromRequest(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func NewAuthenticator(validator JWSValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(ctx, input, validator)
	}
}

func Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput, validator JWSValidator) error {
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("unsupported security scheme: %s", input.SecuritySchemeName)
	}

	// Get the JWS from the request
	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("failed to get JWS from request: %w", err)
	}

	// Validate the JWS
	claims, err := validator.ValidateJWS(jws)
	if err != nil {
		return fmt.Errorf("failed to validate JWS: %w", err)
	}

	err = checkTokenClaims(input.Scopes, claims)
	if err != nil {
		return fmt.Errorf("failed to check token claims: %w", err)
	}

	// Add the claims to the context
	gCtx := middleware.GetGinContext(ctx)
	gCtx.Set(JWTClaimsContextKey, claims)

	return nil
}

// getClaimsFromToken returns a list of claims from the token. We store these
// as a list under the "perms" claim, short for permissions, to keep the token
// shorter.
func getClaimsFromToken(t jwt.Token) ([]string, error) {
	rawPerms, ok := t.Get(PermissionsClaim)
	if !ok {
		return make([]string, 0), nil
	}
	// convert the interface{} to a []interface{}
	perms, ok := rawPerms.([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert perms, unexpected type: %T", rawPerms)
	}

	// convert the []interface{} to a []string
	claims := make([]string, len(perms))

	// make sure each claim is a string
	for i, p := range perms {
		claim, ok := p.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert claim, unexpected type: %T", p)
		}
		claims[i] = claim
	}
	return claims, nil
}

func checkTokenClaims(expectedClaims []string, t jwt.Token) error {
	claims, err := getClaimsFromToken(t)
	if err != nil {
		return fmt.Errorf("failed to get claims from token: %w", err)
	}

	// put claims into a map for easy lookup
	claimsMap := make(map[string]bool)
	for _, c := range claims {
		claimsMap[c] = true
	}

	// check that all expected claims are present
	for _, e := range expectedClaims {
		if !claimsMap[e] {
			return fmt.Errorf("provided claims do not match expected scopes: %s", e)
		}
	}
	return nil
}
