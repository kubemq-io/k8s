package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// CorsConfig is a shared CORS policy applied to an HTTP-family server (the shared
// HTTP server and the REST connector). It is emitted under a caller-supplied env
// prefix (e.g. CONNECTORS_HTTP_CORS or CONNECTORS_REST_CORS).
//
// List fields are emitted ONLY when non-empty, so an empty list never overrides
// the server's non-empty default (which would otherwise fail CORS validation).
type CorsConfig struct {
	// +optional
	AllowOrigins []string `json:"allowOrigins,omitempty" yaml:"allowOrigins,omitempty"`

	// +optional
	AllowMethods []string `json:"allowMethods,omitempty" yaml:"allowMethods,omitempty"`

	// +optional
	AllowHeaders []string `json:"allowHeaders,omitempty" yaml:"allowHeaders,omitempty"`

	// +optional
	AllowCredentials *bool `json:"allowCredentials,omitempty" yaml:"allowCredentials,omitempty"`

	// +optional
	ExposeHeaders []string `json:"exposeHeaders,omitempty" yaml:"exposeHeaders,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxAge *int32 `json:"maxAge,omitempty" yaml:"maxAge,omitempty"`
}

func (c *CorsConfig) DeepCopy() *CorsConfig {
	out := &CorsConfig{}

	if c.AllowOrigins != nil {
		out.AllowOrigins = make([]string, len(c.AllowOrigins))
		copy(out.AllowOrigins, c.AllowOrigins)
	}

	if c.AllowMethods != nil {
		out.AllowMethods = make([]string, len(c.AllowMethods))
		copy(out.AllowMethods, c.AllowMethods)
	}

	if c.AllowHeaders != nil {
		out.AllowHeaders = make([]string, len(c.AllowHeaders))
		copy(out.AllowHeaders, c.AllowHeaders)
	}

	if c.AllowCredentials != nil {
		out.AllowCredentials = new(bool)
		*out.AllowCredentials = *c.AllowCredentials
	}

	if c.ExposeHeaders != nil {
		out.ExposeHeaders = make([]string, len(c.ExposeHeaders))
		copy(out.ExposeHeaders, c.ExposeHeaders)
	}

	if c.MaxAge != nil {
		out.MaxAge = new(int32)
		*out.MaxAge = *c.MaxAge
	}

	return out
}

// setConfig emits the CORS env vars under the given prefix. Lists are guarded by
// len>0 so an empty list is never written; scalar pointers emit only when non-nil.
func (c *CorsConfig) setConfig(config *deployment.Config, prefix string) {
	if len(c.AllowOrigins) > 0 {
		config.SetConfigMapStringValues(config.Name, prefix+"_ALLOW_ORIGINS", strings.Join(c.AllowOrigins, ","))
	}

	if len(c.AllowMethods) > 0 {
		config.SetConfigMapStringValues(config.Name, prefix+"_ALLOW_METHODS", strings.Join(c.AllowMethods, ","))
	}

	if len(c.AllowHeaders) > 0 {
		config.SetConfigMapStringValues(config.Name, prefix+"_ALLOW_HEADERS", strings.Join(c.AllowHeaders, ","))
	}

	if c.AllowCredentials != nil {
		config.SetConfigMapStringValues(config.Name, prefix+"_ALLOW_CREDENTIALS", strconv.FormatBool(*c.AllowCredentials))
	}

	if len(c.ExposeHeaders) > 0 {
		config.SetConfigMapStringValues(config.Name, prefix+"_EXPOSE_HEADERS", strings.Join(c.ExposeHeaders, ","))
	}

	if c.MaxAge != nil {
		config.SetConfigMapStringValues(config.Name, prefix+"_MAX_AGE", fmt.Sprintf("%d", *c.MaxAge))
	}
}
