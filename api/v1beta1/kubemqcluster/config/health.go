package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type HealthConfig struct {
	// +optional
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	// +optional
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty" yaml:"initialDelaySeconds,omitempty"`

	// +optional
	PeriodSeconds int32 `json:"periodSeconds,omitempty" yaml:"periodSeconds,omitempty"`

	// +optional
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty" yaml:"timeoutSeconds,omitempty"`

	// +optional
	SuccessThreshold int32 `json:"successThreshold,omitempty" yaml:"successThreshold,omitempty"`

	// +optional
	FailureThreshold int32 `json:"failureThreshold,omitempty" yaml:"failureThreshold,omitempty"`
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
	// If spec.api.disabled is set together with spec.health.enabled (a user misconfig),
	// ApiConfig.SetConfig returns before SetApiPort so ApiPort stays 0 here; the fallback
	// keeps the probe on the default 8080 (today's hardcoded behavior) instead of port: 0.
	apiPort := config.StatefulSet.ApiPort
	if apiPort == 0 {
		apiPort = 8080
	}
	// The API binds 127.0.0.1 by default; a kubelet httpGet probe targets the pod IP,
	// so the probe is only reachable when the API binds 0.0.0.0.
	config.SetConfigMapStringValues(config.Name, "API_BIND_ADDRESS", "0.0.0.0")

	tmpl := `          livenessProbe:
            httpGet:
              path: /health
              port: %d
            initialDelaySeconds: %d
            periodSeconds: %d
            timeoutSeconds: %d
            successThreshold: %d
            failureThreshold: %d
`
	prob := fmt.Sprintf(tmpl, apiPort,
		c.InitialDelaySeconds,
		c.PeriodSeconds,
		c.TimeoutSeconds,
		c.SuccessThreshold,
		c.FailureThreshold)
	config.StatefulSet.SetHealthProbe(prob)
	return c
}
