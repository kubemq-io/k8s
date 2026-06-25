package config

import (
	"fmt"
	"strconv"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// AuditConfig configures the kubemq-server audit log. Maps to server Audit.
//
// Enable is *bool deliberately: audit defaults ON server-side, so a plain bool
// could never emit "false" to disable it.
type AuditConfig struct {
	// +optional
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	RetentionHours *int32 `json:"retentionHours,omitempty" yaml:"retentionHours,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	CleanupIntervalMinutes *int32 `json:"cleanupIntervalMinutes,omitempty" yaml:"cleanupIntervalMinutes,omitempty"`
}

func (c *AuditConfig) DeepCopy() *AuditConfig {
	out := &AuditConfig{}

	if c.Enable != nil {
		out.Enable = new(bool)
		*out.Enable = *c.Enable
	}

	if c.RetentionHours != nil {
		out.RetentionHours = new(int32)
		*out.RetentionHours = *c.RetentionHours
	}

	if c.CleanupIntervalMinutes != nil {
		out.CleanupIntervalMinutes = new(int32)
		*out.CleanupIntervalMinutes = *c.CleanupIntervalMinutes
	}

	return out
}

func (c *AuditConfig) SetConfig(config *deployment.Config) *AuditConfig {
	if c.Enable != nil {
		config.SetConfigMapStringValues(config.Name, "AUDIT_ENABLE", strconv.FormatBool(*c.Enable))
	}

	if c.RetentionHours != nil {
		config.SetConfigMapStringValues(config.Name, "AUDIT_RETENTION_HOURS", fmt.Sprintf("%d", *c.RetentionHours))
	}

	if c.CleanupIntervalMinutes != nil {
		config.SetConfigMapStringValues(config.Name, "AUDIT_CLEANUP_INTERVAL_MINUTES", fmt.Sprintf("%d", *c.CleanupIntervalMinutes))
	}

	return c
}
