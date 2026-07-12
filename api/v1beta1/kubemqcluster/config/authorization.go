package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type AuthorizationConfig struct {
	// +optional
	Policy string `json:"policy,omitempty" yaml:"policy,omitempty"`

	// +optional
	Url string `json:"url,omitempty" yaml:"url,omitempty"`

	// +optional
	AutoReload int32 `json:"autoReload,omitempty" yaml:"autoReload,omitempty"`
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
			// Raw, NOT base64: the server reads Authorization.Url verbatim and hands it to
			// http.Get (validateURL rejects a base64 blob — "no scheme"). Policy DATA below
			// stays SetDataVariable because the server base64-decodes policy content.
			cmConfig.SetStringVariable("AUTHORIZATION_URL", c.Url)
		}

		if c.AutoReload != 0 {
			cmConfig.SetStringVariable("AUTHORIZATION_AUTO_RELOAD", fmt.Sprintf("%d", c.AutoReload))
		}

	}
	return c
}
