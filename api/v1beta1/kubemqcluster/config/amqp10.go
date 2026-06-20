package config

import (
	"fmt"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// Amqp10Config configures the kubemq-server AMQP 1.0 connector.
// Maps to server Connectors.Amqp10. The connector is enabled by default server-side.
type Amqp10Config struct {
	// +optional
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port *int32 `json:"port,omitempty" yaml:"port,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	TLSPort *int32 `json:"tlsPort,omitempty" yaml:"tlsPort,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=512
	MaxFrameSize *int32 `json:"maxFrameSize,omitempty" yaml:"maxFrameSize,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxMessageSize *int64 `json:"maxMessageSize,omitempty" yaml:"maxMessageSize,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	SessionMax *int32 `json:"sessionMax,omitempty" yaml:"sessionMax,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxLinksPerSession *int32 `json:"maxLinksPerSession,omitempty" yaml:"maxLinksPerSession,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxConnections *int32 `json:"maxConnections,omitempty" yaml:"maxConnections,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	IdleTimeoutSeconds *int32 `json:"idleTimeoutSeconds,omitempty" yaml:"idleTimeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=queues;events;events-store;commands;queries
	DefaultPattern *string `json:"defaultPattern,omitempty" yaml:"defaultPattern,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1024
	GetBatchSize *int32 `json:"getBatchSize,omitempty" yaml:"getBatchSize,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxUnsettledPerLink *int32 `json:"maxUnsettledPerLink,omitempty" yaml:"maxUnsettledPerLink,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	DefaultRPCTimeoutSeconds *int32 `json:"defaultRpcTimeoutSeconds,omitempty" yaml:"defaultRpcTimeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	RPCMaxPending *int32 `json:"rpcMaxPending,omitempty" yaml:"rpcMaxPending,omitempty"`
}

func (c *Amqp10Config) DeepCopy() *Amqp10Config {
	out := &Amqp10Config{}

	out.Disabled = c.Disabled

	if c.Port != nil {
		out.Port = new(int32)
		*out.Port = *c.Port
	}

	if c.TLSPort != nil {
		out.TLSPort = new(int32)
		*out.TLSPort = *c.TLSPort
	}

	if c.MaxFrameSize != nil {
		out.MaxFrameSize = new(int32)
		*out.MaxFrameSize = *c.MaxFrameSize
	}

	if c.MaxMessageSize != nil {
		out.MaxMessageSize = new(int64)
		*out.MaxMessageSize = *c.MaxMessageSize
	}

	if c.SessionMax != nil {
		out.SessionMax = new(int32)
		*out.SessionMax = *c.SessionMax
	}

	if c.MaxLinksPerSession != nil {
		out.MaxLinksPerSession = new(int32)
		*out.MaxLinksPerSession = *c.MaxLinksPerSession
	}

	if c.MaxConnections != nil {
		out.MaxConnections = new(int32)
		*out.MaxConnections = *c.MaxConnections
	}

	if c.IdleTimeoutSeconds != nil {
		out.IdleTimeoutSeconds = new(int32)
		*out.IdleTimeoutSeconds = *c.IdleTimeoutSeconds
	}

	if c.DefaultPattern != nil {
		out.DefaultPattern = new(string)
		*out.DefaultPattern = *c.DefaultPattern
	}

	if c.GetBatchSize != nil {
		out.GetBatchSize = new(int32)
		*out.GetBatchSize = *c.GetBatchSize
	}

	if c.MaxUnsettledPerLink != nil {
		out.MaxUnsettledPerLink = new(int32)
		*out.MaxUnsettledPerLink = *c.MaxUnsettledPerLink
	}

	if c.DefaultRPCTimeoutSeconds != nil {
		out.DefaultRPCTimeoutSeconds = new(int32)
		*out.DefaultRPCTimeoutSeconds = *c.DefaultRPCTimeoutSeconds
	}

	if c.RPCMaxPending != nil {
		out.RPCMaxPending = new(int32)
		*out.RPCMaxPending = *c.RPCMaxPending
	}

	return out
}

func (c *Amqp10Config) SetConfig(config *deployment.Config) *Amqp10Config {
	if c.Disabled {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_ENABLE", "false")
		return c
	}

	if c.Port != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_PORT", fmt.Sprintf("%d", *c.Port))
	}

	if c.TLSPort != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_TLS_PORT", fmt.Sprintf("%d", *c.TLSPort))
	}

	if c.MaxFrameSize != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_MAX_FRAME_SIZE", fmt.Sprintf("%d", *c.MaxFrameSize))
	}

	if c.MaxMessageSize != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_MAX_MESSAGE_SIZE", fmt.Sprintf("%d", *c.MaxMessageSize))
	}

	if c.SessionMax != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_SESSION_MAX", fmt.Sprintf("%d", *c.SessionMax))
	}

	if c.MaxLinksPerSession != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_MAX_LINKS_PER_SESSION", fmt.Sprintf("%d", *c.MaxLinksPerSession))
	}

	if c.MaxConnections != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_MAX_CONNECTIONS", fmt.Sprintf("%d", *c.MaxConnections))
	}

	if c.IdleTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_IDLE_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.IdleTimeoutSeconds))
	}

	if c.DefaultPattern != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_DEFAULT_PATTERN", *c.DefaultPattern)
	}

	if c.GetBatchSize != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_GET_BATCH_SIZE", fmt.Sprintf("%d", *c.GetBatchSize))
	}

	if c.MaxUnsettledPerLink != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_MAX_UNSETTLED_PER_LINK", fmt.Sprintf("%d", *c.MaxUnsettledPerLink))
	}

	if c.DefaultRPCTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_DEFAULT_RPC_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.DefaultRPCTimeoutSeconds))
	}

	if c.RPCMaxPending != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP10_RPC_MAX_PENDING", fmt.Sprintf("%d", *c.RPCMaxPending))
	}

	return c
}
