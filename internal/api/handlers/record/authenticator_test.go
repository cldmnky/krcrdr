// BEGIN: 4c5d6f7g8h9j
package record

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/assert"
)

type mockValidator struct{}

func (m *mockValidator) ValidateJWS(jws string) (jwt.Token, error) {
	if jws == "valid_token" {
		return jwt.New(), nil
	}
	return jwt.New(), errors.New("invalid token")
}

func TestGetJWSFromRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	// Test missing auth header
	_, err := GetJWSFromRequest(req)
	assert.Equal(t, ErrNoAuthHeader, err)

	// Test invalid auth header
	req.Header.Set("Authorization", "invalid")
	_, err = GetJWSFromRequest(req)
	assert.Equal(t, ErrInvalidAuthHeader, err)

	// Test valid auth header
	req.Header.Set("Authorization", "Bearer valid_token")
	jws, err := GetJWSFromRequest(req)
	assert.NoError(t, err)
	assert.Equal(t, "valid_token", jws)
}

func TestAuthenticate(t *testing.T) {
	validator := &mockValidator{}

	// Test unsupported security scheme
	input := &openapi3filter.AuthenticationInput{
		SecuritySchemeName: "InvalidScheme",
	}
	err := Authenticate(context.TODO(), input, validator)
	assert.EqualError(t, err, "unsupported security scheme: InvalidScheme")

	// Test invalid token
	input = &openapi3filter.AuthenticationInput{
		SecuritySchemeName: "BearerAuth",
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request: &http.Request{},
		},
	}
	err = Authenticate(context.TODO(), input, validator)
	assert.EqualError(t, err, "failed to get JWS from request: authorization header is missing")

	// Test valid token
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	input = &openapi3filter.AuthenticationInput{
		SecuritySchemeName: "BearerAuth",
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request: req,
		},
	}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	// initialize gin context
	c.Request = req
	c.Set("oapi-codegen/gin-context", c)

	err = Authenticate(c, input, validator)
	assert.NoError(t, err)
}

func TestCheckTokenClaims(t *testing.T) {
	tests := []struct {
		name           string
		expectedClaims []string
		token          jwt.Token
		wantErr        bool
		tokenClaims    []interface{}
		tokenClaimName string
	}{
		{
			name:           "valid claims",
			expectedClaims: []string{"read", "write"},
			token:          jwt.New(),
			tokenClaims:    []interface{}{"read", "write"},
			tokenClaimName: PermissonClaim,
			wantErr:        false,
		},
		{
			name:           "missing claims",
			expectedClaims: []string{"read", "write", "delete"},
			token:          jwt.New(),
			tokenClaims:    []interface{}{"read", "write"},
			tokenClaimName: PermissonClaim,
			wantErr:        true,
		},
		{
			name:           "extra claims",
			expectedClaims: []string{"read", "write"},
			token:          jwt.New(),
			tokenClaims:    []interface{}{"read", "write", "delete"},
			tokenClaimName: PermissonClaim,
			wantErr:        false,
		},
		{
			name:           "empty claims",
			expectedClaims: []string{},
			token:          jwt.New(),
			tokenClaims:    []interface{}{},
			tokenClaimName: PermissonClaim,
			wantErr:        false,
		},
		{
			name:           "missing perms claim",
			expectedClaims: []string{"read", "write"},
			token:          jwt.New(),
			tokenClaims:    []interface{}{},
			tokenClaimName: "notPerms",
			wantErr:        true,
		},
		{
			name:           "invalid perms type",
			expectedClaims: []string{"read", "write"},
			token:          jwt.New(),
			tokenClaims:    []interface{}{},
			tokenClaimName: PermissonClaim,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.token.Set(tt.tokenClaimName, tt.tokenClaims)
			err := checkTokenClaims(tt.expectedClaims, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
