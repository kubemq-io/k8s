package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type GrpcConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// +optional
	Port int32 `json:"port,omitempty" yaml:"port,omitempty"`

	// +optional
	// +kubebuilder:validation:Pattern=(ClusterIP|NodePort|LoadBalancer)
	Expose string `json:"expose,omitempty" yaml:"expose,omitempty"`

	// +optional
	NodePort int32 `json:"nodePort,omitempty" yaml:"nodePort,omitempty"`

	// +optional
	BufferSize int32 `json:"bufferSize,omitempty" yaml:"bufferSize,omitempty"`

	// +optional
	BodyLimit int32 `json:"bodyLimit,omitempty" yaml:"bodyLimit,omitempty"`
}

func (c *GrpcConfig) getDefaults() *GrpcConfig {
	if c.Port == 0 {
		c.Port = 50000
	}
	if c.Expose == "" {
		c.Expose = "ClusterIP"
	}
	return c
}

func (c *GrpcConfig) SetConfig(config *deployment.Config) *GrpcConfig {
	c.getDefaults()
	if c.Disabled {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GRPC_ENABLE", "false")
		return c
	}

	svc, ok := config.Services["grpc"]
	if ok {
		svc.SetTargetPort(50000).
			SetContainerPort(c.Port)

		if c.Expose == "NodePort" && c.NodePort > 0 {
			svc.SetNodePort(c.NodePort)
		} else {
			svc.SetNodePort(0)
		}

		svc.SetExpose(c.Expose)
	}

	if c.BufferSize != 0 {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GRPC_SUB_BUFF_SIZE", fmt.Sprintf("%d", c.BufferSize))
	}

	if c.BodyLimit != 0 {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_GRPC_BODY_LIMIT", fmt.Sprintf("%d", c.BodyLimit))
	}

	return c
}
