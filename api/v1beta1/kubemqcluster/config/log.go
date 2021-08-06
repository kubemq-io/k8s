package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type LogConfig struct {
	// +optional
	Level *int32 `json:"level,omitempty"`

	// +optional
	File string `json:"file,omitempty"`
}

func (c *LogConfig) DeepCopy() *LogConfig {
	out := &LogConfig{}
	out.File = c.File
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

	if c.File != "" {
		config.SetConfigMapStringValues(config.Name, "LOG_FILE_ENABLE", "true")
		config.SetConfigMapStringValues(config.Name, "LOG_FILE_PATH", c.File)
	}

	return c
}
