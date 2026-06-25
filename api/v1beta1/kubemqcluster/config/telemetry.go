package config

import (
	"strconv"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

// TelemetryTracesConfig configures OpenTelemetry tracing for the kubemq-server.
type TelemetryTracesConfig struct {
	// +optional
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	SamplingRatio *float64 `json:"samplingRatio,omitempty" yaml:"samplingRatio,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=always_on;always_off;trace_id_ratio;parent_based
	Sampler *string `json:"sampler,omitempty" yaml:"sampler,omitempty"`
}

func (c *TelemetryTracesConfig) DeepCopy() *TelemetryTracesConfig {
	out := &TelemetryTracesConfig{}

	if c.Enable != nil {
		out.Enable = new(bool)
		*out.Enable = *c.Enable
	}

	if c.SamplingRatio != nil {
		out.SamplingRatio = new(float64)
		*out.SamplingRatio = *c.SamplingRatio
	}

	if c.Sampler != nil {
		out.Sampler = new(string)
		*out.Sampler = *c.Sampler
	}

	return out
}

// TelemetryMetricsConfig configures OpenTelemetry metrics for the kubemq-server.
type TelemetryMetricsConfig struct {
	// +optional
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`

	// +optional
	ExportInterval *string `json:"exportInterval,omitempty" yaml:"exportInterval,omitempty"`
}

func (c *TelemetryMetricsConfig) DeepCopy() *TelemetryMetricsConfig {
	out := &TelemetryMetricsConfig{}

	if c.Enable != nil {
		out.Enable = new(bool)
		*out.Enable = *c.Enable
	}

	if c.ExportInterval != nil {
		out.ExportInterval = new(string)
		*out.ExportInterval = *c.ExportInterval
	}

	return out
}

// TelemetryExporterConfig configures the OpenTelemetry OTLP exporter.
type TelemetryExporterConfig struct {
	// +optional
	// +kubebuilder:validation:Enum=grpc;http
	Protocol *string `json:"protocol,omitempty" yaml:"protocol,omitempty"`

	// +optional
	Endpoint *string `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`

	// +optional
	Insecure *bool `json:"insecure,omitempty" yaml:"insecure,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=gzip;none
	Compression *string `json:"compression,omitempty" yaml:"compression,omitempty"`

	// +optional
	Timeout *string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (c *TelemetryExporterConfig) DeepCopy() *TelemetryExporterConfig {
	out := &TelemetryExporterConfig{}

	if c.Protocol != nil {
		out.Protocol = new(string)
		*out.Protocol = *c.Protocol
	}

	if c.Endpoint != nil {
		out.Endpoint = new(string)
		*out.Endpoint = *c.Endpoint
	}

	if c.Insecure != nil {
		out.Insecure = new(bool)
		*out.Insecure = *c.Insecure
	}

	if c.Compression != nil {
		out.Compression = new(string)
		*out.Compression = *c.Compression
	}

	if c.Timeout != nil {
		out.Timeout = new(string)
		*out.Timeout = *c.Timeout
	}

	return out
}

// TelemetryConfig configures kubemq-server OpenTelemetry traces & metrics.
// Maps to server Telemetry. The headers/resource maps are intentionally NOT
// exposed here — they are deferred to config.yaml.
type TelemetryConfig struct {
	// +optional
	Enable *bool `json:"enable,omitempty" yaml:"enable,omitempty"`

	// +optional
	ServiceName *string `json:"serviceName,omitempty" yaml:"serviceName,omitempty"`

	// +optional
	Traces *TelemetryTracesConfig `json:"traces,omitempty" yaml:"traces,omitempty"`

	// +optional
	Metrics *TelemetryMetricsConfig `json:"metrics,omitempty" yaml:"metrics,omitempty"`

	// +optional
	Exporter *TelemetryExporterConfig `json:"exporter,omitempty" yaml:"exporter,omitempty"`
}

func (c *TelemetryConfig) DeepCopy() *TelemetryConfig {
	out := &TelemetryConfig{}

	if c.Enable != nil {
		out.Enable = new(bool)
		*out.Enable = *c.Enable
	}

	if c.ServiceName != nil {
		out.ServiceName = new(string)
		*out.ServiceName = *c.ServiceName
	}

	if c.Traces != nil {
		out.Traces = c.Traces.DeepCopy()
	}

	if c.Metrics != nil {
		out.Metrics = c.Metrics.DeepCopy()
	}

	if c.Exporter != nil {
		out.Exporter = c.Exporter.DeepCopy()
	}

	return out
}

func (c *TelemetryConfig) SetConfig(config *deployment.Config) *TelemetryConfig {
	if c.Enable != nil {
		config.SetConfigMapStringValues(config.Name, "TELEMETRY_ENABLE", strconv.FormatBool(*c.Enable))
	}

	if c.ServiceName != nil {
		config.SetConfigMapStringValues(config.Name, "TELEMETRY_SERVICE_NAME", *c.ServiceName)
	}

	if c.Traces != nil {
		if c.Traces.Enable != nil {
			config.SetConfigMapStringValues(config.Name, "TELEMETRY_TRACES_ENABLE", strconv.FormatBool(*c.Traces.Enable))
		}

		if c.Traces.SamplingRatio != nil {
			config.SetConfigMapStringValues(config.Name, "TELEMETRY_TRACES_SAMPLING_RATIO", strconv.FormatFloat(*c.Traces.SamplingRatio, 'f', -1, 64))
		}

		if c.Traces.Sampler != nil {
			config.SetConfigMapStringValues(config.Name, "TELEMETRY_TRACES_SAMPLER", *c.Traces.Sampler)
		}
	}

	if c.Metrics != nil {
		if c.Metrics.Enable != nil {
			config.SetConfigMapStringValues(config.Name, "TELEMETRY_METRICS_ENABLE", strconv.FormatBool(*c.Metrics.Enable))
		}

		if c.Metrics.ExportInterval != nil {
			config.SetConfigMapStringValues(config.Name, "TELEMETRY_METRICS_EXPORT_INTERVAL", *c.Metrics.ExportInterval)
		}
	}

	if c.Exporter != nil {
		if c.Exporter.Protocol != nil {
			config.SetConfigMapStringValues(config.Name, "TELEMETRY_EXPORTER_PROTOCOL", *c.Exporter.Protocol)
		}

		if c.Exporter.Endpoint != nil {
			config.SetConfigMapStringValues(config.Name, "TELEMETRY_EXPORTER_ENDPOINT", *c.Exporter.Endpoint)
		}

		if c.Exporter.Insecure != nil {
			config.SetConfigMapStringValues(config.Name, "TELEMETRY_EXPORTER_INSECURE", strconv.FormatBool(*c.Exporter.Insecure))
		}

		if c.Exporter.Compression != nil {
			config.SetConfigMapStringValues(config.Name, "TELEMETRY_EXPORTER_COMPRESSION", *c.Exporter.Compression)
		}

		if c.Exporter.Timeout != nil {
			config.SetConfigMapStringValues(config.Name, "TELEMETRY_EXPORTER_TIMEOUT", *c.Exporter.Timeout)
		}
	}

	return c
}
