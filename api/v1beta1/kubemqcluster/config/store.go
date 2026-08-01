package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type StoreConfig struct {
	// Engine selects the persistence engine: "legacy" or "next". next-ness is
	// immutable (born-one-mode) — enforced by a CEL transition rule on the store
	// object in the CRD schema.
	//
	// Nil does NOT mean "let the server decide": the server auto-detects the engine
	// from the store directory when it is not told, and a pod starting on an EMPTY
	// volume resolves next. The operator therefore resolves the engine itself
	// (established-engine annotation -> this field -> ConfigMap key, fail-closed) and
	// always pins the result onto the pod. A nil here means "operator-resolved",
	// which for a fresh non-Kafka cluster is legacy.
	//
	// The enum stays legacy|next. The server also accepts "auto", deliberately not
	// exposed: the operator's fail-closed derivation is stricter than the server's
	// detection and stays authoritative for managed clusters.
	// +optional
	// +kubebuilder:validation:Enum=legacy;next
	Engine *string `json:"engine,omitempty" yaml:"engine,omitempty"`

	// +optional
	Clean bool `json:"clean,omitempty" yaml:"clean,omitempty"`

	// +optional
	Path string `json:"path,omitempty" yaml:"path,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxChannels *int32 `json:"maxChannels,omitempty" yaml:"maxChannels,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxSubscribers *int32 `json:"maxSubscribers,omitempty" yaml:"maxSubscribers,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxMessages *int32 `json:"maxMessages,omitempty" yaml:"maxMessages,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxChannelSize *int32 `json:"maxChannelSize,omitempty" yaml:"maxChannelSize,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MessagesRetentionMinutes *int32 `json:"messagesRetentionMinutes,omitempty" yaml:"messagesRetentionMinutes,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	PurgeInactiveMinutes *int32 `json:"purgeInactiveMinutes,omitempty" yaml:"purgeInactiveMinutes,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	IdlePruneCutoffHours *int32 `json:"idlePruneCutoffHours,omitempty" yaml:"idlePruneCutoffHours,omitempty"`
}

func (c *StoreConfig) DeepCopy() *StoreConfig {
	out := &StoreConfig{}

	if c.Engine != nil {
		out.Engine = new(string)
		*out.Engine = *c.Engine
	}

	out.Clean = c.Clean
	out.Path = c.Path

	if c.MaxChannels != nil {
		out.MaxChannels = new(int32)
		*out.MaxChannels = *c.MaxChannels
	}

	if c.MaxSubscribers != nil {
		out.MaxSubscribers = new(int32)
		*out.MaxSubscribers = *c.MaxSubscribers
	}

	if c.MaxMessages != nil {
		out.MaxMessages = new(int32)
		*out.MaxMessages = *c.MaxMessages
	}

	if c.MaxChannelSize != nil {
		out.MaxChannelSize = new(int32)
		*out.MaxChannelSize = *c.MaxChannelSize
	}

	if c.MessagesRetentionMinutes != nil {
		out.MessagesRetentionMinutes = new(int32)
		*out.MessagesRetentionMinutes = *c.MessagesRetentionMinutes
	}

	if c.PurgeInactiveMinutes != nil {
		out.PurgeInactiveMinutes = new(int32)
		*out.PurgeInactiveMinutes = *c.PurgeInactiveMinutes
	}

	if c.IdlePruneCutoffHours != nil {
		out.IdlePruneCutoffHours = new(int32)
		*out.IdlePruneCutoffHours = *c.IdlePruneCutoffHours
	}

	return out
}
func (c *StoreConfig) SetConfig(config *deployment.Config) *StoreConfig {
	// Emit STORE_ENGINE for any explicit value (including "legacy"). A nil Engine
	// emits nothing HERE — the operator pins the resolved engine after this call, so
	// the key still reaches every pod. Do not treat a nil as "server decides": unset
	// now means auto-detect from the store directory, and an empty volume resolves
	// next.
	if c.Engine != nil && *c.Engine != "" {
		config.SetConfigMapStringValues(config.Name, "STORE_ENGINE", *c.Engine)
	}

	if c.Clean {
		config.SetConfigMapStringValues(config.Name, "STORE_CLEAN_STORE", "true")
	}

	if c.Path != "" {
		config.SetConfigMapStringValues(config.Name, "STORE_STORE_PATH", c.Path)
	}

	if c.MaxChannels != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_QUEUES", fmt.Sprintf("%d", *c.MaxChannels))
	}

	if c.MaxSubscribers != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_SUBSCRIBERS", fmt.Sprintf("%d", *c.MaxSubscribers))
	}

	if c.MaxMessages != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_MESSAGES", fmt.Sprintf("%d", *c.MaxMessages))
	}

	if c.MaxChannelSize != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_QUEUE_SIZE", fmt.Sprintf("%d", *c.MaxChannelSize))
	}

	if c.MessagesRetentionMinutes != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_RETENTION", fmt.Sprintf("%d", *c.MessagesRetentionMinutes))
	}

	if c.PurgeInactiveMinutes != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_MAX_PURGE_INACTIVE", fmt.Sprintf("%d", *c.PurgeInactiveMinutes))
	}

	if c.IdlePruneCutoffHours != nil {
		config.SetConfigMapStringValues(config.Name, "STORE_IDLE_PRUNE_CUTOFF_HOURS", fmt.Sprintf("%d", *c.IdlePruneCutoffHours))
	}

	return c
}
