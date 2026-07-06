package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type RoutingConfig struct {
	// +optional
	Data string `json:"data,omitempty" yaml:"data,omitempty"`

	// +optional
	Url string `json:"url,omitempty" yaml:"url,omitempty"`

	// +optional
	AutoReload int32 `json:"autoReload,omitempty" yaml:"autoReload,omitempty"`
}

func (c *RoutingConfig) SetConfig(config *deployment.Config) *RoutingConfig {
	if c.Data == "" && c.Url == "" {
		return c
	}

	cmConfig, ok := config.ConfigMaps[config.Name]
	if ok {
		cmConfig.SetStringVariable("ROUTING_ENABLE", "true")

		// The server binds Routing.Data / Routing.Url and consumes them RAW — it
		// does not base64-decode them (unlike AUTHENTICATION_CONFIG). Emit them as
		// plain string variables; SetDataVariable would base64-encode and the
		// server would then fail to parse the routing table / URL.
		if c.Data != "" {
			cmConfig.SetStringVariable("ROUTING_DATA", c.Data)
		}
		if c.Url != "" {
			cmConfig.SetStringVariable("ROUTING_URL", c.Url)
		}

		if c.AutoReload != 0 {
			cmConfig.SetStringVariable("ROUTING_AUTO_RELOAD", fmt.Sprintf("%d", c.AutoReload))
		}
	}
	return c
}
