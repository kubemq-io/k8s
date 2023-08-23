package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type StoreConfig struct {
	// +optional
	Clean bool `json:"clean,omitempty" yaml:"clean,omitempty"`

	// +optional
	Path string `json:"path,omitempty" yaml:"path,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxChannels *int32 `json:"maxChannels,omitempty" yaml:"maxChannels,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxSubscribers *int32 `json:"maxSubscribers,omitempty" yaml:"maxSubscribers,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxMessages *int32 `json:"maxMessages,omitempty" yaml:"maxMessages,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxChannelSize *int32 `json:"maxChannelSize,omitempty" yaml:"maxChannelSize,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MessagesRetentionMinutes *int32 `json:"messagesRetentionMinutes,omitempty" yaml:"messagesRetentionMinutes,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	PurgeInactiveMinutes *int32 `json:"purgeInactiveMinutes,omitempty" yaml:"purgeInactiveMinutes,omitempty"`
}

func (c *StoreConfig) DeepCopy() *StoreConfig {
	out := &StoreConfig{}

	out.Clean = c.Clean
	out.Path = c.Path

	if c.MaxChannels != nil {
		out.MaxChannels = new(int32)
		*out.MaxChannels = *c.MaxChannels
	}

	if c.MaxSubscribers != nil {
		out.MaxSubscribers = new(int32)
		*out.MaxSubscribers = *c.MaxSubscribers

	}

	if c.MaxMessages != nil {
		out.MaxMessages = new(int32)
		*out.MaxMessages = *c.MaxMessages
	}

	if c.MaxChannelSize != nil {
		out.MaxChannelSize = new(int32)
		*out.MaxChannelSize = *c.MaxChannelSize
	}

	if c.MessagesRetentionMinutes != nil {
		out.MessagesRetentionMinutes = new(int32)
		*out.MessagesRetentionMinutes = *c.MessagesRetentionMinutes
	}

	if c.PurgeInactiveMinutes != nil {
		out.PurgeInactiveMinutes = new(int32)
		*out.PurgeInactiveMinutes = *c.PurgeInactiveMinutes
	}

	return out
}
func (c *StoreConfig) SetConfig(config *deployment.Config) *StoreConfig {
	if c.Clean {
		config.SetConfigMapStringValues(config.Name, "STORE_CLEAN_STORE", "true")
	}

	if c.Path != "" {
		config.SetConfigMapStringValues(config.Name, "STORE_STORE_PATH", c.Path)
	}

	if c.MaxChannels != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_QUEUES", fmt.Sprintf("%d", *c.MaxChannels))
	}

	if c.MaxSubscribers != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_SUBSCRIBERS", fmt.Sprintf("%d", *c.MaxSubscribers))
	}

	if c.MaxMessages != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_MESSAGES", fmt.Sprintf("%d", *c.MaxMessages))
	}

	if c.MaxChannelSize != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_QUEUE_SIZE", fmt.Sprintf("%d", *c.MaxChannelSize))
	}

	if c.MessagesRetentionMinutes != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_RETENTION", fmt.Sprintf("%d", *c.MessagesRetentionMinutes))
	}

	if c.PurgeInactiveMinutes != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_PURGE_INACTIVE", fmt.Sprintf("%d", *c.PurgeInactiveMinutes))
	}

	return c
}
