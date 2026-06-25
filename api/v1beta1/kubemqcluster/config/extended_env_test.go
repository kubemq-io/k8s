package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Env-var assertions for the fields promoted in the config-alignment change that
// did not previously have dedicated env tests. Literals are the exact server
// bindings (convertEnvFormat over kubemq-server/config bindViperEnv keys).

func TestStoreConfig_SetConfig_IdlePruneCutoffHours(t *testing.T) {
	cfg := newTestConfig()
	(&StoreConfig{IdlePruneCutoffHours: ptr32(24)}).SetConfig(cfg)
	assert.Equal(t, "24", vars(cfg)["STORE_IDLE_PRUNE_CUTOFF_HOURS"])
}

func TestQueueConfig_SetConfig_InflightAndPubAck(t *testing.T) {
	cfg := newTestConfig()
	(&QueueConfig{MaxInflight: ptr32(2048), PubAckWaitSeconds: ptr32(60)}).SetConfig(cfg)
	v := vars(cfg)
	assert.Equal(t, "2048", v["QUEUE_MAX_INFLIGHT"])
	assert.Equal(t, "60", v["QUEUE_PUB_ACK_WAIT_SECONDS"])
}

func TestGrpcConfig_SetConfig_EnableReflection(t *testing.T) {
	cfg := newTestConfig()
	(&GrpcConfig{EnableReflection: ptrBool(true)}).SetConfig(cfg)
	assert.Equal(t, "true", vars(cfg)["CONNECTORS_GRPC_ENABLE_REFLECTION"])

	// false must be emitted (a plain bool could not express this)
	cfgFalse := newTestConfig()
	(&GrpcConfig{EnableReflection: ptrBool(false)}).SetConfig(cfgFalse)
	assert.Equal(t, "false", vars(cfgFalse)["CONNECTORS_GRPC_ENABLE_REFLECTION"])
}

func TestApiConfig_SetConfig_AllowOrigins(t *testing.T) {
	cfg := newTestConfig()
	(&ApiConfig{AllowOrigins: []string{"https://a", "https://b"}}).SetConfig(cfg)
	assert.Equal(t, "https://a,https://b", vars(cfg)["API_ALLOW_ORIGINS"])
}

// Additional false-emission coverage for defaults-ON telemetry toggles (server
// defaults true) — they must be disable-able via the CRD.
func TestTelemetryConfig_SetConfig_DisableNestedToggles(t *testing.T) {
	cfg := newTestConfig()
	(&TelemetryConfig{
		Metrics:  &TelemetryMetricsConfig{Enable: ptrBool(false)},
		Exporter: &TelemetryExporterConfig{Insecure: ptrBool(false)},
	}).SetConfig(cfg)
	v := vars(cfg)
	assert.Equal(t, "false", v["TELEMETRY_METRICS_ENABLE"])
	assert.Equal(t, "false", v["TELEMETRY_EXPORTER_INSECURE"])
}
