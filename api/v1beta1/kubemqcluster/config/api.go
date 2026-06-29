package config

import (
	"fmt"
	"strconv"
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

	// +optional
	// Auth configures the opt-in management-API / dashboard authentication
	// (server [Api.Auth] / API_AUTH_*). Disabled by default.
	Auth *ApiAuthConfig `json:"auth,omitempty" yaml:"auth,omitempty"`
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

	if c.Auth != nil {
		out.Auth = c.Auth.DeepCopy()
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

	if c.Auth != nil {
		c.Auth.SetConfig(config)
	}

	return c
}

// ApiAuthConfig configures the opt-in management-API / dashboard authentication
// layer (server [Api.Auth] / API_AUTH_*). It is disabled by default. When
// enabled, a seed admin credential is required: either reference an existing
// Secret via AdminSecretRef, or let the operator generate one once and retain
// it (see the operator's seed-admin Secret reconciliation). The session/proxy
// fields and the admin username map straight to server env vars; the admin
// PASSWORD is never placed here — it is injected into the pod from a Secret.
type ApiAuthConfig struct {
	// +optional
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`

	// +optional
	SessionIdleMinutes *int32 `json:"sessionIdleMinutes,omitempty" yaml:"sessionIdleMinutes,omitempty"`

	// +optional
	SessionAbsoluteHours *int32 `json:"sessionAbsoluteHours,omitempty" yaml:"sessionAbsoluteHours,omitempty"`

	// +optional
	StorePath *string `json:"storePath,omitempty" yaml:"storePath,omitempty"`

	// +optional
	TrustedTLSProxy *bool `json:"trustedTLSProxy,omitempty" yaml:"trustedTLSProxy,omitempty"`

	// +optional
	// AdminUsername sets KUBEMQ_API_ADMIN_USERNAME (server default: "admin").
	AdminUsername *string `json:"adminUsername,omitempty" yaml:"adminUsername,omitempty"`

	// +optional
	// AdminSecretRef names an existing Secret holding the seed admin password.
	// When empty (and auth enabled) the operator generates one once and retains
	// it across reconciles. The key defaults to AdminSecretKey ("admin-password").
	AdminSecretRef *string `json:"adminSecretRef,omitempty" yaml:"adminSecretRef,omitempty"`

	// +optional
	// AdminSecretKey is the key within AdminSecretRef (or the operator-generated
	// Secret) holding the password. Defaults to "admin-password".
	AdminSecretKey *string `json:"adminSecretKey,omitempty" yaml:"adminSecretKey,omitempty"`
}

func (c *ApiAuthConfig) DeepCopy() *ApiAuthConfig {
	out := &ApiAuthConfig{}

	if c.Enable != nil {
		out.Enable = new(bool)
		*out.Enable = *c.Enable
	}
	if c.SessionIdleMinutes != nil {
		out.SessionIdleMinutes = new(int32)
		*out.SessionIdleMinutes = *c.SessionIdleMinutes
	}
	if c.SessionAbsoluteHours != nil {
		out.SessionAbsoluteHours = new(int32)
		*out.SessionAbsoluteHours = *c.SessionAbsoluteHours
	}
	if c.StorePath != nil {
		out.StorePath = new(string)
		*out.StorePath = *c.StorePath
	}
	if c.TrustedTLSProxy != nil {
		out.TrustedTLSProxy = new(bool)
		*out.TrustedTLSProxy = *c.TrustedTLSProxy
	}
	if c.AdminUsername != nil {
		out.AdminUsername = new(string)
		*out.AdminUsername = *c.AdminUsername
	}
	if c.AdminSecretRef != nil {
		out.AdminSecretRef = new(string)
		*out.AdminSecretRef = *c.AdminSecretRef
	}
	if c.AdminSecretKey != nil {
		out.AdminSecretKey = new(string)
		*out.AdminSecretKey = *c.AdminSecretKey
	}

	return out
}

// SetConfig emits the API_AUTH_* env vars (consumed by the server via viper) and
// KUBEMQ_API_ADMIN_USERNAME (read by the server via os.Getenv) into the cluster
// ConfigMap. The admin PASSWORD is intentionally NOT emitted here: the operator
// injects KUBEMQ_API_ADMIN_PASSWORD into the pod from a Secret (valueFrom).
func (c *ApiAuthConfig) SetConfig(config *deployment.Config) *ApiAuthConfig {
	if c.Enable != nil {
		config.SetConfigMapStringValues(config.Name, "API_AUTH_ENABLE", strconv.FormatBool(*c.Enable))
	}
	if c.SessionIdleMinutes != nil {
		config.SetConfigMapStringValues(config.Name, "API_AUTH_SESSION_IDLE_MINUTES", fmt.Sprintf("%d", *c.SessionIdleMinutes))
	}
	if c.SessionAbsoluteHours != nil {
		config.SetConfigMapStringValues(config.Name, "API_AUTH_SESSION_ABSOLUTE_HOURS", fmt.Sprintf("%d", *c.SessionAbsoluteHours))
	}
	if c.StorePath != nil {
		config.SetConfigMapStringValues(config.Name, "API_AUTH_STORE_PATH", *c.StorePath)
	}
	if c.TrustedTLSProxy != nil {
		config.SetConfigMapStringValues(config.Name, "API_AUTH_TRUSTED_TLS_PROXY", strconv.FormatBool(*c.TrustedTLSProxy))
	}
	if c.AdminUsername != nil {
		config.SetConfigMapStringValues(config.Name, "KUBEMQ_API_ADMIN_USERNAME", *c.AdminUsername)
	}

	return c
}
