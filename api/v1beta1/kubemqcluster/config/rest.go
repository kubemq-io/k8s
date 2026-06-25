package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type RestConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// +optional
	Port *int32 `json:"port,omitempty" yaml:"port,omitempty"`

	// +optional
	// +kubebuilder:validation:Pattern=(ClusterIP|NodePort|LoadBalancer)
	Expose string `json:"expose,omitempty" yaml:"expose,omitempty"`

	// +optional
	BufferSize int32 `json:"bufferSize,omitempty" yaml:"bufferSize,omitempty"`

	// +optional
	BodyLimit int32 `json:"bodyLimit,omitempty" yaml:"bodyLimit,omitempty"`

	// +optional
	NodePort int32 `json:"nodePort,omitempty" yaml:"nodePort,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	ReadTimeout *int32 `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	WriteTimeout *int32 `json:"writeTimeout,omitempty" yaml:"writeTimeout,omitempty"`

	// +optional
	Cors *CorsConfig `json:"cors,omitempty" yaml:"cors,omitempty"`
}

func (c *RestConfig) DeepCopy() *RestConfig {
	out := &RestConfig{}

	out.Disabled = c.Disabled
	if c.Port != nil {
		out.Port = new(int32)
		*out.Port = *c.Port
	}
	out.Expose = c.Expose
	out.BufferSize = c.BufferSize
	out.BodyLimit = c.BodyLimit
	out.NodePort = c.NodePort

	if c.ReadTimeout != nil {
		out.ReadTimeout = new(int32)
		*out.ReadTimeout = *c.ReadTimeout
	}

	if c.WriteTimeout != nil {
		out.WriteTimeout = new(int32)
		*out.WriteTimeout = *c.WriteTimeout
	}

	if c.Cors != nil {
		out.Cors = c.Cors.DeepCopy()
	}

	return out
}

func (c *RestConfig) getDefaults() *RestConfig {
	if c.Expose == "" {
		c.Expose = "ClusterIP"
	}
	return c
}

func (c *RestConfig) SetConfig(config *deployment.Config) *RestConfig {
	c.getDefaults()
	if c.Disabled {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_REST_ENABLE", "false")
		return c
	}
	svc, ok := config.Services["rest"]
	if ok {
		if c.Port != nil {
			svc.SetContainerPort(*c.Port).SetTargetPort(*c.Port)
			config.SetConfigMapStringValues(config.Name, "CONNECTORS_REST_PORT", fmt.Sprintf("%d", *c.Port))
			config.StatefulSet.SetRestPort(*c.Port) // containerPort + prometheus annotation
		}

		if c.Expose == "NodePort" && c.NodePort > 0 {
			svc.SetNodePort(c.NodePort)
		} else {
			svc.SetNodePort(0)
		}
		svc.SetExpose(c.Expose)
	}

	if c.BufferSize != 0 {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_REST_SUB_BUFF_SIZE", fmt.Sprintf("%d", c.BufferSize))
	}

	if c.BodyLimit != 0 {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_REST_BODY_LIMIT", fmt.Sprintf("%d", c.BodyLimit))
	}

	if c.ReadTimeout != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_REST_READ_TIMEOUT", fmt.Sprintf("%d", *c.ReadTimeout))
	}

	if c.WriteTimeout != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_REST_WRITE_TIMEOUT", fmt.Sprintf("%d", *c.WriteTimeout))
	}

	if c.Cors != nil {
		c.Cors.setConfig(config, "CONNECTORS_REST_CORS")
	}

	return c
}
