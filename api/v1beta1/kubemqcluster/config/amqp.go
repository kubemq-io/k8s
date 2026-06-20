package config

import (
	"fmt"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// AmqpConfig configures the kubemq-server AMQP 0.9.1 connector.
// Maps to server Connectors.Amqp. The connector is enabled by default server-side.
type AmqpConfig struct {
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
	// +kubebuilder:validation:Minimum=0
	HeartbeatSeconds *int32 `json:"heartbeatSeconds,omitempty" yaml:"heartbeatSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=4096
	FrameMax *int32 `json:"frameMax,omitempty" yaml:"frameMax,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	ChannelMax *int32 `json:"channelMax,omitempty" yaml:"channelMax,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxConnections *int32 `json:"maxConnections,omitempty" yaml:"maxConnections,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxBodySize *int32 `json:"maxBodySize,omitempty" yaml:"maxBodySize,omitempty"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	DefaultVhost *string `json:"defaultVhost,omitempty" yaml:"defaultVhost,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1024
	GetBatchSize *int32 `json:"getBatchSize,omitempty" yaml:"getBatchSize,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	DeadLetterMaxHops *int32 `json:"deadLetterMaxHops,omitempty" yaml:"deadLetterMaxHops,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxReceiveCount *int32 `json:"maxReceiveCount,omitempty" yaml:"maxReceiveCount,omitempty"`
}

func (c *AmqpConfig) DeepCopy() *AmqpConfig {
	out := &AmqpConfig{}

	out.Disabled = c.Disabled

	if c.Port != nil {
		out.Port = new(int32)
		*out.Port = *c.Port
	}

	if c.TLSPort != nil {
		out.TLSPort = new(int32)
		*out.TLSPort = *c.TLSPort
	}

	if c.HeartbeatSeconds != nil {
		out.HeartbeatSeconds = new(int32)
		*out.HeartbeatSeconds = *c.HeartbeatSeconds
	}

	if c.FrameMax != nil {
		out.FrameMax = new(int32)
		*out.FrameMax = *c.FrameMax
	}

	if c.ChannelMax != nil {
		out.ChannelMax = new(int32)
		*out.ChannelMax = *c.ChannelMax
	}

	if c.MaxConnections != nil {
		out.MaxConnections = new(int32)
		*out.MaxConnections = *c.MaxConnections
	}

	if c.MaxBodySize != nil {
		out.MaxBodySize = new(int32)
		*out.MaxBodySize = *c.MaxBodySize
	}

	if c.DefaultVhost != nil {
		out.DefaultVhost = new(string)
		*out.DefaultVhost = *c.DefaultVhost
	}

	if c.GetBatchSize != nil {
		out.GetBatchSize = new(int32)
		*out.GetBatchSize = *c.GetBatchSize
	}

	if c.DeadLetterMaxHops != nil {
		out.DeadLetterMaxHops = new(int32)
		*out.DeadLetterMaxHops = *c.DeadLetterMaxHops
	}

	if c.MaxReceiveCount != nil {
		out.MaxReceiveCount = new(int32)
		*out.MaxReceiveCount = *c.MaxReceiveCount
	}

	return out
}

func (c *AmqpConfig) SetConfig(config *deployment.Config) *AmqpConfig {
	if c.Disabled {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_ENABLE", "false")
		return c
	}

	if c.Port != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_PORT", fmt.Sprintf("%d", *c.Port))
	}

	if c.TLSPort != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_TLS_PORT", fmt.Sprintf("%d", *c.TLSPort))
	}

	if c.HeartbeatSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_HEARTBEAT_SECONDS", fmt.Sprintf("%d", *c.HeartbeatSeconds))
	}

	if c.FrameMax != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_FRAME_MAX", fmt.Sprintf("%d", *c.FrameMax))
	}

	if c.ChannelMax != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_CHANNEL_MAX", fmt.Sprintf("%d", *c.ChannelMax))
	}

	if c.MaxConnections != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_MAX_CONNECTIONS", fmt.Sprintf("%d", *c.MaxConnections))
	}

	if c.MaxBodySize != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_MAX_BODY_SIZE", fmt.Sprintf("%d", *c.MaxBodySize))
	}

	if c.DefaultVhost != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_DEFAULT_VHOST", *c.DefaultVhost)
	}

	if c.GetBatchSize != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_GET_BATCH_SIZE", fmt.Sprintf("%d", *c.GetBatchSize))
	}

	if c.DeadLetterMaxHops != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_DEAD_LETTER_MAX_HOPS", fmt.Sprintf("%d", *c.DeadLetterMaxHops))
	}

	if c.MaxReceiveCount != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AMQP_MAX_RECEIVE_COUNT", fmt.Sprintf("%d", *c.MaxReceiveCount))
	}

	return c
}
