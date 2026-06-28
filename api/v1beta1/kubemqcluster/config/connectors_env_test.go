package config

import (
	"encoding/base64"
	"testing"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Env-var names below are the exact literals the kubemq-server binds via
// convertEnvFormat (kubemq-server/config/env.go): ToSnakeCase -> strip "." -> upper.
// They include the acronym traps (CONNECTORSMCP_, CONNECTORSA2_A_, CONNECTORSCE_).
// If the server's key naming changes, update these AND the SetConfig literals together.

const testClusterName = "kubemq-cluster"

func newTestConfig() *deployment.Config {
	return deployment.DefaultKubeMQManifestConfig("", testClusterName, "kubemq", false)
}

func vars(cfg *deployment.Config) map[string]string {
	return cfg.ConfigMaps[testClusterName].Variables
}

// secData returns the cluster Secret's data variables (base64-encoded values).
func secData(cfg *deployment.Config) map[string]string {
	return cfg.Secrets[testClusterName].DataVariables
}

func ptr32(v int32) *int32    { return &v }
func ptr64(v int64) *int64    { return &v }
func strptr(v string) *string { return &v }
func boolptr(v bool) *bool    { return &v }

func TestMcpConfig_SetConfig_Disabled(t *testing.T) {
	cfg := newTestConfig()
	(&McpConfig{Disabled: true}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORSMCP_ENABLE"])
}

func TestMcpConfig_SetConfig_Scalars(t *testing.T) {
	cfg := newTestConfig()
	(&McpConfig{
		ToolTimeoutSeconds: ptr32(30),
		TrustedOrigins:     []string{"*", "x"},
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 2)
	assert.Equal(t, "30", v["CONNECTORSMCP_TOOL_TIMEOUT_SECONDS"])
	assert.Equal(t, "*,x", v["CONNECTORSMCP_TRUSTED_ORIGINS"])
	_, hasEnable := v["CONNECTORSMCP_ENABLE"]
	assert.False(t, hasEnable, "enable must not be written when not disabled")
}

func TestMcpConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&McpConfig{}).SetConfig(cfg)
	assert.Len(t, vars(cfg), 0)
}

func TestAgentsConfig_SetConfig_Disabled(t *testing.T) {
	cfg := newTestConfig()
	(&AgentsConfig{Disabled: true}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORSA2_A_ENABLE"])
}

func TestAgentsConfig_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&AgentsConfig{
		AgentTTLSeconds:       ptr32(300),
		DefaultTimeoutSeconds: ptr32(300),
		MaxTimeoutSeconds:     ptr32(3600),
		MaxAgents:             ptr32(0),
		MaxSSEIdleSeconds:     ptr32(300),
		TrustedOrigins:        []string{"https://a", "https://b"},
		AgentMaxResponseBytes: ptr64(10485760),
		AgentTLSSkipVerify:    true,
		AgentMaxConcurrency:   ptr32(100),
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 9) // all 9 non-disable fields written
	assert.Equal(t, "300", v["CONNECTORSA2_A_AGENT_TTL_SECONDS"])
	assert.Equal(t, "300", v["CONNECTORSA2_A_DEFAULT_TIMEOUT_SECONDS"])
	assert.Equal(t, "3600", v["CONNECTORSA2_A_MAX_TIMEOUT_SECONDS"])
	assert.Equal(t, "0", v["CONNECTORSA2_A_MAX_AGENTS"])
	assert.Equal(t, "300", v["CONNECTORSA2_A_MAX_SSE_IDLE_SECONDS"])
	assert.Equal(t, "https://a,https://b", v["CONNECTORSA2_A_TRUSTED_ORIGINS"])
	assert.Equal(t, "10485760", v["CONNECTORSA2_A_AGENT_MAX_RESPONSE_BYTES"])
	assert.Equal(t, "true", v["CONNECTORSA2_A_AGENT_TLS_SKIP_VERIFY"])
	assert.Equal(t, "100", v["CONNECTORSA2_A_AGENT_MAX_CONCURRENCY"])
	_, hasEnable := v["CONNECTORSA2_A_ENABLE"]
	assert.False(t, hasEnable)
}

func TestAgentsConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&AgentsConfig{}).SetConfig(cfg)
	assert.Len(t, vars(cfg), 0)
}

