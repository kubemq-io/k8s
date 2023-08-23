package config

import "github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"

type NotificationConfig struct {
	// +optional
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	// +optional
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty"`

	// +optional
	Log bool `json:"log,omitempty" yaml:"log,omitempty"`
}

func (c *NotificationConfig) SetConfig(config *deployment.Config) *NotificationConfig {

	if c.Enabled {
		config.SetConfigMapStringValues(config.Name, "NOTIFICATION_ENABLE", "true")
	}

	if c.Prefix != "" {
		config.SetConfigMapStringValues(config.Name, "NOTIFICATION_REPORT_CHANNEL_PREFIX", c.Prefix)
	}

	if c.Log {
		config.SetConfigMapStringValues(config.Name, "NOTIFICATION_LOG", "true")
	}
	return c
}
