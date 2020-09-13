package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1alpha1/kubemqcluster/deployment"
	"strings"
)

type NodeSelectorConfig struct {
	// +optional
	Keys map[string]string `json:"keys,omitempty"`
}

func (c *NodeSelectorConfig) SetConfig(config *deployment.Config) *NodeSelectorConfig {
	if len(c.Keys) == 0 {
		return nil
	}
	tmpl := []string{"      nodeSelector:\n"}
	for key, value := range c.Keys {
		tmpl = append(tmpl, fmt.Sprintf("        %s: %s\n", key, value))
	}
	config.StatefulSet.SetNodeSelectors(strings.Join(tmpl, ""))
	return c
}

func (c *NodeSelectorConfig) DeepCopy() *NodeSelectorConfig {
	out := &NodeSelectorConfig{
		Keys: map[string]string{},
	}
	if out.Keys == nil {
		out.Keys = map[string]string{}
	}
	for key, value := range c.Keys {
		out.Keys[key] = value
	}
	return out
}