func TestCeConfig_SetConfig_Disabled(t *testing.T) {
	cfg := newTestConfig()
	(&CeConfig{Disabled: true}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORSCE_ENABLE"])
}

func TestCeConfig_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&CeConfig{
		TimeoutSeconds:    ptr32(60),
		SubBuffSize:       ptr32(100),
		MaxSSEIdleSeconds: ptr32(300),
		MaxSSEConnections: ptr32(0),
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 4)
	assert.Equal(t, "60", v["CONNECTORSCE_TIMEOUT_SECONDS"])
	assert.Equal(t, "100", v["CONNECTORSCE_SUB_BUFF_SIZE"])
	assert.Equal(t, "300", v["CONNECTORSCE_MAX_SSE_IDLE_SECONDS"])
	assert.Equal(t, "0", v["CONNECTORSCE_MAX_SSE_CONNECTIONS"])
}

func TestCeConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&CeConfig{}).SetConfig(cfg)
	assert.Len(t, vars(cfg), 0)
}

func TestMqttConfig_SetConfig_NilEnabled(t *testing.T) {
	cfg := newTestConfig()
	(&MqttConfig{}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORSMQTT_ENABLE"])
}

func TestMqttConfig_SetConfig_ExplicitFalse(t *testing.T) {
	cfg := newTestConfig()
	(&MqttConfig{Enabled: boolptr(false)}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORSMQTT_ENABLE"])
}

func TestMqttConfig_SetConfig_ExplicitTrue(t *testing.T) {
	cfg := newTestConfig()
	(&MqttConfig{Enabled: boolptr(true)}).SetConfig(cfg)
	v := vars(cfg)
	assert.Equal(t, "true", v["CONNECTORSMQTT_ENABLE"])
}

func TestMqttConfig_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&MqttConfig{
		Enabled:                boolptr(true),
		Port:                   ptr32(1883),
		TLSPort:                ptr32(8883),
		WSPort:                 ptr32(8083),
		DefaultPattern:         strptr("events"),
		SubBuffSize:            ptr32(100),
		QueueAckTimeoutSeconds: ptr32(30),
		RPCTimeoutSeconds:      ptr32(30),
		RPCMaxPending:          ptr32(1024),
		Capabilities: &MqttCapabilitiesConfig{
			MaxClients:              ptr64(0),
			MaxPacketSizeBytes:      ptr64(4194304),
			ReceiveMaximum:          ptr32(1024),
			MaxInflight:             ptr32(8192),
			MaxSessionExpirySeconds: ptr64(3600),
			MaxMessageExpirySeconds: ptr64(86400),
			MaxQos:                  ptr32(2),
			MinProtocolVersion:      ptr32(4),
		},
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 17) // ENABLE + 8 scalar + 8 capability fields
	assert.Equal(t, "true", v["CONNECTORSMQTT_ENABLE"])
	assert.Equal(t, "1883", v["CONNECTORSMQTT_PORT"])
	assert.Equal(t, "8883", v["CONNECTORSMQTT_TLS_PORT"])
	assert.Equal(t, "8083", v["CONNECTORSMQTT_WS_PORT"])
	assert.Equal(t, "events", v["CONNECTORSMQTT_DEFAULT_PATTERN"])
	assert.Equal(t, "100", v["CONNECTORSMQTT_SUB_BUFF_SIZE"])
	assert.Equal(t, "30", v["CONNECTORSMQTT_QUEUE_ACK_TIMEOUT_SECONDS"])
	assert.Equal(t, "30", v["CONNECTORSMQTT_RPC_TIMEOUT_SECONDS"])
	assert.Equal(t, "1024", v["CONNECTORSMQTT_RPC_MAX_PENDING"])
	assert.Equal(t, "0", v["CONNECTORSMQTT_CAPABILITIES_MAX_CLIENTS"])
	assert.Equal(t, "4194304", v["CONNECTORSMQTT_CAPABILITIES_MAX_PACKET_SIZE_BYTES"])
	assert.Equal(t, "1024", v["CONNECTORSMQTT_CAPABILITIES_RECEIVE_MAXIMUM"])
	assert.Equal(t, "8192", v["CONNECTORSMQTT_CAPABILITIES_MAX_INFLIGHT"])
	assert.Equal(t, "3600", v["CONNECTORSMQTT_CAPABILITIES_MAX_SESSION_EXPIRY_SECONDS"])
	assert.Equal(t, "86400", v["CONNECTORSMQTT_CAPABILITIES_MAX_MESSAGE_EXPIRY_SECONDS"])
	assert.Equal(t, "2", v["CONNECTORSMQTT_CAPABILITIES_MAX_QOS"])
	assert.Equal(t, "4", v["CONNECTORSMQTT_CAPABILITIES_MIN_PROTOCOL_VERSION"])
}

func TestMqttConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&MqttConfig{}).SetConfig(cfg)
	// present-empty block always emits ENABLE=false (opt-in: nil → off)
	assert.Len(t, vars(cfg), 1)
	assert.Equal(t, "false", vars(cfg)["CONNECTORSMQTT_ENABLE"])
}

