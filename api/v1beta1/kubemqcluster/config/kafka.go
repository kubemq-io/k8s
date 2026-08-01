package config

import (
	"fmt"
	"strconv"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// KafkaConfig configures the kubemq-server Kafka drop-in connector.
// Maps to server Connectors.Kafka. The connector is opt-in (disabled by default);
// set enabled: true to activate it (opens ports 9092/9093). The `enabled` key
// matches every sibling connector's CRD convention (mqtt/amqp/stomp/aws/gcp),
// reconciling the proposed `enable:` from the server repo's deploy/kafka sample.
//
// The field set mirrors the env-mapped block authored in the kubemq-server repo
// (deploy/kafka/crd-sample.yaml + deploy/kafka/README.md §"Env mapping table",
// WP-7.3/D105). SASL Credentials are deliberately NOT here — they have no env
// binding by design (file/structured config only) and need a Secret-mount path.
//
// ENV-NAME NOTE (LOAD-BEARING): "Kafka" is Titlecase, so the server's
// convertEnvFormat (kubemq-server config/env.go) yields UNDERSCORED keys —
// CONNECTORS_KAFKA_ENABLE / _PORT / _TLS_PORT / _ADVERTISED_HOST /
// _ADVERTISED_PORT / _MAX_CONNECTIONS / _MAX_MESSAGE_BYTES — NOT the glued
// CONNECTORSKAFKA_* form that all-caps connectors (MQTT/CE) use. A single missing
// underscore leaves the connector silently unconfigured (the M-23 footgun class).
type KafkaConfig struct {
	// +optional
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	// +optional
	Port *string `json:"port,omitempty" yaml:"port,omitempty"`

	// +optional
	TLSPort *string `json:"tlsPort,omitempty" yaml:"tlsPort,omitempty"`

	// AdvertisedHost is the host clients are told to reconnect to (Metadata). SET
	// it to the external LB DNS / NodePort IP for out-of-cluster access, or the
	// ClusterIP Service DNS for in-cluster-only — leaving it "" hangs every
	// external client (the M-23 connect-then-hang footgun). The global
	// Security.Cert SAN must include this value or TLS on 9093 fails hostname
	// verification.
	// +optional
	AdvertisedHost *string `json:"advertisedHost,omitempty" yaml:"advertisedHost,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	AdvertisedPort *int32 `json:"advertisedPort,omitempty" yaml:"advertisedPort,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxConnections *int32 `json:"maxConnections,omitempty" yaml:"maxConnections,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1073741824
	MaxMessageBytes *int64 `json:"maxMessageBytes,omitempty" yaml:"maxMessageBytes,omitempty"`

	// Peers is the per-broker Kafka advertised-address map for a CLUSTERED cluster,
	// `id@host:port` comma-separated, one entry per replica, ids matching the raft
	// replica ids (ordinal+1). Kafka clients are handed a broker address per broker,
	// so one round-robin Service address cannot serve a multi-broker cluster — this
	// is the standard Kafka advertised-listeners problem on Kubernetes.
	//
	// LEAVE IT UNSET for in-cluster access: the operator derives it from the
	// StatefulSet's stable pod DNS names and gives each pod its own advertised host.
	// SET IT for external access, where each broker needs its own client-reachable
	// address (a Service or LoadBalancer per broker, or a NodePort per pod).
	//
	// When it is set the server runs in cluster mode, where a broker's own address —
	// for both the Metadata broker list and the FindCoordinator answer — comes from
	// its entry in THIS map. advertisedHost is not consulted for addressing and is
	// best left unset; a single value could only ever be right for one broker.
	// +optional
	Peers *string `json:"peers,omitempty" yaml:"peers,omitempty"`

	// Expose controls the K8s Service .spec.type for the kafka/kafka-tls ports.
	// Unset leaves the catalog default (ClusterIP) untouched.
	//
	// NodePort/LoadBalancer DOES NOT give a multi-replica cluster working external
	// Kafka. A Service is ONE address that round-robins across brokers, while Kafka
	// hands clients an address per broker. External multi-broker access additionally
	// requires per-broker addressing that YOU provision (a Service/LoadBalancer per pod,
	// or a per-pod NodePort), listed in Peers. Exposing without Peers on a clustered
	// cluster is rejected by the operator rather than left to fail at runtime.
	// +optional
	// +kubebuilder:default=ClusterIP
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer
	Expose *string `json:"expose,omitempty" yaml:"expose,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=None;ClientIP
	SessionAffinity *string `json:"sessionAffinity,omitempty" yaml:"sessionAffinity,omitempty"`

	// NodePort / TLSNodePort pin the node ports for 9092 / 9093. Honoured only when
	// expose is NodePort; unset leaves them kernel-assigned. Note that a single
	// NodePort Service still round-robins across brokers — a multi-broker external
	// deployment needs per-broker addressing via peers.
	// +optional
	// +kubebuilder:validation:Minimum=30000
	// +kubebuilder:validation:Maximum=32767
	NodePort *int32 `json:"nodePort,omitempty" yaml:"nodePort,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=30000
	// +kubebuilder:validation:Maximum=32767
	TLSNodePort *int32 `json:"tlsNodePort,omitempty" yaml:"tlsNodePort,omitempty"`
}

func (c *KafkaConfig) DeepCopy() *KafkaConfig {
	out := &KafkaConfig{}

	if c.Enabled != nil {
		out.Enabled = new(bool)
		*out.Enabled = *c.Enabled
	}

	if c.Port != nil {
		out.Port = new(string)
		*out.Port = *c.Port
	}

	if c.TLSPort != nil {
		out.TLSPort = new(string)
		*out.TLSPort = *c.TLSPort
	}

	if c.AdvertisedHost != nil {
		out.AdvertisedHost = new(string)
		*out.AdvertisedHost = *c.AdvertisedHost
	}

	if c.AdvertisedPort != nil {
		out.AdvertisedPort = new(int32)
		*out.AdvertisedPort = *c.AdvertisedPort
	}

	if c.MaxConnections != nil {
		out.MaxConnections = new(int32)
		*out.MaxConnections = *c.MaxConnections
	}

	if c.MaxMessageBytes != nil {
		out.MaxMessageBytes = new(int64)
		*out.MaxMessageBytes = *c.MaxMessageBytes
	}

	if c.Peers != nil {
		out.Peers = new(string)
		*out.Peers = *c.Peers
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

	if c.Expose != nil {
		out.Expose = new(string)
		*out.Expose = *c.Expose
	}

	return out
}

func (c *KafkaConfig) SetConfig(config *deployment.Config) *KafkaConfig {
	effective := c.Enabled != nil && *c.Enabled
	config.SetConfigMapStringValues(config.Name, "CONNECTORS_KAFKA_ENABLE", strconv.FormatBool(effective))
	if !effective {
		return c
	}

	// Reflect custom ports onto the K8s Service so traffic reaches the listener
	// (Port/TLSPort are strings on the server — "9092"/"9093"; a non-numeric value
	// leaves the catalog default untouched rather than erroring).
	if svc, ok := config.Services["kafka"]; ok {
		if c.Port != nil {
			if p, err := strconv.Atoi(*c.Port); err == nil {
				svc.SetPort("kafka", int32(p))
			}
		}
		if c.TLSPort != nil {
			if p, err := strconv.Atoi(*c.TLSPort); err == nil {
				svc.SetPort("kafka-tls", int32(p))
			}
		}
		applyServiceExposure(svc, c.Expose, c.SessionAffinity, map[string]*int32{
			"kafka":     c.NodePort,
			"kafka-tls": c.TLSNodePort,
		})
	}

	if c.Peers != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_KAFKA_PEERS", *c.Peers)
	}

	if c.Port != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_KAFKA_PORT", *c.Port)
	}

	if c.TLSPort != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_KAFKA_TLS_PORT", *c.TLSPort)
	}

	if c.AdvertisedHost != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_KAFKA_ADVERTISED_HOST", *c.AdvertisedHost)
	} else {
		// F4: leaving advertisedHost unset hangs every external client on
		// reconnect (the M-23 connect-then-hang footgun) — default to the
		// in-cluster short-form Service DNS name so at least in-cluster
		// clients work out of the box. Custom-DNS-domain safe (short form,
		// no ".cluster.local" suffix). Explicit value above always wins.
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_KAFKA_ADVERTISED_HOST", config.Name+"-kafka."+config.Namespace+".svc")
	}

	if c.AdvertisedPort != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_KAFKA_ADVERTISED_PORT", fmt.Sprintf("%d", *c.AdvertisedPort))
	}

	if c.MaxConnections != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_KAFKA_MAX_CONNECTIONS", fmt.Sprintf("%d", *c.MaxConnections))
	}

	if c.MaxMessageBytes != nil {
		config.SetConfigMapStringValues(config.Name, "CONNECTORS_KAFKA_MAX_MESSAGE_BYTES", fmt.Sprintf("%d", *c.MaxMessageBytes))
	}

	return c
}
