package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ptrBool and ptrFloat extend the shared test helpers (ptr32/ptr64/strptr in
// connectors_env_test.go) for the *bool and *float64 fields introduced by the
// telemetry/audit/http configs.
func ptrBool(v bool) *bool        { return &v }
func ptrFloat(v float64) *float64 { return &v }

func TestTelemetryConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&TelemetryConfig{}).SetConfig(cfg)
	assert.Len(t, vars(cfg), 0)
}

func TestTelemetryConfig_SetConfig_DisabledFalse(t *testing.T) {
	cfg := newTestConfig()
	(&TelemetryConfig{
		Enable: ptrBool(false),
		Traces: &TelemetryTracesConfig{Enable: ptrBool(false)},
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 2)
	// CRITICAL false-emission cases: explicit false MUST be written.
	assert.Equal(t, "false", v["TELEMETRY_ENABLE"])
	assert.Equal(t, "false", v["TELEMETRY_TRACES_ENABLE"])
}

func TestTelemetryConfig_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&TelemetryConfig{
		Enable:      ptrBool(true),
		ServiceName: strptr("kubemq-server"),
		Traces: &TelemetryTracesConfig{
			Enable:        ptrBool(true),
			SamplingRatio: ptrFloat(0.25),
			Sampler:       strptr("trace_id_ratio"),
		},
		Metrics: &TelemetryMetricsConfig{
			Enable:         ptrBool(true),
			ExportInterval: strptr("60s"),
		},
		Exporter: &TelemetryExporterConfig{
			Protocol:    strptr("grpc"),
			Endpoint:    strptr("otel-collector:4317"),
			Insecure:    ptrBool(true),
			Compression: strptr("gzip"),
			Timeout:     strptr("10s"),
		},
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 12)
	assert.Equal(t, "true", v["TELEMETRY_ENABLE"])
	assert.Equal(t, "kubemq-server", v["TELEMETRY_SERVICE_NAME"])
	assert.Equal(t, "true", v["TELEMETRY_TRACES_ENABLE"])
	assert.Equal(t, "0.25", v["TELEMETRY_TRACES_SAMPLING_RATIO"])
	assert.Equal(t, "trace_id_ratio", v["TELEMETRY_TRACES_SAMPLER"])
	assert.Equal(t, "true", v["TELEMETRY_METRICS_ENABLE"])
	assert.Equal(t, "60s", v["TELEMETRY_METRICS_EXPORT_INTERVAL"])
	assert.Equal(t, "grpc", v["TELEMETRY_EXPORTER_PROTOCOL"])
	assert.Equal(t, "otel-collector:4317", v["TELEMETRY_EXPORTER_ENDPOINT"])
	assert.Equal(t, "true", v["TELEMETRY_EXPORTER_INSECURE"])
	assert.Equal(t, "gzip", v["TELEMETRY_EXPORTER_COMPRESSION"])
	assert.Equal(t, "10s", v["TELEMETRY_EXPORTER_TIMEOUT"])
}

func TestTelemetryConfig_DeepCopy_Independent(t *testing.T) {
	src := &TelemetryConfig{
		Enable:      ptrBool(true),
		ServiceName: strptr("svc"),
		Traces:      &TelemetryTracesConfig{Enable: ptrBool(true), SamplingRatio: ptrFloat(0.5), Sampler: strptr("always_on")},
		Metrics:     &TelemetryMetricsConfig{Enable: ptrBool(true), ExportInterval: strptr("30s")},
		Exporter:    &TelemetryExporterConfig{Protocol: strptr("http"), Insecure: ptrBool(true)},
	}
	cp := src.DeepCopy()
	*src.Enable = false
	*src.ServiceName = "mutated"
	*src.Traces.SamplingRatio = 0.1
	*src.Metrics.Enable = false
	*src.Exporter.Insecure = false
	src.Traces = nil
	src.Metrics = nil
	src.Exporter = nil

	assert.Equal(t, true, *cp.Enable)
	assert.Equal(t, "svc", *cp.ServiceName)
	require.NotNil(t, cp.Traces)
	assert.Equal(t, 0.5, *cp.Traces.SamplingRatio)
	require.NotNil(t, cp.Metrics)
	assert.Equal(t, true, *cp.Metrics.Enable)
	require.NotNil(t, cp.Exporter)
	assert.Equal(t, true, *cp.Exporter.Insecure)
}
