package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1alpha1/kubemqcluster/deployment"
)

type AuthorizationConfig struct {
	// +optional
	Policy string `json:"policy,omitempty"`

	// +optional
	Url string `json:"url,omitempty"`

	// +optional
	AutoReload int32 `json:"autoReload,omitempty"`
}

func (c *AuthorizationConfig) SetConfig(config *deployment.Config) *AuthorizationConfig {
	if c.Policy == "" && c.Url == "" {
		return c
	}
	cmConfig, ok := config.ConfigMaps[config.Name]
	if ok {
		cmConfig.SetStringVariable("AUTHORIZATION_ENABLE", "true")
		if c.Policy != "" {
			cmConfig.SetDataVariable("AUTHORIZATION_POLICY_DATA", c.Policy)
		}
		if c.Url != "" {
			cmConfig.SetDataVariable("AUTHORIZATION_URL", c.Url)
		}

		if c.AutoReload != 0 {
			cmConfig.SetStringVariable("AUTHORIZATION_AUTO_RELOAD", fmt.Sprintf("%d", c.AutoReload))
		}

	}
	return c
}
