package config

import (
	"fmt"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// StompConfig configures the kubemq-server STOMP connector.
// Maps to server Connectors.Stomp. The connector is enabled by default server-side.
type StompConfig struct {
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
	// +kubebuilder:validation:Enum=events;queues;store;none
	DefaultPattern *string `json:"defaultPattern,omitempty" yaml:"defaultPattern,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10000
	SubBuffSize *int32 `json:"subBuffSize,omitempty" yaml:"subBuffSize,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxConnections *int32 `json:"maxConnections,omitempty" yaml:"maxConnections,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxBodySize *int32 `json:"maxBodySize,omitempty" yaml:"maxBodySize,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	HeartbeatMs *int32 `json:"heartbeatMs,omitempty" yaml:"heartbeatMs,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	QueueAckTimeoutSeconds *int32 `json:"queueAckTimeoutSeconds,omitempty" yaml:"queueAckTimeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	RPCTimeoutSeconds *int32 `json:"rpcTimeoutSeconds,omitempty" yaml:"rpcTimeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	RPCMaxPending *int32 `json:"rpcMaxPending,omitempty" yaml:"rpcMaxPending,omitempty"`
}

func (c *StompConfig) DeepCopy() *StompConfig {
	out := &StompConfig{}

	out.Disabled = c.Disabled

	if c.Port != nil {
		out.Port = new(int32)
		*out.Port = *c.Port
	}

	if c.TLSPort != nil {
		out.TLSPort = new(int32)
		*out.TLSPort = *c.TLSPort
	}

	if c.DefaultPattern != nil {
		out.DefaultPattern = new(string)
		*out.DefaultPattern = *c.DefaultPattern
	}

	if c.SubBuffSize != nil {
		out.SubBuffSize = new(int32)
		*out.SubBuffSize = *c.SubBuffSize
	}

	if c.MaxConnections != nil {
		out.MaxConnections = new(int32)
		*out.MaxConnections = *c.MaxConnections
	}

	if c.MaxBodySize != nil {
		out.MaxBodySize = new(int32)
		*out.MaxBodySize = *c.MaxBodySize
	}

	if c.HeartbeatMs != nil {
		out.HeartbeatMs = new(int32)
		*out.HeartbeatMs = *c.HeartbeatMs
	}

	if c.QueueAckTimeoutSeconds != nil {
		out.QueueAckTimeoutSeconds = new(int32)
		*out.QueueAckTimeoutSeconds = *c.QueueAckTimeoutSeconds
	}

	if c.RPCTimeoutSeconds != nil {
		out.RPCTimeoutSeconds = new(int32)
		*out.RPCTimeoutSeconds = *c.RPCTimeoutSeconds
	}

	if c.RPCMaxPending != nil {
		out.RPCMaxPending = new(int32)
		*out.RPCMaxPending = *c.RPCMaxPending
	}

	return out
}

func (c *StompConfig) SetConfig(config *deployment.Config) *StompConfig {
	if c.Disabled {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_ENABLE", "false")
		return c
	}

	if c.Port != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_PORT", fmt.Sprintf("%d", *c.Port))
	}

	if c.TLSPort != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_TLS_PORT", fmt.Sprintf("%d", *c.TLSPort))
	}

	if c.DefaultPattern != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_DEFAULT_PATTERN", *c.DefaultPattern)
	}

	if c.SubBuffSize != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_SUB_BUFF_SIZE", fmt.Sprintf("%d", *c.SubBuffSize))
	}

	if c.MaxConnections != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_MAX_CONNECTIONS", fmt.Sprintf("%d", *c.MaxConnections))
	}

	if c.MaxBodySize != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_MAX_BODY_SIZE", fmt.Sprintf("%d", *c.MaxBodySize))
	}

	if c.HeartbeatMs != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_HEARTBEAT_MS", fmt.Sprintf("%d", *c.HeartbeatMs))
	}

	if c.QueueAckTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_QUEUE_ACK_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.QueueAckTimeoutSeconds))
	}

	if c.RPCTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_RPC_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.RPCTimeoutSeconds))
	}

	if c.RPCMaxPending != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_STOMP_RPC_MAX_PENDING", fmt.Sprintf("%d", *c.RPCMaxPending))
	}

	return c
}
