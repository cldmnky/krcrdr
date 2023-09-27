package record

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/ecdsafile"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

const PrivateKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIN2dALnjdcZaIZg4QuA6Dw+kxiSW502kJfmBN3priIhPoAoGCCqGSM49
AwEHoUQDQgAE4pPyvrB9ghqkT1Llk0A42lixkugFd/TBdOp6wf69O9Nndnp4+HcR
s9SlG/8hjB2Hz42v4p3haKWv3uS1C6ahCQ==
-----END EC PRIVATE KEY-----` // notsecret

const KeyID = `fake-key-id`
const FakeIssuer = "fake-issuer"
const FakeAudience = "example-users"
const PermissionsClaim = "perm"

type FakeAuthenticator struct {
	PrivateKey *ecdsa.PrivateKey
	KeySet     jwk.Set
}

var _ JWSValidator = (*FakeAuthenticator)(nil)

func (f *FakeAuthenticator) ValidateJWS(jwsString string) (jwt.Token, error) {
	return jwt.Parse([]byte(jwsString), jwt.WithKeySet(f.KeySet), jwt.WithAudience(FakeAudience), jwt.WithIssuer(FakeIssuer))
}

func NewFakeAuthenticator() (*FakeAuthenticator, error) {
	privKey, err := ecdsafile.LoadEcdsaPrivateKey([]byte(PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}
	set := jwk.NewSet()
	pubKey := jwk.NewECDSAPublicKey()
	if err := pubKey.FromRaw(&privKey.PublicKey); err != nil {
		return nil, fmt.Errorf("failed to create public key: %w", err)
	}

	if err := pubKey.Set(jwk.AlgorithmKey, jwa.ES256); err != nil {
		return nil, fmt.Errorf("failed to set algorithm: %w", err)
	}

	if err := pubKey.Set(jwk.KeyIDKey, KeyID); err != nil {
		return nil, fmt.Errorf("failed to set key id: %w", err)
	}

	set.Add(pubKey)

	return &FakeAuthenticator{
		PrivateKey: privKey,
		KeySet:     set,
	}, nil
}

// SignToken takes a jwt token and signs it with the private key, returning a JWS
func (f *FakeAuthenticator) SignToken(token jwt.Token) ([]byte, error) {
	header := jws.NewHeaders()
	// Set Key ID
	if err := header.Set(jwk.KeyIDKey, KeyID); err != nil {
		return nil, fmt.Errorf("failed to set key id: %w", err)
	}
	// Set Algorithm
	if err := header.Set(jwk.AlgorithmKey, jwa.ES256); err != nil {
		return nil, fmt.Errorf("failed to set algorithm: %w", err)
	}
	// Set type
	if err := header.Set(jws.TypeKey, "JWT"); err != nil {
		return nil, fmt.Errorf("failed to set type: %w", err)
	}

	return jwt.Sign(token, jwa.ES256, f.PrivateKey, jwt.WithHeaders(header))
}

// CreateJWSWithClaims is a helper function to create JWT's with the specified
// claims.
func (f *FakeAuthenticator) CreateJWSWithClaims(claims []string) ([]byte, error) {
	t := jwt.New()
	if err := t.Set(jwt.IssuerKey, FakeIssuer); err != nil {
		return nil, fmt.Errorf("failed to set issuer: %w", err)
	}
	if err := t.Set(jwt.AudienceKey, FakeAudience); err != nil {
		return nil, fmt.Errorf("failed to set audience: %w", err)
	}
	if err := t.Set(PermissionsClaim, claims); err != nil {
		return nil, fmt.Errorf("failed to set permissions claim: %w", err)
	}
	return f.SignToken(t)
}
