package config

import (
	"encoding/json"
	"strconv"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// OidcConfig mirrors the server-side OIDCConfig (see config_data.go OIDCConfig)
// exactly, so its JSON serialization matches what the server expects in
// AUTHENTICATION_CONFIG.
type OidcConfig struct {
	Issuer                     string `json:"issuer" yaml:"issuer"`
	ClientID                   string `json:"clientID" yaml:"clientID"`
	SkipClientIDCheck          bool   `json:"skipClientIDCheck" yaml:"skipClientIDCheck"`
	SkipExpiryCheck            bool   `json:"skipExpiryCheck" yaml:"skipExpiryCheck"`
	SkipIssuerCheck            bool   `json:"skipIssuerCheck" yaml:"skipIssuerCheck"`
	InsecureSkipSignatureCheck bool   `json:"insecureSkipSignatureCheck" yaml:"insecureSkipSignatureCheck"`
}

func (o *OidcConfig) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}

func (o *OidcConfig) DeepCopy() *OidcConfig {
	out := &OidcConfig{}
	*out = *o
	return out
}

// AuthenticationConfig configures the kubemq-server authentication. Maps to
// server Authentication. Two mutually-exclusive modes:
//   - JWT  (Type empty): Key + SignatureType select the JWT verification key/algorithm.
//   - OIDC (Type "oidc"): Oidc carries the OIDC verification settings.
type AuthenticationConfig struct {
	// +optional
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`

	// +optional
	// Type is the server MODE selector: empty=jwt, "oidc"=oidc.
	Type *string `json:"type,omitempty" yaml:"type,omitempty"`

	// +optional
	// Key is the JWT verification key.
	Key *string `json:"key,omitempty" yaml:"key,omitempty"`

	// +optional
	// SignatureType is the JWT signing algorithm (e.g. HS256, RS256).
	SignatureType *string `json:"signatureType,omitempty" yaml:"signatureType,omitempty"`

	// +optional
	Oidc *OidcConfig `json:"oidc,omitempty" yaml:"oidc,omitempty"`
}

func (c *AuthenticationConfig) DeepCopy() *AuthenticationConfig {
	out := &AuthenticationConfig{}

	if c.Enable != nil {
		out.Enable = new(bool)
		*out.Enable = *c.Enable
	}

	if c.Type != nil {
		out.Type = new(string)
		*out.Type = *c.Type
	}

	if c.Key != nil {
		out.Key = new(string)
		*out.Key = *c.Key
	}

	if c.SignatureType != nil {
		out.SignatureType = new(string)
		*out.SignatureType = *c.SignatureType
	}

	if c.Oidc != nil {
		out.Oidc = c.Oidc.DeepCopy()
	}

	return out
}

func (c *AuthenticationConfig) SetConfig(config *deployment.Config) *AuthenticationConfig {
	// Mutually-exclusive error case: OIDC and JWT settings cannot both be set.
	// Emit nothing and return — reconcile-level validation should reject this
	// before it ever reaches SetConfig; guarding here avoids emitting a
	// contradictory mix of OIDC and JWT env vars.
	if c.Oidc != nil && (c.Key != nil || c.SignatureType != nil) {
		return c
	}

	// OIDC path.
	if c.Oidc != nil {
		config.SetConfigMapStringValues(config.Name, "AUTHENTICATION_ENABLE", "true")
		config.SetConfigMapStringValues(config.Name, "AUTHENTICATION_TYPE", "oidc")
		config.SetConfigMapDataValues(config.Name, "AUTHENTICATION_CONFIG", c.Oidc.String())
		return c
	}

	// JWT / explicit path.
	if c.Enable != nil {
		config.SetConfigMapStringValues(config.Name, "AUTHENTICATION_ENABLE", strconv.FormatBool(*c.Enable))
	}

	if c.Type != nil {
		config.SetConfigMapStringValues(config.Name, "AUTHENTICATION_TYPE", *c.Type)
	}

	if c.Key != nil {
		config.SetSecretDataValues(config.Name, "AUTHENTICATION_JWT_CONFIG_KEY", *c.Key)
	}

	if c.SignatureType != nil {
		config.SetSecretStringValues(config.Name, "AUTHENTICATION_JWT_CONFIG_SIGNATURE_TYPE", *c.SignatureType)
	}

	return c
}
