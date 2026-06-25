package config

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationConfig_OidcTypeOnly_NoJwtSignature(t *testing.T) {
	cfg := newTestConfig()
	// Type=oidc selector alone (no Key/SignatureType): JWT-path Type is written,
	// but NO JWT key/signature secrets are emitted, and no OIDC config either
	// (Oidc is nil so this falls through the explicit JWT/Type path).
	(&AuthenticationConfig{Type: strptr("oidc")}).SetConfig(cfg)

	v := vars(cfg)
	assert.Equal(t, "oidc", v["AUTHENTICATION_TYPE"])

	sec := cfg.Secrets[testClusterName]
	_, hasKey := sec.DataVariables["AUTHENTICATION_JWT_CONFIG_KEY"]
	_, hasSig := sec.StringVariables["AUTHENTICATION_JWT_CONFIG_SIGNATURE_TYPE"]
	assert.False(t, hasKey, "no JWT key must be emitted when only Type is set")
	assert.False(t, hasSig, "no JWT signature type must be emitted when only Type is set")
}

func TestAuthenticationConfig_Oidc(t *testing.T) {
	cfg := newTestConfig()
	(&AuthenticationConfig{
		Oidc: &OidcConfig{
			Issuer:   "https://issuer.example.com",
			ClientID: "kubemq",
		},
	}).SetConfig(cfg)

	v := vars(cfg)
	assert.Equal(t, "true", v["AUTHENTICATION_ENABLE"])
	assert.Equal(t, "oidc", v["AUTHENTICATION_TYPE"])
	// AUTHENTICATION_CONFIG is written to the ConfigMap as a data value (base64).
	require.Contains(t, v, "AUTHENTICATION_CONFIG")
	decoded, err := base64.StdEncoding.DecodeString(v["AUTHENTICATION_CONFIG"])
	require.NoError(t, err)
	assert.JSONEq(t,
		`{"issuer":"https://issuer.example.com","clientID":"kubemq","skipClientIDCheck":false,"skipExpiryCheck":false,"skipIssuerCheck":false,"insecureSkipSignatureCheck":false}`,
		string(decoded))

	// No JWT secrets should be present in the OIDC path.
	sec := cfg.Secrets[testClusterName]
	_, hasKey := sec.DataVariables["AUTHENTICATION_JWT_CONFIG_KEY"]
	_, hasSig := sec.StringVariables["AUTHENTICATION_JWT_CONFIG_SIGNATURE_TYPE"]
	assert.False(t, hasKey)
	assert.False(t, hasSig)
}

func TestAuthenticationConfig_Jwt(t *testing.T) {
	cfg := newTestConfig()
	(&AuthenticationConfig{
		Enable:        ptrBool(true),
		Key:           strptr("my-jwt-key"),
		SignatureType: strptr("HS256"),
	}).SetConfig(cfg)

	v := vars(cfg)
	assert.Equal(t, "true", v["AUTHENTICATION_ENABLE"])
	// No AUTHENTICATION_TYPE in the JWT path (Type not set).
	_, hasType := v["AUTHENTICATION_TYPE"]
	assert.False(t, hasType, "AUTHENTICATION_TYPE must not be written for the JWT path when Type is unset")

	sec := cfg.Secrets[testClusterName]
	// Key is a Secret DATA value (base64-encoded).
	keyEnc, hasKey := sec.DataVariables["AUTHENTICATION_JWT_CONFIG_KEY"]
	require.True(t, hasKey, "AUTHENTICATION_JWT_CONFIG_KEY must be in the Secret data")
	keyDec, err := base64.StdEncoding.DecodeString(keyEnc)
	require.NoError(t, err)
	assert.Equal(t, "my-jwt-key", string(keyDec))
	// SignatureType is a Secret STRING value (raw).
	assert.Equal(t, "HS256", sec.StringVariables["AUTHENTICATION_JWT_CONFIG_SIGNATURE_TYPE"])
}

func TestAuthenticationConfig_MutuallyExclusive_EmitsNothing(t *testing.T) {
	cfg := newTestConfig()
	// Oidc + JWT settings together is the error case: nothing must be emitted.
	(&AuthenticationConfig{
		Oidc:          &OidcConfig{Issuer: "https://x", ClientID: "y"},
		Key:           strptr("k"),
		SignatureType: strptr("HS256"),
	}).SetConfig(cfg)

	assert.Len(t, vars(cfg), 0)
	sec := cfg.Secrets[testClusterName]
	assert.Len(t, sec.DataVariables, 0)
	assert.Len(t, sec.StringVariables, 0)
}

func TestAuthenticationConfig_DeepCopy_Independent(t *testing.T) {
	src := &AuthenticationConfig{
		Enable:        ptrBool(true),
		Type:          strptr("oidc"),
		Key:           strptr("k"),
		SignatureType: strptr("HS256"),
		Oidc:          &OidcConfig{Issuer: "https://x", ClientID: "y"},
	}
	cp := src.DeepCopy()
	*src.Enable = false
	*src.Type = "mutated"
	*src.Key = "mutated"
	*src.SignatureType = "mutated"
	src.Oidc.Issuer = "mutated"
	src.Oidc = nil

	assert.Equal(t, true, *cp.Enable)
	assert.Equal(t, "oidc", *cp.Type)
	assert.Equal(t, "k", *cp.Key)
	assert.Equal(t, "HS256", *cp.SignatureType)
	require.NotNil(t, cp.Oidc)
	assert.Equal(t, "https://x", cp.Oidc.Issuer)
}
