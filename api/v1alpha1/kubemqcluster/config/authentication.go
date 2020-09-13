package config

import "github.com/kubemq-io/k8s/api/v1alpha1/kubemqcluster/deployment"

type AuthenticationConfig struct {
	// +optional
	Key string `json:"key,omitempty"`

	// +optional
	Type string `json:"type,omitempty"`
}

func (c *AuthenticationConfig) SetConfig(config *deployment.Config) *AuthenticationConfig {
	if c.Key == "" || c.Type == "" {
		return c
	}
	secConfig, ok := config.Secrets[config.Name]
	if ok {
		secConfig.SetStringVariable("AUTHENTICATION_ENABLE", "true").
			SetDataVariable("AUTHENTICATION_JWT_CONFIG_KEY", c.Key).
			SetStringVariable("AUTHENTICATION_JWT_CONFIG_SIGNATURE_TYPE", c.Type)
	}
	return c
}
