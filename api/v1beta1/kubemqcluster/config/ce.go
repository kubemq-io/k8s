package config

import (
	"fmt"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// CeConfig configures the kubemq-server CloudEvents (CE) connector.
// Maps to server Connectors.CE. The connector is enabled by default server-side.
type CeConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty" yaml:"timeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10000
	SubBuffSize *int32 `json:"subBuffSize,omitempty" yaml:"subBuffSize,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxSSEIdleSeconds *int32 `json:"maxSseIdleSeconds,omitempty" yaml:"maxSseIdleSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxSSEConnections *int32 `json:"maxSseConnections,omitempty" yaml:"maxSseConnections,omitempty"`
}

func (c *CeConfig) DeepCopy() *CeConfig {
	out := &CeConfig{}

	out.Disabled = c.Disabled

	if c.TimeoutSeconds != nil {
		out.TimeoutSeconds = new(int32)
		*out.TimeoutSeconds = *c.TimeoutSeconds
	}

	if c.SubBuffSize != nil {
		out.SubBuffSize = new(int32)
		*out.SubBuffSize = *c.SubBuffSize
	}

	if c.MaxSSEIdleSeconds != nil {
		out.MaxSSEIdleSeconds = new(int32)
		*out.MaxSSEIdleSeconds = *c.MaxSSEIdleSeconds
	}

	if c.MaxSSEConnections != nil {
		out.MaxSSEConnections = new(int32)
		*out.MaxSSEConnections = *c.MaxSSEConnections
	}

	return out
}

func (c *CeConfig) SetConfig(config *deployment.Config) *CeConfig {
	if c.Disabled {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSCE_ENABLE", "false")
		return c
	}

	if c.TimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSCE_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.TimeoutSeconds))
	}

	if c.SubBuffSize != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSCE_SUB_BUFF_SIZE", fmt.Sprintf("%d", *c.SubBuffSize))
	}

	if c.MaxSSEIdleSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSCE_MAX_SSE_IDLE_SECONDS", fmt.Sprintf("%d", *c.MaxSSEIdleSeconds))
	}

	if c.MaxSSEConnections != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSCE_MAX_SSE_CONNECTIONS", fmt.Sprintf("%d", *c.MaxSSEConnections))
	}

	return c
}