func TestAmqpConfig_SetConfig_NilEnabled(t *testing.T) {
	cfg := newTestConfig()
	(&AmqpConfig{}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_AMQP_ENABLE"])
}

func TestAmqpConfig_SetConfig_ExplicitFalse(t *testing.T) {
	cfg := newTestConfig()
	(&AmqpConfig{Enabled: boolptr(false)}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_AMQP_ENABLE"])
}

func TestAmqpConfig_SetConfig_ExplicitTrue(t *testing.T) {
	cfg := newTestConfig()
	(&AmqpConfig{Enabled: boolptr(true)}).SetConfig(cfg)
	v := vars(cfg)
	assert.Equal(t, "true", v["CONNECTORS_AMQP_ENABLE"])
}

func TestAmqpConfig_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&AmqpConfig{
		Enabled:           boolptr(true),
		Port:              ptr32(5672),
		TLSPort:           ptr32(5671),
		HeartbeatSeconds:  ptr32(60),
		FrameMax:          ptr32(131072),
		ChannelMax:        ptr32(2047),
		MaxConnections:    ptr32(1000),
		MaxBodySize:       ptr32(104857600),
		DefaultVhost:      strptr("default"),
		GetBatchSize:      ptr32(32),
		DeadLetterMaxHops: ptr32(16),
		MaxReceiveCount:   ptr32(0),
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 12) // ENABLE + 11 scalar fields
	assert.Equal(t, "true", v["CONNECTORS_AMQP_ENABLE"])
	assert.Equal(t, "5672", v["CONNECTORS_AMQP_PORT"])
	assert.Equal(t, "5671", v["CONNECTORS_AMQP_TLS_PORT"])
	assert.Equal(t, "60", v["CONNECTORS_AMQP_HEARTBEAT_SECONDS"])
	assert.Equal(t, "131072", v["CONNECTORS_AMQP_FRAME_MAX"])
	assert.Equal(t, "2047", v["CONNECTORS_AMQP_CHANNEL_MAX"])
	assert.Equal(t, "1000", v["CONNECTORS_AMQP_MAX_CONNECTIONS"])
	assert.Equal(t, "104857600", v["CONNECTORS_AMQP_MAX_BODY_SIZE"])
	assert.Equal(t, "default", v["CONNECTORS_AMQP_DEFAULT_VHOST"])
	assert.Equal(t, "32", v["CONNECTORS_AMQP_GET_BATCH_SIZE"])
	assert.Equal(t, "16", v["CONNECTORS_AMQP_DEAD_LETTER_MAX_HOPS"])
	assert.Equal(t, "0", v["CONNECTORS_AMQP_MAX_RECEIVE_COUNT"])
}

func TestAmqpConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&AmqpConfig{}).SetConfig(cfg)
	// present-empty block always emits ENABLE=false (opt-in: nil → off)
	assert.Len(t, vars(cfg), 1)
	assert.Equal(t, "false", vars(cfg)["CONNECTORS_AMQP_ENABLE"])
}

