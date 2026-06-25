package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type QueueConfig struct {

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxReceiveMessagesRequest *int32 `json:"maxReceiveMessagesRequest,omitempty" yaml:"maxReceiveMessagesRequest,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxWaitTimeoutSeconds *int32 `json:"maxWaitTimeoutSeconds,omitempty" yaml:"maxWaitTimeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxExpirationSeconds *int32 `json:"maxExpirationSeconds,omitempty" yaml:"maxExpirationSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxDelaySeconds *int32 `json:"maxDelaySeconds,omitempty" yaml:"maxDelaySeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxReQueues *int32 `json:"maxReQueues,omitempty" yaml:"maxReQueues,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxVisibilitySeconds *int32 `json:"maxVisibilitySeconds,omitempty" yaml:"maxVisibilitySeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	DefaultVisibilitySeconds *int32 `json:"defaultVisibilitySeconds,omitempty" yaml:"defaultVisibilitySeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	DefaultWaitTimeoutSeconds *int32 `json:"defaultWaitTimeoutSeconds,omitempty" yaml:"defaultWaitTimeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxInflight *int32 `json:"maxInflight,omitempty" yaml:"maxInflight,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	PubAckWaitSeconds *int32 `json:"pubAckWaitSeconds,omitempty" yaml:"pubAckWaitSeconds,omitempty"`
}

func (c *QueueConfig) DeepCopy() *QueueConfig {
	out := &QueueConfig{}

	if c.MaxReceiveMessagesRequest != nil {
		out.MaxReceiveMessagesRequest = new(int32)
		*out.MaxReceiveMessagesRequest = *c.MaxReceiveMessagesRequest
	}

	if c.MaxWaitTimeoutSeconds != nil {
		out.MaxWaitTimeoutSeconds = new(int32)
		*out.MaxWaitTimeoutSeconds = *c.MaxWaitTimeoutSeconds
	}

	if c.MaxExpirationSeconds != nil {
		out.MaxExpirationSeconds = new(int32)
		*out.MaxExpirationSeconds = *c.MaxExpirationSeconds
	}

	if c.MaxDelaySeconds != nil {
		out.MaxDelaySeconds = new(int32)
		*out.MaxDelaySeconds = *c.MaxDelaySeconds
	}

	if c.MaxReQueues != nil {
		out.MaxReQueues = new(int32)
		*out.MaxReQueues = *c.MaxReQueues
	}

	if c.MaxVisibilitySeconds != nil {
		out.MaxVisibilitySeconds = new(int32)
		*out.MaxVisibilitySeconds = *c.MaxVisibilitySeconds
	}

	if c.DefaultVisibilitySeconds != nil {
		out.DefaultVisibilitySeconds = new(int32)
		*out.DefaultVisibilitySeconds = *c.DefaultVisibilitySeconds
	}

	if c.DefaultWaitTimeoutSeconds != nil {
		out.DefaultWaitTimeoutSeconds = new(int32)
		*out.DefaultWaitTimeoutSeconds = *c.DefaultWaitTimeoutSeconds
	}

	if c.MaxInflight != nil {
		out.MaxInflight = new(int32)
		*out.MaxInflight = *c.MaxInflight
	}

	if c.PubAckWaitSeconds != nil {
		out.PubAckWaitSeconds = new(int32)
		*out.PubAckWaitSeconds = *c.PubAckWaitSeconds
	}

	return out
}

func (c *QueueConfig) SetConfig(config *deployment.Config) *QueueConfig {

	if c.MaxReceiveMessagesRequest != nil {
		config.SetConfigMapStringValues(config.Name, "QUEUE_MAX_NUMBER_OF_MESSAGES", fmt.Sprintf("%d", *c.MaxReceiveMessagesRequest))
	}

	if c.MaxWaitTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "QUEUE_MAX_WAIT_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.MaxWaitTimeoutSeconds))
	}

	if c.MaxExpirationSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "QUEUE_MAX_EXPIRATION_SECONDS", fmt.Sprintf("%d", *c.MaxExpirationSeconds))
	}

	if c.MaxDelaySeconds != nil {
		config.SetConfigMapStringValues(config.Name, "QUEUE_MAX_DELAY_SECONDS", fmt.Sprintf("%d", *c.MaxDelaySeconds))
	}

	if c.MaxReQueues != nil {
		config.SetConfigMapStringValues(config.Name, "QUEUE_MAX_RECEIVE_COUNT", fmt.Sprintf("%d", *c.MaxReQueues))
	}

	if c.MaxVisibilitySeconds != nil {
		config.SetConfigMapStringValues(config.Name, "QUEUE_MAX_VISIBILITY_SECONDS", fmt.Sprintf("%d", *c.MaxVisibilitySeconds))
	}

	if c.DefaultVisibilitySeconds != nil {
		config.SetConfigMapStringValues(config.Name, "QUEUE_DEFAULT_VISIBILITY_SECONDS", fmt.Sprintf("%d", *c.DefaultVisibilitySeconds))
	}

	if c.DefaultWaitTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "QUEUE_DEFAULT_WAIT_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.DefaultWaitTimeoutSeconds))
	}

	if c.MaxInflight != nil {
		config.SetConfigMapStringValues(config.Name, "QUEUE_MAX_INFLIGHT", fmt.Sprintf("%d", *c.MaxInflight))
	}

	if c.PubAckWaitSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "QUEUE_PUB_ACK_WAIT_SECONDS", fmt.Sprintf("%d", *c.PubAckWaitSeconds))
	}

	return c
}
