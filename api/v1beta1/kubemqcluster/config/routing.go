package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type RoutingConfig struct {
	// +optional
	Data string `json:"data,omitempty"`

	// +optional
	Url string `json:"url,omitempty"`

	// +optional
	AutoReload int32 `json:"autoReload,omitempty"`
}

func (c *RoutingConfig) SetConfig(config *deployment.Config) *RoutingConfig {
	if c.Data == "" && c.Url == "" {
		return c
	}

	cmConfig, ok := config.ConfigMaps[config.Name]
	if ok {
		cmConfig.SetStringVariable("ROUTING_ENABLE", "true")

		if c.Data != "" {
			cmConfig.SetDataVariable("ROUTING_DATA", c.Data)
		}
		if c.Url != "" {
			cmConfig.SetDataVariable("ROUTING_URL", c.Url)
		}

		if c.AutoReload != 0 {
			cmConfig.SetStringVariable("ROUTING_AUTO_RELOAD", fmt.Sprintf("%d", c.AutoReload))
		}
	}
	return c
}
