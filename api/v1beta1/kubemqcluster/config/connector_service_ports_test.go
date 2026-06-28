package config

import (
	"testing"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A custom connector port must be reflected onto the multi-port K8s Service
// (both Port and TargetPort) so traffic through the ClusterIP Service reaches
// the listener — the connector also emits its *_PORT env, moving the server.

func svcPort(cfg *deployment.Config, svc, name string) (int32, int32, bool) {
	s, ok := cfg.Services[svc]
	if !ok {
		return 0, 0, false
	}
	for _, p := range s.Ports {
		if p.Name == name {
			return p.Port, p.TargetPort, true
		}
	}
	return 0, 0, false
}

func assertSvcPort(t *testing.T, cfg *deployment.Config, svc, name string, want int32) {
	t.Helper()
	port, target, ok := svcPort(cfg, svc, name)
	require.Truef(t, ok, "service %q port %q not found", svc, name)
	assert.Equalf(t, want, port, "service %q port %q .port", svc, name)
	assert.Equalf(t, want, target, "service %q port %q .targetPort", svc, name)
}

// assertScalarSvcPort checks the single-port Services (api/grpc/rest) whose port
// lives on the scalar ServiceConfig.ContainerPort/.TargetPort fields rather than
// in the multi-port .Ports slice.
func assertScalarSvcPort(t *testing.T, cfg *deployment.Config, svc string, want int32) {
	t.Helper()
	s, ok := cfg.Services[svc]
	require.Truef(t, ok, "service %q", svc)
	assert.Equalf(t, want, s.ContainerPort, "%q containerPort", svc)
	assert.Equalf(t, want, s.TargetPort, "%q targetPort", svc)
}

func TestMqttConfig_SetConfig_ServicePorts(t *testing.T) {
	cfg := newTestConfig()
	// Enabled: true required to reach the port-reflection path.
	(&MqttConfig{Enabled: boolptr(true), Port: ptr32(1884), TLSPort: ptr32(8884), WSPort: ptr32(8084)}).SetConfig(cfg)
	assertSvcPort(t, cfg, "mqtt", "mqtt", 1884)
	assertSvcPort(t, cfg, "mqtt", "mqtt-tls", 8884)
	assertSvcPort(t, cfg, "mqtt", "mqtt-ws", 8084)
}

func TestAmqpConfig_SetConfig_ServicePorts(t *testing.T) {
	cfg := newTestConfig()
	// Enabled: true required to reach the port-reflection path.
	(&AmqpConfig{Enabled: boolptr(true), Port: ptr32(5673), TLSPort: ptr32(5674)}).SetConfig(cfg)
	assertSvcPort(t, cfg, "amqp", "amqp", 5673)
	assertSvcPort(t, cfg, "amqp", "amqp-tls", 5674)
}

func TestAmqp10Config_SetConfig_ServicePorts(t *testing.T) {
	cfg := newTestConfig()
	// AMQP 1.0 shares the "amqp" Service with AMQP 0.9.1.
	// Enabled: true required to reach the port-reflection path.
	(&Amqp10Config{Enabled: boolptr(true), Port: ptr32(5680), TLSPort: ptr32(5681)}).SetConfig(cfg)
	assertSvcPort(t, cfg, "amqp", "amqp", 5680)
	assertSvcPort(t, cfg, "amqp", "amqp-tls", 5681)
}

func TestStompConfig_SetConfig_ServicePorts(t *testing.T) {
	cfg := newTestConfig()
	// Enabled: true required to reach the port-reflection path.
	(&StompConfig{Enabled: boolptr(true), Port: ptr32(61615), TLSPort: ptr32(61616)}).SetConfig(cfg)
	assertSvcPort(t, cfg, "stomp", "stomp", 61615)
	assertSvcPort(t, cfg, "stomp", "stomp-tls", 61616)
}

func TestAwsConfig_SetConfig_ServicePorts(t *testing.T) {
	cfg := newTestConfig()
	// Enabled: true required to reach the port-reflection path.
	(&AwsConfig{Enabled: boolptr(true), Port: ptr32(4567)}).SetConfig(cfg)
	assertSvcPort(t, cfg, "aws", "aws-http", 4567)
}

func TestGcpConfig_SetConfig_ServicePorts(t *testing.T) {
	cfg := newTestConfig()
	// gcp is a multi-port Service (port name "gcp-grpc"), so the existing
	// slice-based helper applies. Enabled: true required to reach the port-reflection path.
	(&GcpConfig{Enabled: boolptr(true), Port: ptr32(9000)}).SetConfig(cfg)
	assertSvcPort(t, cfg, "gcp", "gcp-grpc", 9000)
}

func TestApiConfig_SetConfig_ServicePort(t *testing.T) {
	t.Run("custom port moves scalar", func(t *testing.T) {
		cfg := newTestConfig()
		(&ApiConfig{Port: ptr32(9000)}).SetConfig(cfg)
		assertScalarSvcPort(t, cfg, "api", 9000)
	})
	t.Run("default unchanged", func(t *testing.T) {
		cfg := newTestConfig()
		(&ApiConfig{}).SetConfig(cfg)
		assertScalarSvcPort(t, cfg, "api", 8080)
	})
}

func TestGrpcConfig_SetConfig_ServicePort(t *testing.T) {
	t.Run("custom port moves scalar", func(t *testing.T) {
		cfg := newTestConfig()
		(&GrpcConfig{Port: ptr32(9000)}).SetConfig(cfg)
		assertScalarSvcPort(t, cfg, "grpc", 9000)
	})
	t.Run("default unchanged", func(t *testing.T) {
		cfg := newTestConfig()
		(&GrpcConfig{}).SetConfig(cfg)
		assertScalarSvcPort(t, cfg, "grpc", 50000)
	})
}

func TestRestConfig_SetConfig_ServicePort(t *testing.T) {
	t.Run("custom port moves scalar", func(t *testing.T) {
		cfg := newTestConfig()
		(&RestConfig{Port: ptr32(9000)}).SetConfig(cfg)
		assertScalarSvcPort(t, cfg, "rest", 9000)
	})
	t.Run("default unchanged", func(t *testing.T) {
		cfg := newTestConfig()
		(&RestConfig{}).SetConfig(cfg)
		assertScalarSvcPort(t, cfg, "rest", 9090)
	})
}

func TestConnector_SetConfig_ServicePorts_DefaultUnchanged(t *testing.T) {
	cfg := newTestConfig()
	(&MqttConfig{}).SetConfig(cfg)
	assertSvcPort(t, cfg, "mqtt", "mqtt", 1883)
}
