package config

import (
	"testing"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConnectorConfig_SetConfig_Expose pins F5 (DEC-1): all 7 wire connectors
// (kafka, mqtt, amqp, amqp10, stomp, aws, gcp) share the same Expose *string
// -> svc.Expose mechanism via SetConfig. A nil Expose leaves the catalog
// default ("ClusterIP") untouched; an explicit "LoadBalancer"/"NodePort"
// value propagates onto the K8s Service. DEC-1: no pinned NodePort field on
// any connector (aws/gcp included) — expose-only, driven through the shared
// ServiceConfig.Expose field the multi-port template renders as
// `type: {{.Expose}}`.
func TestConnectorConfig_SetConfig_Expose(t *testing.T) {
	tests := []struct {
		name      string
		svcKey    string
		setConfig func(cfg *deployment.Config, expose *string)
	}{
		{
			name:   "kafka",
			svcKey: "kafka",
			setConfig: func(cfg *deployment.Config, expose *string) {
				(&KafkaConfig{Enabled: boolptr(true), Expose: expose}).SetConfig(cfg)
			},
		},
		{
			name:   "mqtt",
			svcKey: "mqtt",
			setConfig: func(cfg *deployment.Config, expose *string) {
				(&MqttConfig{Enabled: boolptr(true), Expose: expose}).SetConfig(cfg)
			},
		},
		{
			name:   "amqp",
			svcKey: "amqp",
			setConfig: func(cfg *deployment.Config, expose *string) {
				(&AmqpConfig{Enabled: boolptr(true), Expose: expose}).SetConfig(cfg)
			},
		},
		{
			// AMQP 1.0 shares the "amqp" Service with AMQP 0.9.1 — see the
			// dedicated shared-mechanism subtests below for last-writer-wins.
			name:   "amqp10",
			svcKey: "amqp",
			setConfig: func(cfg *deployment.Config, expose *string) {
				(&Amqp10Config{Enabled: boolptr(true), Expose: expose}).SetConfig(cfg)
			},
		},
		{
			name:   "stomp",
			svcKey: "stomp",
			setConfig: func(cfg *deployment.Config, expose *string) {
				(&StompConfig{Enabled: boolptr(true), Expose: expose}).SetConfig(cfg)
			},
		},
		{
			// DEC-1: aws gets Expose ONLY, no pinned NodePort field (renders
			// via the multi-port template which has no nodePort line).
			name:   "aws",
			svcKey: "aws",
			setConfig: func(cfg *deployment.Config, expose *string) {
				(&AwsConfig{Enabled: boolptr(true), Expose: expose}).SetConfig(cfg)
			},
		},
		{
			// DEC-1: gcp gets Expose ONLY, no pinned NodePort field (renders
			// via the multi-port template which has no nodePort line).
			name:   "gcp",
			svcKey: "gcp",
			setConfig: func(cfg *deployment.Config, expose *string) {
				(&GcpConfig{Enabled: boolptr(true), Expose: expose}).SetConfig(cfg)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("nil Expose keeps default ClusterIP", func(t *testing.T) {
				cfg := newTestConfig()
				tt.setConfig(cfg, nil)
				svc, ok := cfg.Services[tt.svcKey]
				require.Truef(t, ok, "service %q not found", tt.svcKey)
				assert.Equal(t, "ClusterIP", svc.Expose)
			})

			t.Run("LoadBalancer propagates", func(t *testing.T) {
				cfg := newTestConfig()
				tt.setConfig(cfg, strptr("LoadBalancer"))
				svc, ok := cfg.Services[tt.svcKey]
				require.Truef(t, ok, "service %q not found", tt.svcKey)
				assert.Equal(t, "LoadBalancer", svc.Expose)
			})

			t.Run("NodePort propagates", func(t *testing.T) {
				cfg := newTestConfig()
				tt.setConfig(cfg, strptr("NodePort"))
				svc, ok := cfg.Services[tt.svcKey]
				require.Truef(t, ok, "service %q not found", tt.svcKey)
				assert.Equal(t, "NodePort", svc.Expose)
			})
		})
	}
}

// TestConnectorConfig_SetConfig_Expose_SharedAmqpService documents F5's
// last-writer-wins semantics on the "amqp" Service shared by AmqpConfig
// (AMQP 0.9.1) and Amqp10Config (AMQP 1.0): whichever connector's SetConfig
// call runs last determines the final Expose value on the shared Service.
func TestConnectorConfig_SetConfig_Expose_SharedAmqpService(t *testing.T) {
	t.Run("amqp10 runs last -> amqp10's Expose wins", func(t *testing.T) {
		cfg := newTestConfig()
		(&AmqpConfig{Enabled: boolptr(true), Expose: strptr("LoadBalancer")}).SetConfig(cfg)
		(&Amqp10Config{Enabled: boolptr(true), Expose: strptr("NodePort")}).SetConfig(cfg)

		svc, ok := cfg.Services["amqp"]
		require.True(t, ok, "amqp service not found")
		assert.Equal(t, "NodePort", svc.Expose, "the last SetConfig call wins on the shared Service")
	})

	t.Run("amqp runs last -> amqp's Expose wins", func(t *testing.T) {
		cfg := newTestConfig()
		(&Amqp10Config{Enabled: boolptr(true), Expose: strptr("NodePort")}).SetConfig(cfg)
		(&AmqpConfig{Enabled: boolptr(true), Expose: strptr("LoadBalancer")}).SetConfig(cfg)

		svc, ok := cfg.Services["amqp"]
		require.True(t, ok, "amqp service not found")
		assert.Equal(t, "LoadBalancer", svc.Expose, "the last SetConfig call wins on the shared Service")
	})
}