func TestAmqp10Config_SetConfig_NilEnabled(t *testing.T) {
	cfg := newTestConfig()
	(&Amqp10Config{}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_AMQP10_ENABLE"])
}

func TestAmqp10Config_SetConfig_ExplicitFalse(t *testing.T) {
	cfg := newTestConfig()
	(&Amqp10Config{Enabled: boolptr(false)}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_AMQP10_ENABLE"])
}

func TestAmqp10Config_SetConfig_ExplicitTrue(t *testing.T) {
	cfg := newTestConfig()
	(&Amqp10Config{Enabled: boolptr(true)}).SetConfig(cfg)
	v := vars(cfg)
	assert.Equal(t, "true", v["CONNECTORS_AMQP10_ENABLE"])
}

func TestAmqp10Config_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&Amqp10Config{
		Enabled:                  boolptr(true),
		Port:                     ptr32(5672),
		TLSPort:                  ptr32(5671),
		MaxFrameSize:             ptr32(131072),
		MaxMessageSize:           ptr64(104857600),
		SessionMax:               ptr32(256),
		MaxLinksPerSession:       ptr32(256),
		MaxConnections:           ptr32(1000),
		IdleTimeoutSeconds:       ptr32(120),
		DefaultPattern:           strptr("queues"),
		GetBatchSize:             ptr32(32),
		MaxUnsettledPerLink:      ptr32(1024),
		DefaultRPCTimeoutSeconds: ptr32(30),
		RPCMaxPending:            ptr32(512),
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 14) // ENABLE + 13 scalar fields
	assert.Equal(t, "true", v["CONNECTORS_AMQP10_ENABLE"])
	assert.Equal(t, "5672", v["CONNECTORS_AMQP10_PORT"])
	assert.Equal(t, "5671", v["CONNECTORS_AMQP10_TLS_PORT"])
	assert.Equal(t, "131072", v["CONNECTORS_AMQP10_MAX_FRAME_SIZE"])
	assert.Equal(t, "104857600", v["CONNECTORS_AMQP10_MAX_MESSAGE_SIZE"])
	assert.Equal(t, "256", v["CONNECTORS_AMQP10_SESSION_MAX"])
	assert.Equal(t, "256", v["CONNECTORS_AMQP10_MAX_LINKS_PER_SESSION"])
	assert.Equal(t, "1000", v["CONNECTORS_AMQP10_MAX_CONNECTIONS"])
	assert.Equal(t, "120", v["CONNECTORS_AMQP10_IDLE_TIMEOUT_SECONDS"])
	assert.Equal(t, "queues", v["CONNECTORS_AMQP10_DEFAULT_PATTERN"])
	assert.Equal(t, "32", v["CONNECTORS_AMQP10_GET_BATCH_SIZE"])
	assert.Equal(t, "1024", v["CONNECTORS_AMQP10_MAX_UNSETTLED_PER_LINK"])
	assert.Equal(t, "30", v["CONNECTORS_AMQP10_DEFAULT_RPC_TIMEOUT_SECONDS"])
	assert.Equal(t, "512", v["CONNECTORS_AMQP10_RPC_MAX_PENDING"])
}

func TestAmqp10Config_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&Amqp10Config{}).SetConfig(cfg)
	// present-empty block always emits ENABLE=false (opt-in: nil → off)
	assert.Len(t, vars(cfg), 1)
	assert.Equal(t, "false", vars(cfg)["CONNECTORS_AMQP10_ENABLE"])
}

