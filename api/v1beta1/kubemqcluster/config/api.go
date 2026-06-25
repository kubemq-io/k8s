package config

import (
	"fmt"
	"strings"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type ApiConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// +optional
	Port *int32 `json:"port,omitempty" yaml:"port,omitempty"`

	// +optional
	// +kubebuilder:validation:Pattern=(ClusterIP|NodePort|LoadBalancer)
	Expose string `json:"expose,omitempty" yaml:"expose,omitempty"`

	// +optional
	NodePort int32 `json:"nodePort,omitempty" yaml:"nodePort,omitempty"`

	// +optional
	AllowOrigins []string `json:"allowOrigins,omitempty" yaml:"allowOrigins,omitempty"`
}

func (c *ApiConfig) DeepCopy() *ApiConfig {
	out := &ApiConfig{}

	out.Disabled = c.Disabled
	if c.Port != nil {
		out.Port = new(int32)
		*out.Port = *c.Port
	}
	out.Expose = c.Expose
	out.NodePort = c.NodePort

	if c.AllowOrigins != nil {
		out.AllowOrigins = make([]string, len(c.AllowOrigins))
		copy(out.AllowOrigins, c.AllowOrigins)
	}

	return out
}

func (c *ApiConfig) getDefaults() *ApiConfig {
	if c.Expose == "" {
		c.Expose = "ClusterIP"
	}
	return c
}
func (c *ApiConfig) SetConfig(config *deployment.Config) *ApiConfig {
	c.getDefaults()
	if c.Disabled {
		config.SetConfigMapStringValues(config.Name, "API_ENABLE", "false")
		return c
	}

	svc, ok := config.Services["api"]
	if ok {

		if c.Port != nil {
			svc.SetContainerPort(*c.Port).SetTargetPort(*c.Port)
			config.SetConfigMapStringValues(config.Name, "API_PORT", fmt.Sprintf("%d", *c.Port))
			config.StatefulSet.SetApiPort(*c.Port) // containerPort + prometheus annotation
		}

		if c.Expose == "NodePort" && c.NodePort > 0 {
			svc.SetNodePort(c.NodePort)
		} else {
			svc.SetNodePort(0)
		}
		svc.SetExpose(c.Expose)
	}

	if len(c.AllowOrigins) > 0 {
		config.SetConfigMapStringValues(config.Name, "API_ALLOW_ORIGINS", strings.Join(c.AllowOrigins, ","))
	}

	return c
}
