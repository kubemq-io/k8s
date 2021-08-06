package config

import "github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"

type ApiConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// +optional
	Port int32 `json:"port,omitempty"`

	// +optional
	// +kubebuilder:validation:Pattern=(ClusterIP|NodePort|LoadBalancer)
	Expose string `json:"expose,omitempty"`

	// +optional
	NodePort int32 `json:"nodePort,omitempty"`
}

func (c *ApiConfig) getDefaults() *ApiConfig {
	if c.Port == 0 {
		c.Port = 8080
	}
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

		svc.SetTargetPort(8080).
			SetContainerPort(c.Port)

		if c.Expose == "NodePort" && c.NodePort > 0 {
			svc.SetNodePort(c.NodePort)
		} else {
			svc.SetNodePort(0)
		}
		svc.SetExpose(c.Expose)
	}
	return c
}
