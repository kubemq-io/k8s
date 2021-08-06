package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type HealthConfig struct {
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// +optional
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty"`

	// +optional
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`

	// +optional
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`

	// +optional
	SuccessThreshold int32 `json:"successThreshold,omitempty"`

	// +optional
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
}

func (c *HealthConfig) getDefaults() *HealthConfig {
	if c.SuccessThreshold <= 0 {
		c.SuccessThreshold = 1
	}
	if c.TimeoutSeconds <= 0 {
		c.TimeoutSeconds = 5
	}
	if c.PeriodSeconds <= 0 {
		c.PeriodSeconds = 10
	}
	if c.InitialDelaySeconds <= 0 {
		c.InitialDelaySeconds = 5
	}
	if c.FailureThreshold <= 0 {
		c.FailureThreshold = 12
	}
	return c
}
func (c *HealthConfig) SetConfig(config *deployment.Config) *HealthConfig {
	c.getDefaults()
	if !c.Enabled {
		return c
	}

	tmpl := `          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: %d
            periodSeconds: %d
            timeoutSeconds: %d
            successThreshold: %d
            failureThreshold: %d
`
	prob := fmt.Sprintf(tmpl,
		c.InitialDelaySeconds,
		c.TimeoutSeconds,
		c.PeriodSeconds,
		c.SuccessThreshold,
		c.FailureThreshold)
	config.StatefulSet.SetHealthProbe(prob)
	return c
}