func TestStompConfig_SetConfig_NilEnabled(t *testing.T) {
	cfg := newTestConfig()
	(&StompConfig{}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_STOMP_ENABLE"])
}

func TestStompConfig_SetConfig_ExplicitFalse(t *testing.T) {
	cfg := newTestConfig()
	(&StompConfig{Enabled: boolptr(false)}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_STOMP_ENABLE"])
}

func TestStompConfig_SetConfig_ExplicitTrue(t *testing.T) {
	cfg := newTestConfig()
	(&StompConfig{Enabled: boolptr(true)}).SetConfig(cfg)
	v := vars(cfg)
	assert.Equal(t, "true", v["CONNECTORS_STOMP_ENABLE"])
}

func TestStompConfig_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&StompConfig{
		Enabled:                boolptr(true),
		Port:                   ptr32(61613),
		TLSPort:                ptr32(61614),
		DefaultPattern:         strptr("events"),
		SubBuffSize:            ptr32(100),
		MaxConnections:         ptr32(1000),
		MaxBodySize:            ptr32(104857600),
		HeartbeatMs:            ptr32(10000),
		QueueAckTimeoutSeconds: ptr32(30),
		RPCTimeoutSeconds:      ptr32(30),
		RPCMaxPending:          ptr32(1024),
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 11) // ENABLE + 10 scalar fields
	assert.Equal(t, "true", v["CONNECTORS_STOMP_ENABLE"])
	assert.Equal(t, "61613", v["CONNECTORS_STOMP_PORT"])
	assert.Equal(t, "61614", v["CONNECTORS_STOMP_TLS_PORT"])
	assert.Equal(t, "events", v["CONNECTORS_STOMP_DEFAULT_PATTERN"])
	assert.Equal(t, "100", v["CONNECTORS_STOMP_SUB_BUFF_SIZE"])
	assert.Equal(t, "1000", v["CONNECTORS_STOMP_MAX_CONNECTIONS"])
	assert.Equal(t, "104857600", v["CONNECTORS_STOMP_MAX_BODY_SIZE"])
	assert.Equal(t, "10000", v["CONNECTORS_STOMP_HEARTBEAT_MS"])
	assert.Equal(t, "30", v["CONNECTORS_STOMP_QUEUE_ACK_TIMEOUT_SECONDS"])
	assert.Equal(t, "30", v["CONNECTORS_STOMP_RPC_TIMEOUT_SECONDS"])
	assert.Equal(t, "1024", v["CONNECTORS_STOMP_RPC_MAX_PENDING"])
}

func TestStompConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&StompConfig{}).SetConfig(cfg)
	// present-empty block always emits ENABLE=false (opt-in: nil → off)
	assert.Len(t, vars(cfg), 1)
	assert.Equal(t, "false", vars(cfg)["CONNECTORS_STOMP_ENABLE"])
}

func TestAwsConfig_SetConfig_NilEnabled(t *testing.T) {
	cfg := newTestConfig()
	(&AwsConfig{}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_AWS_ENABLE"])
}

func TestAwsConfig_SetConfig_ExplicitFalse(t *testing.T) {
	cfg := newTestConfig()
	(&AwsConfig{Enabled: boolptr(false)}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_AWS_ENABLE"])
}

func TestAwsConfig_SetConfig_ExplicitTrue(t *testing.T) {
	cfg := newTestConfig()
	(&AwsConfig{Enabled: boolptr(true)}).SetConfig(cfg)
	v := vars(cfg)
	assert.Equal(t, "true", v["CONNECTORS_AWS_ENABLE"])
}

func TestAwsConfig_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&AwsConfig{
		Enabled:             boolptr(true),
		Port:                ptr32(4566),
		Region:              strptr("kubemq"),
		AccountID:           strptr("000000000000"),
		AdvertisedURL:       strptr("http://localhost:4566"),
		MaxInflightPerQueue: ptr32(20000),
		MaxConcurrentPolls:  ptr32(1024),
		ReadTimeout:         ptr32(60),
		BodyLimit:           strptr("2M"),
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 9) // ENABLE + 8 scalar fields
	assert.Equal(t, "true", v["CONNECTORS_AWS_ENABLE"])
	assert.Equal(t, "4566", v["CONNECTORS_AWS_PORT"])
	assert.Equal(t, "kubemq", v["CONNECTORS_AWS_REGION"])
	assert.Equal(t, "000000000000", v["CONNECTORS_AWS_ACCOUNT_ID"])
	assert.Equal(t, "http://localhost:4566", v["CONNECTORS_AWS_ADVERTISED_URL"])
	assert.Equal(t, "20000", v["CONNECTORS_AWS_MAX_INFLIGHT_PER_QUEUE"])
	assert.Equal(t, "1024", v["CONNECTORS_AWS_MAX_CONCURRENT_POLLS"])
	assert.Equal(t, "60", v["CONNECTORS_AWS_READ_TIMEOUT"])
	assert.Equal(t, "2M", v["CONNECTORS_AWS_BODY_LIMIT"])
}

func TestAwsConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&AwsConfig{}).SetConfig(cfg)
	// present-empty block always emits ENABLE=false (opt-in: nil → off)
	assert.Len(t, vars(cfg), 1)
	assert.Equal(t, "false", vars(cfg)["CONNECTORS_AWS_ENABLE"])
}

