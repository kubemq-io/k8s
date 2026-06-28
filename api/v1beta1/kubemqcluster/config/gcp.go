package config

import (
	"fmt"
	"strconv"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// GcpConfig configures the kubemq-server Google Cloud Pub/Sub wire-protocol
// connector. Maps to server Connectors.PubSub. The connector is opt-in
// (disabled by default); set enabled: true to activate it (opens port 8085).
type GcpConfig struct {
	// +optional
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port *int32 `json:"port,omitempty" yaml:"port,omitempty"`

	// +optional
	AdvertisedEndpoint *string `json:"advertisedEndpoint,omitempty" yaml:"advertisedEndpoint,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxMessageBytes *int32 `json:"maxMessageBytes,omitempty" yaml:"maxMessageBytes,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=10
	// +kubebuilder:validation:Maximum=600
	DefaultAckDeadlineSeconds *int32 `json:"defaultAckDeadlineSeconds,omitempty" yaml:"defaultAckDeadlineSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxOutstandingMessages *int32 `json:"maxOutstandingMessages,omitempty" yaml:"maxOutstandingMessages,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxInflightPerSubscription *int32 `json:"maxInflightPerSubscription,omitempty" yaml:"maxInflightPerSubscription,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxConcurrentPolls *int32 `json:"maxConcurrentPolls,omitempty" yaml:"maxConcurrentPolls,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=256
	DeliveryShards *int32 `json:"deliveryShards,omitempty" yaml:"deliveryShards,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=3600
	MaxAckExtensionSeconds *int32 `json:"maxAckExtensionSeconds,omitempty" yaml:"maxAckExtensionSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	StreamCloseSeconds *int32 `json:"streamCloseSeconds,omitempty" yaml:"streamCloseSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxSeekReplay *int32 `json:"maxSeekReplay,omitempty" yaml:"maxSeekReplay,omitempty"`

	// +optional
	EnableReflection *bool `json:"enableReflection,omitempty" yaml:"enableReflection,omitempty"`
}

func (c *GcpConfig) DeepCopy() *GcpConfig {
	out := &GcpConfig{}

	if c.Enabled != nil {
		out.Enabled = new(bool)
		*out.Enabled = *c.Enabled
	}

	if c.Port != nil {
		out.Port = new(int32)
		*out.Port = *c.Port
	}

	if c.AdvertisedEndpoint != nil {
		out.AdvertisedEndpoint = new(string)
		*out.AdvertisedEndpoint = *c.AdvertisedEndpoint
	}

	if c.MaxMessageBytes != nil {
		out.MaxMessageBytes = new(int32)
		*out.MaxMessageBytes = *c.MaxMessageBytes
	}

	if c.DefaultAckDeadlineSeconds != nil {
		out.DefaultAckDeadlineSeconds = new(int32)
		*out.DefaultAckDeadlineSeconds = *c.DefaultAckDeadlineSeconds
	}

	if c.MaxOutstandingMessages != nil {
		out.MaxOutstandingMessages = new(int32)
		*out.MaxOutstandingMessages = *c.MaxOutstandingMessages
	}

	if c.MaxInflightPerSubscription != nil {
		out.MaxInflightPerSubscription = new(int32)
		*out.MaxInflightPerSubscription = *c.MaxInflightPerSubscription
	}

	if c.MaxConcurrentPolls != nil {
		out.MaxConcurrentPolls = new(int32)
		*out.MaxConcurrentPolls = *c.MaxConcurrentPolls
	}

	if c.DeliveryShards != nil {
		out.DeliveryShards = new(int32)
		*out.DeliveryShards = *c.DeliveryShards
	}

	if c.MaxAckExtensionSeconds != nil {
		out.MaxAckExtensionSeconds = new(int32)
		*out.MaxAckExtensionSeconds = *c.MaxAckExtensionSeconds
	}

	if c.StreamCloseSeconds != nil {
		out.StreamCloseSeconds = new(int32)
		*out.StreamCloseSeconds = *c.StreamCloseSeconds
	}

	if c.MaxSeekReplay != nil {
		out.MaxSeekReplay = new(int32)
		*out.MaxSeekReplay = *c.MaxSeekReplay
	}

	if c.EnableReflection != nil {
		out.EnableReflection = new(bool)
		*out.EnableReflection = *c.EnableReflection
	}

	return out
}

func (c *GcpConfig) SetConfig(config *deployment.Config) *GcpConfig {
	effective := c.Enabled != nil && *c.Enabled
	config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_ENABLE", strconv.FormatBool(effective))
	if !effective {
		return c
	}

	// Reflect a custom port onto the K8s Service so traffic reaches the listener.
	if svc, ok := config.Services["gcp"]; ok {
		if c.Port != nil {
			svc.SetPort("gcp-grpc", *c.Port)
		}
	}

	if c.Port != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_PORT", fmt.Sprintf("%d", *c.Port))
	}

	if c.AdvertisedEndpoint != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_ADVERTISED_ENDPOINT", *c.AdvertisedEndpoint)
	}

	if c.MaxMessageBytes != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_MAX_MESSAGE_BYTES", fmt.Sprintf("%d", *c.MaxMessageBytes))
	}

	if c.DefaultAckDeadlineSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_DEFAULT_ACK_DEADLINE_SECONDS", fmt.Sprintf("%d", *c.DefaultAckDeadlineSeconds))
	}

	if c.MaxOutstandingMessages != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_MAX_OUTSTANDING_MESSAGES", fmt.Sprintf("%d", *c.MaxOutstandingMessages))
	}

	if c.MaxInflightPerSubscription != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_MAX_INFLIGHT_PER_SUBSCRIPTION", fmt.Sprintf("%d", *c.MaxInflightPerSubscription))
	}

	if c.MaxConcurrentPolls != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_MAX_CONCURRENT_POLLS", fmt.Sprintf("%d", *c.MaxConcurrentPolls))
	}

	if c.DeliveryShards != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_DELIVERY_SHARDS", fmt.Sprintf("%d", *c.DeliveryShards))
	}

	if c.MaxAckExtensionSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_MAX_ACK_EXTENSION_SECONDS", fmt.Sprintf("%d", *c.MaxAckExtensionSeconds))
	}

	if c.StreamCloseSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_STREAM_CLOSE_SECONDS", fmt.Sprintf("%d", *c.StreamCloseSeconds))
	}

	if c.MaxSeekReplay != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_MAX_SEEK_REPLAY", fmt.Sprintf("%d", *c.MaxSeekReplay))
	}

	if c.EnableReflection != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GCP_ENABLE_REFLECTION", strconv.FormatBool(*c.EnableReflection))
	}

	return c
}
