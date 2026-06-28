package config

import (
	"fmt"
	"strconv"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// AwsConfig configures the kubemq-server AWS (SQS/SNS) connector.
// Maps to server Connectors.Aws. The connector is opt-in (disabled by default);
// set enabled: true to activate it (opens port 4566).
// Credentials are intentionally omitted; they are supplied through the cluster
// Secret / server config, not the CR.
type AwsConfig struct {
	// +optional
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port *int32 `json:"port,omitempty" yaml:"port,omitempty"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	Region *string `json:"region,omitempty" yaml:"region,omitempty"`

	// +optional
	// +kubebuilder:validation:Pattern=^[0-9]{12}$
	AccountID *string `json:"accountId,omitempty" yaml:"accountId,omitempty"`

	// +optional
	AdvertisedURL *string `json:"advertisedUrl,omitempty" yaml:"advertisedUrl,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxInflightPerQueue *int32 `json:"maxInflightPerQueue,omitempty" yaml:"maxInflightPerQueue,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxConcurrentPolls *int32 `json:"maxConcurrentPolls,omitempty" yaml:"maxConcurrentPolls,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	ReadTimeout *int32 `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	BodyLimit *string `json:"bodyLimit,omitempty" yaml:"bodyLimit,omitempty"`

	// +optional
	MessageSigning *bool `json:"messageSigning,omitempty" yaml:"messageSigning,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	SigningCertTtlHours *int32 `json:"signingCertTtlHours,omitempty" yaml:"signingCertTtlHours,omitempty"`

	// +optional
	CredentialsData *string `json:"credentialsData,omitempty" yaml:"credentialsData,omitempty"`
}

func (c *AwsConfig) DeepCopy() *AwsConfig {
	out := &AwsConfig{}

	if c.Enabled != nil {
		out.Enabled = new(bool)
		*out.Enabled = *c.Enabled
	}

	if c.Port != nil {
		out.Port = new(int32)
		*out.Port = *c.Port
	}

	if c.Region != nil {
		out.Region = new(string)
		*out.Region = *c.Region
	}

	if c.AccountID != nil {
		out.AccountID = new(string)
		*out.AccountID = *c.AccountID
	}

	if c.AdvertisedURL != nil {
		out.AdvertisedURL = new(string)
		*out.AdvertisedURL = *c.AdvertisedURL
	}

	if c.MaxInflightPerQueue != nil {
		out.MaxInflightPerQueue = new(int32)
		*out.MaxInflightPerQueue = *c.MaxInflightPerQueue
	}

	if c.MaxConcurrentPolls != nil {
		out.MaxConcurrentPolls = new(int32)
		*out.MaxConcurrentPolls = *c.MaxConcurrentPolls
	}

	if c.ReadTimeout != nil {
		out.ReadTimeout = new(int32)
		*out.ReadTimeout = *c.ReadTimeout
	}

	if c.BodyLimit != nil {
		out.BodyLimit = new(string)
		*out.BodyLimit = *c.BodyLimit
	}

	if c.MessageSigning != nil {
		out.MessageSigning = new(bool)
		*out.MessageSigning = *c.MessageSigning
	}

	if c.SigningCertTtlHours != nil {
		out.SigningCertTtlHours = new(int32)
		*out.SigningCertTtlHours = *c.SigningCertTtlHours
	}

	if c.CredentialsData != nil {
		out.CredentialsData = new(string)
		*out.CredentialsData = *c.CredentialsData
	}

	return out
}

func (c *AwsConfig) SetConfig(config *deployment.Config) *AwsConfig {
	effective := c.Enabled != nil && *c.Enabled
	config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_ENABLE", strconv.FormatBool(effective))
	if !effective {
		return c
	}

	// Reflect a custom port onto the K8s Service so traffic reaches the listener.
	if svc, ok := config.Services["aws"]; ok {
		if c.Port != nil {
			svc.SetPort("aws-http", *c.Port)
		}
	}

	if c.Port != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_PORT", fmt.Sprintf("%d", *c.Port))
	}

	if c.Region != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_REGION", *c.Region)
	}

	if c.AccountID != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_ACCOUNT_ID", *c.AccountID)
	}

	if c.AdvertisedURL != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_ADVERTISED_URL", *c.AdvertisedURL)
	}

	if c.MaxInflightPerQueue != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_MAX_INFLIGHT_PER_QUEUE", fmt.Sprintf("%d", *c.MaxInflightPerQueue))
	}

	if c.MaxConcurrentPolls != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_MAX_CONCURRENT_POLLS", fmt.Sprintf("%d", *c.MaxConcurrentPolls))
	}

	if c.ReadTimeout != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_READ_TIMEOUT", fmt.Sprintf("%d", *c.ReadTimeout))
	}

	if c.BodyLimit != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_BODY_LIMIT", *c.BodyLimit)
	}

	if c.MessageSigning != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_MESSAGE_SIGNING", strconv.FormatBool(*c.MessageSigning))
	}

	if c.SigningCertTtlHours != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_AWS_SIGNING_CERT_TTL_HOURS", fmt.Sprintf("%d", *c.SigningCertTtlHours))
	}

	if c.CredentialsData != nil {
		config.SetSecretDataValues(config.Name, "CONNECTORS_AWS_CREDENTIALS_DATA", *c.CredentialsData)
	}

	return c
}