func TestAwsConfig_SetConfig_SigningFields(t *testing.T) {
	cfg := newTestConfig()
	(&AwsConfig{
		Enabled:             boolptr(true),
		MessageSigning:      boolptr(true),
		SigningCertTtlHours: ptr32(24),
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 3) // ENABLE + 2 signing fields
	assert.Equal(t, "true", v["CONNECTORS_AWS_ENABLE"])
	assert.Equal(t, "true", v["CONNECTORS_AWS_MESSAGE_SIGNING"])
	assert.Equal(t, "24", v["CONNECTORS_AWS_SIGNING_CERT_TTL_HOURS"])
}

// CredentialsData must land in the cluster Secret (base64-encoded), never in the ConfigMap.
func TestAwsConfig_SetConfig_CredentialsData_Secret(t *testing.T) {
	cfg := newTestConfig()
	const raw = `{"accessKeyId":"AKIA","secretAccessKey":"shh"}`
	(&AwsConfig{Enabled: boolptr(true), CredentialsData: strptr(raw)}).SetConfig(cfg)

	// Value lands in the Secret, base64-encoded under the uppercased key.
	sd := secData(cfg)
	require.Len(t, sd, 1)
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(raw)), sd["CONNECTORS_AWS_CREDENTIALS_DATA"])

	// ConfigMap gets only ENABLE=true (never the credentials).
	v := vars(cfg)
	_, inConfigMap := v["CONNECTORS_AWS_CREDENTIALS_DATA"]
	assert.False(t, inConfigMap, "credentials must not be written to the ConfigMap")
	assert.Len(t, v, 1)
	assert.Equal(t, "true", v["CONNECTORS_AWS_ENABLE"])
}

func TestGcpConfig_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&GcpConfig{
		Enabled:                    boolptr(true),
		Port:                       ptr32(8085),
		AdvertisedEndpoint:         strptr("pubsub.example.com:443"),
		MaxMessageBytes:            ptr32(10485760),
		DefaultAckDeadlineSeconds:  ptr32(60),
		MaxOutstandingMessages:     ptr32(1000),
		MaxInflightPerSubscription: ptr32(2000),
		MaxConcurrentPolls:         ptr32(16),
		DeliveryShards:             ptr32(32),
		MaxAckExtensionSeconds:     ptr32(600),
		StreamCloseSeconds:         ptr32(30),
		MaxSeekReplay:              ptr32(100000),
		EnableReflection:           boolptr(true),
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 13) // ENABLE + 12 non-enable fields
	assert.Equal(t, "true", v["CONNECTORS_GCP_ENABLE"])
	assert.Equal(t, "8085", v["CONNECTORS_GCP_PORT"])
	assert.Equal(t, "pubsub.example.com:443", v["CONNECTORS_GCP_ADVERTISED_ENDPOINT"])
	assert.Equal(t, "10485760", v["CONNECTORS_GCP_MAX_MESSAGE_BYTES"])
	assert.Equal(t, "60", v["CONNECTORS_GCP_DEFAULT_ACK_DEADLINE_SECONDS"])
	assert.Equal(t, "1000", v["CONNECTORS_GCP_MAX_OUTSTANDING_MESSAGES"])
	assert.Equal(t, "2000", v["CONNECTORS_GCP_MAX_INFLIGHT_PER_SUBSCRIPTION"])
	assert.Equal(t, "16", v["CONNECTORS_GCP_MAX_CONCURRENT_POLLS"])
	assert.Equal(t, "32", v["CONNECTORS_GCP_DELIVERY_SHARDS"])
	assert.Equal(t, "600", v["CONNECTORS_GCP_MAX_ACK_EXTENSION_SECONDS"])
	assert.Equal(t, "30", v["CONNECTORS_GCP_STREAM_CLOSE_SECONDS"])
	assert.Equal(t, "100000", v["CONNECTORS_GCP_MAX_SEEK_REPLAY"])
	assert.Equal(t, "true", v["CONNECTORS_GCP_ENABLE_REFLECTION"])
}

