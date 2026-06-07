package config

import (
	"fmt"
	"strings"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// AgentsConfig configures the kubemq-server Agents (A2A) platform connector.
// Maps to server Connectors.A2A. The connector is enabled by default server-side.
type AgentsConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	AgentTTLSeconds *int32 `json:"agentTtlSeconds,omitempty" yaml:"agentTtlSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	DefaultTimeoutSeconds *int32 `json:"defaultTimeoutSeconds,omitempty" yaml:"defaultTimeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxTimeoutSeconds *int32 `json:"maxTimeoutSeconds,omitempty" yaml:"maxTimeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxAgents *int32 `json:"maxAgents,omitempty" yaml:"maxAgents,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxSSEIdleSeconds *int32 `json:"maxSseIdleSeconds,omitempty" yaml:"maxSseIdleSeconds,omitempty"`

	// +optional
	TrustedOrigins []string `json:"trustedOrigins,omitempty" yaml:"trustedOrigins,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	AgentMaxResponseBytes *int64 `json:"agentMaxResponseBytes,omitempty" yaml:"agentMaxResponseBytes,omitempty"`

	// +optional
	AgentTLSSkipVerify bool `json:"agentTlsSkipVerify,omitempty" yaml:"agentTlsSkipVerify,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	AgentMaxConcurrency *int32 `json:"agentMaxConcurrency,omitempty" yaml:"agentMaxConcurrency,omitempty"`
}

func (c *AgentsConfig) DeepCopy() *AgentsConfig {
	out := &AgentsConfig{}

	out.Disabled = c.Disabled
	out.AgentTLSSkipVerify = c.AgentTLSSkipVerify

	if c.AgentTTLSeconds != nil {
		out.AgentTTLSeconds = new(int32)
		*out.AgentTTLSeconds = *c.AgentTTLSeconds
	}

	if c.DefaultTimeoutSeconds != nil {
		out.DefaultTimeoutSeconds = new(int32)
		*out.DefaultTimeoutSeconds = *c.DefaultTimeoutSeconds
	}

	if c.MaxTimeoutSeconds != nil {
		out.MaxTimeoutSeconds = new(int32)
		*out.MaxTimeoutSeconds = *c.MaxTimeoutSeconds
	}

	if c.MaxAgents != nil {
		out.MaxAgents = new(int32)
		*out.MaxAgents = *c.MaxAgents
	}

	if c.MaxSSEIdleSeconds != nil {
		out.MaxSSEIdleSeconds = new(int32)
		*out.MaxSSEIdleSeconds = *c.MaxSSEIdleSeconds
	}

	if c.TrustedOrigins != nil {
		out.TrustedOrigins = make([]string, len(c.TrustedOrigins))
		copy(out.TrustedOrigins, c.TrustedOrigins)
	}

	if c.AgentMaxResponseBytes != nil {
		out.AgentMaxResponseBytes = new(int64)
		*out.AgentMaxResponseBytes = *c.AgentMaxResponseBytes
	}

	if c.AgentMaxConcurrency != nil {
		out.AgentMaxConcurrency = new(int32)
		*out.AgentMaxConcurrency = *c.AgentMaxConcurrency
	}

	return out
}

func (c *AgentsConfig) SetConfig(config *deployment.Config) *AgentsConfig {
	if c.Disabled {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSA2_A_ENABLE", "false")
		return c
	}

	if c.AgentTTLSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSA2_A_AGENT_TTL_SECONDS", fmt.Sprintf("%d", *c.AgentTTLSeconds))
	}

	if c.DefaultTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSA2_A_DEFAULT_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.DefaultTimeoutSeconds))
	}

	if c.MaxTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSA2_A_MAX_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.MaxTimeoutSeconds))
	}

	if c.MaxAgents != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSA2_A_MAX_AGENTS", fmt.Sprintf("%d", *c.MaxAgents))
	}

	if c.MaxSSEIdleSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSA2_A_MAX_SSE_IDLE_SECONDS", fmt.Sprintf("%d", *c.MaxSSEIdleSeconds))
	}

	if len(c.TrustedOrigins) > 0 {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSA2_A_TRUSTED_ORIGINS", strings.Join(c.TrustedOrigins, ","))
	}

	if c.AgentMaxResponseBytes != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSA2_A_AGENT_MAX_RESPONSE_BYTES", fmt.Sprintf("%d", *c.AgentMaxResponseBytes))
	}

	if c.AgentTLSSkipVerify {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSA2_A_AGENT_TLS_SKIP_VERIFY", "true")
	}

	if c.AgentMaxConcurrency != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSA2_A_AGENT_MAX_CONCURRENCY", fmt.Sprintf("%d", *c.AgentMaxConcurrency))
	}

	return c
}
