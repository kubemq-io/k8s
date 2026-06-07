package config

import (
	"fmt"
	"strings"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// McpConfig configures the kubemq-server MCP (Model Context Protocol) connector.
// Maps to server Connectors.MCP. The connector is enabled by default server-side.
type McpConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	ToolTimeoutSeconds *int32 `json:"toolTimeoutSeconds,omitempty" yaml:"toolTimeoutSeconds,omitempty"`

	// +optional
	TrustedOrigins []string `json:"trustedOrigins,omitempty" yaml:"trustedOrigins,omitempty"`
}

func (c *McpConfig) DeepCopy() *McpConfig {
	out := &McpConfig{}

	out.Disabled = c.Disabled

	if c.ToolTimeoutSeconds != nil {
		out.ToolTimeoutSeconds = new(int32)
		*out.ToolTimeoutSeconds = *c.ToolTimeoutSeconds
	}

	if c.TrustedOrigins != nil {
		out.TrustedOrigins = make([]string, len(c.TrustedOrigins))
		copy(out.TrustedOrigins, c.TrustedOrigins)
	}

	return out
}

func (c *McpConfig) SetConfig(config *deployment.Config) *McpConfig {
	if c.Disabled {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMCP_ENABLE", "false")
		return c
	}

	if c.ToolTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMCP_TOOL_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.ToolTimeoutSeconds))
	}

	if len(c.TrustedOrigins) > 0 {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMCP_TRUSTED_ORIGINS", strings.Join(c.TrustedOrigins, ","))
	}

	return c
}
