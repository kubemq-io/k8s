package config

import (
	"fmt"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// HttpConfig configures the kubemq-server shared HTTP server (used by the
// HTTP-family connectors) and its CORS policy. Maps to server Connectors.HTTP.
type HttpConfig struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port *int32 `json:"port,omitempty" yaml:"port,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	ReadTimeout *int32 `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty"`

	// +optional
	BodyLimit *string `json:"bodyLimit,omitempty" yaml:"bodyLimit,omitempty"`

	// +optional
	BaseURL *string `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty"`

	// +optional
	Cors *CorsConfig `json:"cors,omitempty" yaml:"cors,omitempty"`
}

func (c *HttpConfig) DeepCopy() *HttpConfig {
	out := &HttpConfig{}

	if c.Port != nil {
		out.Port = new(int32)
		*out.Port = *c.Port
	}

	if c.ReadTimeout != nil {
		out.ReadTimeout = new(int32)
		*out.ReadTimeout = *c.ReadTimeout
	}

	if c.BodyLimit != nil {
		out.BodyLimit = new(string)
		*out.BodyLimit = *c.BodyLimit
	}

	if c.BaseURL != nil {
		out.BaseURL = new(string)
		*out.BaseURL = *c.BaseURL
	}

	if c.Cors != nil {
		out.Cors = c.Cors.DeepCopy()
	}

	return out
}

func (c *HttpConfig) SetConfig(config *deployment.Config) *HttpConfig {
	if c.Port != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_HTTP_PORT", fmt.Sprintf("%d", *c.Port))
	}

	if c.ReadTimeout != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_HTTP_READ_TIMEOUT", fmt.Sprintf("%d", *c.ReadTimeout))
	}

	if c.BodyLimit != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_HTTP_BODY_LIMIT", *c.BodyLimit)
	}

	if c.BaseURL != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_HTTP_BASE_URL", *c.BaseURL)
	}

	if c.Cors != nil {
		c.Cors.setConfig(config, "CONNECTORS_HTTP_CORS")
	}

	return c
}