func TestGcpConfig_SetConfig_NilEnabled(t *testing.T) {
	cfg := newTestConfig()
	(&GcpConfig{}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_GCP_ENABLE"])
}

func TestGcpConfig_SetConfig_ExplicitFalse(t *testing.T) {
	cfg := newTestConfig()
	(&GcpConfig{Enabled: boolptr(false)}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_GCP_ENABLE"])
}

func TestGcpConfig_SetConfig_ExplicitTrue(t *testing.T) {
	cfg := newTestConfig()
	(&GcpConfig{Enabled: boolptr(true)}).SetConfig(cfg)
	v := vars(cfg)
	assert.Equal(t, "true", v["CONNECTORS_GCP_ENABLE"])
}

func TestGcpConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&GcpConfig{}).SetConfig(cfg)
	// present-empty block always emits ENABLE=false (opt-in: nil → off)
	assert.Len(t, vars(cfg), 1)
	assert.Equal(t, "false", vars(cfg)["CONNECTORS_GCP_ENABLE"])
}

func TestApiConfig_SetConfig_Disabled(t *testing.T) {
	cfg := newTestConfig()
	(&ApiConfig{Disabled: true}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["API_ENABLE"])
}

func TestApiConfig_SetConfig_PortMoves(t *testing.T) {
	cfg := newTestConfig()
	(&ApiConfig{Port: ptr32(9000)}).SetConfig(cfg)
	assert.Equal(t, "9000", vars(cfg)["API_PORT"])
}

func TestApiConfig_SetConfig_PortUnset(t *testing.T) {
	cfg := newTestConfig()
	(&ApiConfig{}).SetConfig(cfg)
	_, has := vars(cfg)["API_PORT"]
	assert.False(t, has, "API_PORT must not be emitted when Port is unset")
}

