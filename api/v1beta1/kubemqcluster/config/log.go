package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type LogConfig struct {
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=5
	// Level: 0=Trace 1=Debug 2=Info 3=Warn 4=Error 5=Fatal
	Level *int32 `json:"level,omitempty" yaml:"level,omitempty"`
}

func (c *LogConfig) DeepCopy() *LogConfig {
	out := &LogConfig{}
	if c.Level != nil {
		out.Level = new(int32)
		*out.Level = *c.Level
	}
	return out
}
func (c *LogConfig) SetConfig(config *deployment.Config) *LogConfig {

	if c.Level != nil {
		config.SetConfigMapStringValues(config.Name, "LOG_LEVEL", fmt.Sprintf("%d", *c.Level))
	}

	return c
}
