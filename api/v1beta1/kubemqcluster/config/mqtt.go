package config

import (
	"fmt"
	"strconv"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// MqttConfig configures the kubemq-server MQTT connector.
// Maps to server Connectors.MQTT. The connector is opt-in (disabled by default);
// set enabled: true to activate it (opens ports 1883/8883/8083).
type MqttConfig struct {
	// +optional
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port *int32 `json:"port,omitempty" yaml:"port,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	TLSPort *int32 `json:"tlsPort,omitempty" yaml:"tlsPort,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	WSPort *int32 `json:"wsPort,omitempty" yaml:"wsPort,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=events;store;none
	DefaultPattern *string `json:"defaultPattern,omitempty" yaml:"defaultPattern,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10000
	SubBuffSize *int32 `json:"subBuffSize,omitempty" yaml:"subBuffSize,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	QueueAckTimeoutSeconds *int32 `json:"queueAckTimeoutSeconds,omitempty" yaml:"queueAckTimeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	RPCTimeoutSeconds *int32 `json:"rpcTimeoutSeconds,omitempty" yaml:"rpcTimeoutSeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	RPCMaxPending *int32 `json:"rpcMaxPending,omitempty" yaml:"rpcMaxPending,omitempty"`

	// +optional
	DetailHistoryEnabled *bool `json:"detailHistoryEnabled,omitempty" yaml:"detailHistoryEnabled,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	DetailHistoryMaxEntities *int32 `json:"detailHistoryMaxEntities,omitempty" yaml:"detailHistoryMaxEntities,omitempty"`

	// +optional
	Capabilities *MqttCapabilitiesConfig `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer
	Expose *string `json:"expose,omitempty" yaml:"expose,omitempty"`

	// SessionAffinity SHOULD be ClientIP when clients rely on persistent MQTT
	// sessions: session state lives on the replica that accepted the connection.
	// +optional
	// +kubebuilder:validation:Enum=None;ClientIP
	SessionAffinity *string `json:"sessionAffinity,omitempty" yaml:"sessionAffinity,omitempty"`

	// NodePort / TLSNodePort / WSNodePort pin the node ports for 1883 / 8883 / 8083.
	// Honoured only when expose is NodePort; unset leaves them kernel-assigned.
	// +optional
	// +kubebuilder:validation:Minimum=30000
	// +kubebuilder:validation:Maximum=32767
	NodePort *int32 `json:"nodePort,omitempty" yaml:"nodePort,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=30000
	// +kubebuilder:validation:Maximum=32767
	TLSNodePort *int32 `json:"tlsNodePort,omitempty" yaml:"tlsNodePort,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=30000
	// +kubebuilder:validation:Maximum=32767
	WSNodePort *int32 `json:"wsNodePort,omitempty" yaml:"wsNodePort,omitempty"`
}

// MqttCapabilitiesConfig configures the kubemq-server MQTT connector capabilities.
// Maps to server Connectors.MQTT.Capabilities.
type MqttCapabilitiesConfig struct {
	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxClients *int64 `json:"maxClients,omitempty" yaml:"maxClients,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=4294967295
	MaxPacketSizeBytes *int64 `json:"maxPacketSizeBytes,omitempty" yaml:"maxPacketSizeBytes,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65535
	ReceiveMaximum *int32 `json:"receiveMaximum,omitempty" yaml:"receiveMaximum,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65535
	MaxInflight *int32 `json:"maxInflight,omitempty" yaml:"maxInflight,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=4294967295
	MaxSessionExpirySeconds *int64 `json:"maxSessionExpirySeconds,omitempty" yaml:"maxSessionExpirySeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxMessageExpirySeconds *int64 `json:"maxMessageExpirySeconds,omitempty" yaml:"maxMessageExpirySeconds,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2
	MaxQos *int32 `json:"maxQos,omitempty" yaml:"maxQos,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxSubscriptionsPerClient *int32 `json:"maxSubscriptionsPerClient,omitempty" yaml:"maxSubscriptionsPerClient,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=4;5
	MinProtocolVersion *int32 `json:"minProtocolVersion,omitempty" yaml:"minProtocolVersion,omitempty"`
}

func (c *MqttConfig) DeepCopy() *MqttConfig {
	out := &MqttConfig{}

	if c.Enabled != nil {
		out.Enabled = new(bool)
		*out.Enabled = *c.Enabled
	}

	if c.Port != nil {
		out.Port = new(int32)
		*out.Port = *c.Port
	}

	if c.TLSPort != nil {
		out.TLSPort = new(int32)
		*out.TLSPort = *c.TLSPort
	}

	if c.WSPort != nil {
		out.WSPort = new(int32)
		*out.WSPort = *c.WSPort
	}

	if c.DefaultPattern != nil {
		out.DefaultPattern = new(string)
		*out.DefaultPattern = *c.DefaultPattern
	}

	if c.SubBuffSize != nil {
		out.SubBuffSize = new(int32)
		*out.SubBuffSize = *c.SubBuffSize
	}

	if c.QueueAckTimeoutSeconds != nil {
		out.QueueAckTimeoutSeconds = new(int32)
		*out.QueueAckTimeoutSeconds = *c.QueueAckTimeoutSeconds
	}

	if c.RPCTimeoutSeconds != nil {
		out.RPCTimeoutSeconds = new(int32)
		*out.RPCTimeoutSeconds = *c.RPCTimeoutSeconds
	}

	if c.RPCMaxPending != nil {
		out.RPCMaxPending = new(int32)
		*out.RPCMaxPending = *c.RPCMaxPending
	}

	if c.DetailHistoryEnabled != nil {
		out.DetailHistoryEnabled = new(bool)
		*out.DetailHistoryEnabled = *c.DetailHistoryEnabled
	}

	if c.DetailHistoryMaxEntities != nil {
		out.DetailHistoryMaxEntities = new(int32)
		*out.DetailHistoryMaxEntities = *c.DetailHistoryMaxEntities
	}

	if c.Capabilities != nil {
		out.Capabilities = c.Capabilities.DeepCopy()
	}

	if c.Expose != nil {
		out.Expose = new(string)
		*out.Expose = *c.Expose
	}

	if c.SessionAffinity != nil {
		out.SessionAffinity = new(string)
		*out.SessionAffinity = *c.SessionAffinity
	}

	if c.NodePort != nil {
		out.NodePort = new(int32)
		*out.NodePort = *c.NodePort
	}

	if c.TLSNodePort != nil {
		out.TLSNodePort = new(int32)
		*out.TLSNodePort = *c.TLSNodePort
	}

	if c.WSNodePort != nil {
		out.WSNodePort = new(int32)
		*out.WSNodePort = *c.WSNodePort
	}

	return out
}

func (c *MqttCapabilitiesConfig) DeepCopy() *MqttCapabilitiesConfig {
	out := &MqttCapabilitiesConfig{}

	if c.MaxClients != nil {
		out.MaxClients = new(int64)
		*out.MaxClients = *c.MaxClients
	}

	if c.MaxPacketSizeBytes != nil {
		out.MaxPacketSizeBytes = new(int64)
		*out.MaxPacketSizeBytes = *c.MaxPacketSizeBytes
	}

	if c.ReceiveMaximum != nil {
		out.ReceiveMaximum = new(int32)
		*out.ReceiveMaximum = *c.ReceiveMaximum
	}

	if c.MaxInflight != nil {
		out.MaxInflight = new(int32)
		*out.MaxInflight = *c.MaxInflight
	}

	if c.MaxSessionExpirySeconds != nil {
		out.MaxSessionExpirySeconds = new(int64)
		*out.MaxSessionExpirySeconds = *c.MaxSessionExpirySeconds
	}

	if c.MaxMessageExpirySeconds != nil {
		out.MaxMessageExpirySeconds = new(int64)
		*out.MaxMessageExpirySeconds = *c.MaxMessageExpirySeconds
	}

	if c.MaxQos != nil {
		out.MaxQos = new(int32)
		*out.MaxQos = *c.MaxQos
	}

	if c.MaxSubscriptionsPerClient != nil {
		out.MaxSubscriptionsPerClient = new(int32)
		*out.MaxSubscriptionsPerClient = *c.MaxSubscriptionsPerClient
	}

	if c.MinProtocolVersion != nil {
		out.MinProtocolVersion = new(int32)
		*out.MinProtocolVersion = *c.MinProtocolVersion
	}

	return out
}

func (c *MqttConfig) SetConfig(config *deployment.Config) *MqttConfig {
	effective := c.Enabled != nil && *c.Enabled
	config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_ENABLE", strconv.FormatBool(effective))
	if !effective {
		return c
	}

	// Reflect custom ports onto the K8s Service so traffic reaches the listener.
	if svc, ok := config.Services["mqtt"]; ok {
		if c.Port != nil {
			svc.SetPort("mqtt", *c.Port)
		}
		if c.TLSPort != nil {
			svc.SetPort("mqtt-tls", *c.TLSPort)
		}
		if c.WSPort != nil {
			svc.SetPort("mqtt-ws", *c.WSPort)
		}
		applyServiceExposure(svc, c.Expose, c.SessionAffinity, map[string]*int32{
			"mqtt":     c.NodePort,
			"mqtt-tls": c.TLSNodePort,
			"mqtt-ws":  c.WSNodePort,
		})
	}

	if c.Port != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_PORT", fmt.Sprintf("%d", *c.Port))
	}

	if c.TLSPort != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_TLS_PORT", fmt.Sprintf("%d", *c.TLSPort))
	}

	if c.WSPort != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_WS_PORT", fmt.Sprintf("%d", *c.WSPort))
	}

	if c.DefaultPattern != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_DEFAULT_PATTERN", *c.DefaultPattern)
	}

	if c.SubBuffSize != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_SUB_BUFF_SIZE", fmt.Sprintf("%d", *c.SubBuffSize))
	}

	if c.QueueAckTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_QUEUE_ACK_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.QueueAckTimeoutSeconds))
	}

	if c.RPCTimeoutSeconds != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_RPC_TIMEOUT_SECONDS", fmt.Sprintf("%d", *c.RPCTimeoutSeconds))
	}

	if c.RPCMaxPending != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_RPC_MAX_PENDING", fmt.Sprintf("%d", *c.RPCMaxPending))
	}

	if c.DetailHistoryEnabled != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_DETAIL_HISTORY_ENABLED", strconv.FormatBool(*c.DetailHistoryEnabled))
	}

	if c.DetailHistoryMaxEntities != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_DETAIL_HISTORY_MAX_ENTITIES", fmt.Sprintf("%d", *c.DetailHistoryMaxEntities))
	}

	if c.Capabilities != nil {
		if c.Capabilities.MaxClients != nil {
			config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_CAPABILITIES_MAX_CLIENTS", fmt.Sprintf("%d", *c.Capabilities.MaxClients))
		}

		if c.Capabilities.MaxPacketSizeBytes != nil {
			config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_CAPABILITIES_MAX_PACKET_SIZE_BYTES", fmt.Sprintf("%d", *c.Capabilities.MaxPacketSizeBytes))
		}

		if c.Capabilities.ReceiveMaximum != nil {
			config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_CAPABILITIES_RECEIVE_MAXIMUM", fmt.Sprintf("%d", *c.Capabilities.ReceiveMaximum))
		}

		if c.Capabilities.MaxInflight != nil {
			config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_CAPABILITIES_MAX_INFLIGHT", fmt.Sprintf("%d", *c.Capabilities.MaxInflight))
		}

		if c.Capabilities.MaxSessionExpirySeconds != nil {
			config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_CAPABILITIES_MAX_SESSION_EXPIRY_SECONDS", fmt.Sprintf("%d", *c.Capabilities.MaxSessionExpirySeconds))
		}

		if c.Capabilities.MaxMessageExpirySeconds != nil {
			config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_CAPABILITIES_MAX_MESSAGE_EXPIRY_SECONDS", fmt.Sprintf("%d", *c.Capabilities.MaxMessageExpirySeconds))
		}

		if c.Capabilities.MaxQos != nil {
			config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_CAPABILITIES_MAX_QOS", fmt.Sprintf("%d", *c.Capabilities.MaxQos))
		}

		if c.Capabilities.MaxSubscriptionsPerClient != nil {
			config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_CAPABILITIES_MAX_SUBSCRIPTIONS_PER_CLIENT", fmt.Sprintf("%d", *c.Capabilities.MaxSubscriptionsPerClient))
		}

		if c.Capabilities.MinProtocolVersion != nil {
			config.SetConfigMapStringValues(config.Name, "CONNECTORSMQTT_CAPABILITIES_MIN_PROTOCOL_VERSION", fmt.Sprintf("%d", *c.Capabilities.MinProtocolVersion))
		}
	}

	return c
}
