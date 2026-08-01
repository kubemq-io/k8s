package config

import (
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// antiAffinityTopologyKey spreads replicas across nodes. Node-level is the only
// spread that turns a single node failure back into a single replica failure.
const antiAffinityTopologyKey = "kubernetes.io/hostname"

// PodDisruptionBudgetConfig controls the operator-generated PodDisruptionBudget.
// It is emitted only for a multi-replica, non-standalone cluster — a single pod
// has no quorum to protect and a budget on it would block every node drain.
type PodDisruptionBudgetConfig struct {
	// Enabled turns the budget off. Unset means ON: an unprotected multi-replica
	// cluster can lose raft quorum to a routine node-pool upgrade.
	// +optional
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	// MinAvailable overrides the derived raft majority (replicas/2+1). Setting it
	// below the majority permits an eviction that costs quorum — the operator
	// honours the value but it is not a safe choice.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinAvailable *int32 `json:"minAvailable,omitempty" yaml:"minAvailable,omitempty"`
}

func (c *PodDisruptionBudgetConfig) DeepCopy() *PodDisruptionBudgetConfig {
	out := &PodDisruptionBudgetConfig{}
	if c.Enabled != nil {
		out.Enabled = new(bool)
		*out.Enabled = *c.Enabled
	}
	if c.MinAvailable != nil {
		out.MinAvailable = new(int32)
		*out.MinAvailable = *c.MinAvailable
	}
	return out
}

// IsEnabled reports whether the budget should be emitted. A nil block is enabled.
func (c *PodDisruptionBudgetConfig) IsEnabled() bool {
	return c == nil || c.Enabled == nil || *c.Enabled
}

// PodAntiAffinityConfig controls how replicas are spread across nodes.
type PodAntiAffinityConfig struct {
	// Enabled turns the spread off. Unset means ON.
	// +optional
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	// Required picks the trade-off, and there is no free choice here.
	// false (the default) renders preferredDuringSchedulingIgnoredDuringExecution:
	// pods always schedule, but two replicas MAY land on one node, so a single node
	// failure can still cost quorum. true renders
	// requiredDuringSchedulingIgnoredDuringExecution: the spread is guaranteed, but
	// a replica stays Pending while there are fewer schedulable nodes than replicas.
	// +optional
	Required *bool `json:"required,omitempty" yaml:"required,omitempty"`
}

func (c *PodAntiAffinityConfig) DeepCopy() *PodAntiAffinityConfig {
	out := &PodAntiAffinityConfig{}
	if c.Enabled != nil {
		out.Enabled = new(bool)
		*out.Enabled = *c.Enabled
	}
	if c.Required != nil {
		out.Required = new(bool)
		*out.Required = *c.Required
	}
	return out
}

// IsEnabled reports whether the anti-affinity block should be rendered. A nil
// block is enabled.
func (c *PodAntiAffinityConfig) IsEnabled() bool {
	return c == nil || c.Enabled == nil || *c.Enabled
}

// SetConfig renders the pod-spec `affinity:` block onto the StatefulSet. The
// receiver may be nil (the default-on path). appName is the pod label the
// selector matches.
func (c *PodAntiAffinityConfig) SetConfig(config *deployment.Config) *PodAntiAffinityConfig {
	if !c.IsEnabled() {
		return c
	}
	required := c != nil && c.Required != nil && *c.Required
	config.StatefulSet.SetAffinity(renderPodAntiAffinity(config.Name, required))
	return c
}

// renderPodAntiAffinity builds the 6-space-indented affinity block (the same
// pre-rendered-YAML convention as nodeSelector/resources).
func renderPodAntiAffinity(appName string, required bool) string {
	head := "      affinity:\n        podAntiAffinity:\n"
	if required {
		return head +
			"          requiredDuringSchedulingIgnoredDuringExecution:\n" +
			"            - topologyKey: " + antiAffinityTopologyKey + "\n" +
			"              labelSelector:\n" +
			"                matchLabels:\n" +
			"                  app: " + appName + "\n"
	}
	return head +
		"          preferredDuringSchedulingIgnoredDuringExecution:\n" +
		"            - weight: 100\n" +
		"              podAffinityTerm:\n" +
		"                topologyKey: " + antiAffinityTopologyKey + "\n" +
		"                labelSelector:\n" +
		"                  matchLabels:\n" +
		"                    app: " + appName + "\n"
}