func TestGrpcConfig_SetConfig_Disabled(t *testing.T) {
	cfg := newTestConfig()
	(&GrpcConfig{Disabled: true}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_GRPC_ENABLE"])
}

func TestGrpcConfig_SetConfig_PortMoves(t *testing.T) {
	cfg := newTestConfig()
	(&GrpcConfig{Port: ptr32(9000)}).SetConfig(cfg)
	assert.Equal(t, "9000", vars(cfg)["CONNECTORS_GRPC_PORT"])
}

func TestRestConfig_SetConfig_Disabled(t *testing.T) {
	cfg := newTestConfig()
	(&RestConfig{Disabled: true}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["CONNECTORS_REST_ENABLE"])
}

func TestRestConfig_SetConfig_PortMoves(t *testing.T) {
	cfg := newTestConfig()
	(&RestConfig{Port: ptr32(9000)}).SetConfig(cfg)
	assert.Equal(t, "9000", vars(cfg)["CONNECTORS_REST_PORT"])
}

// DeepCopy independence: mutating the source after copy must not affect the copy.
func TestConnectorConfigs_DeepCopy_Independent(t *testing.T) {
	mcp := &McpConfig{ToolTimeoutSeconds: ptr32(30), TrustedOrigins: []string{"a", "b"}}
	mcpCopy := mcp.DeepCopy()
	*mcp.ToolTimeoutSeconds = 99
	mcp.TrustedOrigins[0] = "mutated"
	assert.Equal(t, int32(30), *mcpCopy.ToolTimeoutSeconds)
	assert.Equal(t, "a", mcpCopy.TrustedOrigins[0])

	ag := &AgentsConfig{AgentMaxResponseBytes: ptr64(1), TrustedOrigins: []string{"x"}}
	agCopy := ag.DeepCopy()
	*ag.AgentMaxResponseBytes = 2
	ag.TrustedOrigins[0] = "y"
	assert.Equal(t, int64(1), *agCopy.AgentMaxResponseBytes)
	assert.Equal(t, "x", agCopy.TrustedOrigins[0])

	ce := &CeConfig{SubBuffSize: ptr32(100)}
	ceCopy := ce.DeepCopy()
	*ce.SubBuffSize = 1
	assert.Equal(t, int32(100), *ceCopy.SubBuffSize)

	mqtt := &MqttConfig{
		Port:           ptr32(1883),
		DefaultPattern: strptr("events"),
		Capabilities: &MqttCapabilitiesConfig{
			MaxClients: ptr64(10),
			MaxQos:     ptr32(2),
		},
	}
	mqttCopy := mqtt.DeepCopy()
	*mqtt.Port = 9999
	*mqtt.DefaultPattern = "mutated"
	*mqtt.Capabilities.MaxClients = 999 // mutate nested pointer field
	*mqtt.Capabilities.MaxQos = 0       // mutate nested pointer field
	mqtt.Capabilities = nil             // detach source nested struct entirely
	assert.Equal(t, int32(1883), *mqttCopy.Port)
	assert.Equal(t, "events", *mqttCopy.DefaultPattern)
	require.NotNil(t, mqttCopy.Capabilities)
	assert.Equal(t, int64(10), *mqttCopy.Capabilities.MaxClients)
	assert.Equal(t, int32(2), *mqttCopy.Capabilities.MaxQos)

	amqp := &AmqpConfig{Port: ptr32(5672), DefaultVhost: strptr("default")}
	amqpCopy := amqp.DeepCopy()
	*amqp.Port = 1
	*amqp.DefaultVhost = "mutated"
	assert.Equal(t, int32(5672), *amqpCopy.Port)
	assert.Equal(t, "default", *amqpCopy.DefaultVhost)

	amqp10 := &Amqp10Config{Port: ptr32(5672), MaxMessageSize: ptr64(104857600)}
	amqp10Copy := amqp10.DeepCopy()
	*amqp10.Port = 1
	*amqp10.MaxMessageSize = 2
	assert.Equal(t, int32(5672), *amqp10Copy.Port)
	assert.Equal(t, int64(104857600), *amqp10Copy.MaxMessageSize)

	stomp := &StompConfig{Port: ptr32(61613), DefaultPattern: strptr("events")}
	stompCopy := stomp.DeepCopy()
	*stomp.Port = 1
	*stomp.DefaultPattern = "mutated"
	assert.Equal(t, int32(61613), *stompCopy.Port)
	assert.Equal(t, "events", *stompCopy.DefaultPattern)

	aws := &AwsConfig{Port: ptr32(4566), Region: strptr("kubemq")}
	awsCopy := aws.DeepCopy()
	*aws.Port = 1
	*aws.Region = "mutated"
	assert.Equal(t, int32(4566), *awsCopy.Port)
	assert.Equal(t, "kubemq", *awsCopy.Region)

	pubsub := &GcpConfig{
		Port:               ptr32(8085),
		AdvertisedEndpoint: strptr("pubsub.example.com:443"),
		EnableReflection:   boolptr(true),
	}
	pubsubCopy := pubsub.DeepCopy()
	*pubsub.Port = 1
	*pubsub.AdvertisedEndpoint = "mutated"
	*pubsub.EnableReflection = false
	assert.Equal(t, int32(8085), *pubsubCopy.Port)
	assert.Equal(t, "pubsub.example.com:443", *pubsubCopy.AdvertisedEndpoint)
	assert.Equal(t, true, *pubsubCopy.EnableReflection)
}
