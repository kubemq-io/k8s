package config

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// containerPortByName returns the containerPort value of the named port in the
// rendered StatefulSet's first container, plus whether a port with that name exists.
func containerPortByName(c corev1.Container, name string) (int32, bool) {
	for _, p := range c.Ports {
		if p.Name == name {
			return p.ContainerPort, true
		}
	}
	return 0, false
}

// TestStatefulSet_ApiPort_Render proves that moving the API port (FR-2) flows all
// the way through to the rendered StatefulSet: the api-port containerPort, the
// liveness probe httpGet port, and the prometheus.io/port pod annotation all follow
// the configured port, and API_BIND_ADDRESS is emitted when health is enabled.
// It also guards the P1-A arg-swap fix: the rendered probe must carry the DEFAULT
// periodSeconds=10 / timeoutSeconds=5 (before the fix these were swapped to 5/10).
func TestStatefulSet_ApiPort_Render(t *testing.T) {
	cfg := newTestConfig()
	(&ApiConfig{Port: ptr32(9000)}).SetConfig(cfg)
	(&HealthConfig{Enabled: true}).SetConfig(cfg)

	sts, err := cfg.StatefulSet.Get()
	require.NoError(t, err)
	require.NotEmpty(t, sts.Spec.Template.Spec.Containers, "expected at least one container")

	c := sts.Spec.Template.Spec.Containers[0]

	// api-port containerPort follows the configured API port.
	apiPort, ok := containerPortByName(c, "api-port")
	require.Truef(t, ok, "container port %q not found", "api-port")
	assert.Equal(t, int32(9000), apiPort, "api-port containerPort")

	// Liveness probe httpGet port follows the configured API port.
	probe := c.LivenessProbe
	require.NotNil(t, probe, "expected a liveness probe when health is enabled")
	require.NotNil(t, probe.HTTPGet, "expected an httpGet liveness probe")
	assert.Equal(t, 9000, probe.HTTPGet.Port.IntValue(), "liveness probe httpGet port")

	// P1-A regression: probe period/timeout must NOT be swapped (defaults 10/5).
	assert.Equal(t, int32(10), probe.PeriodSeconds, "probe periodSeconds (arg-swap regression)")
	assert.Equal(t, int32(5), probe.TimeoutSeconds, "probe timeoutSeconds (arg-swap regression)")

	// Prometheus scrape annotation port follows the configured API port.
	assert.Equal(t, "9000", sts.Spec.Template.Annotations["prometheus.io/port"],
		"prometheus.io/port annotation")

	// Health enabled => the API must bind 0.0.0.0 so the kubelet probe is reachable.
	assert.Equal(t, "0.0.0.0", vars(cfg)["API_BIND_ADDRESS"], "API_BIND_ADDRESS env")
}

// TestStatefulSet_DefaultPorts_Render guards the StatefulSet port defaults: with no
// Api/Grpc/Rest blocks the rendered containerPorts must be 50000/8080/9090 and the
// prometheus annotation 8080 (never 0 / an invalid pod spec).
func TestStatefulSet_DefaultPorts_Render(t *testing.T) {
	cfg := newTestConfig()

	sts, err := cfg.StatefulSet.Get()
	require.NoError(t, err)
	require.NotEmpty(t, sts.Spec.Template.Spec.Containers, "expected at least one container")

	c := sts.Spec.Template.Spec.Containers[0]

	grpcPort, ok := containerPortByName(c, "grpc-port")
	require.Truef(t, ok, "container port %q not found", "grpc-port")
	assert.Equal(t, int32(50000), grpcPort, "grpc-port containerPort default")

	apiPort, ok := containerPortByName(c, "api-port")
	require.Truef(t, ok, "container port %q not found", "api-port")
	assert.Equal(t, int32(8080), apiPort, "api-port containerPort default")

	restPort, ok := containerPortByName(c, "rest-port")
	require.Truef(t, ok, "container port %q not found", "rest-port")
	assert.Equal(t, int32(9090), restPort, "rest-port containerPort default")

	assert.Equal(t, "8080", sts.Spec.Template.Annotations["prometheus.io/port"],
		"prometheus.io/port annotation default")
}

// TestStatefulSet_HealthDisabled_NoBindAddress proves that with health disabled no
// liveness probe is rendered and API_BIND_ADDRESS is not emitted.
func TestStatefulSet_HealthDisabled_NoBindAddress(t *testing.T) {
	cfg := newTestConfig()
	(&HealthConfig{Enabled: false}).SetConfig(cfg)

	_, ok := vars(cfg)["API_BIND_ADDRESS"]
	assert.False(t, ok, "API_BIND_ADDRESS must not be set when health is disabled")

	sts, err := cfg.StatefulSet.Get()
	require.NoError(t, err)
	require.NotEmpty(t, sts.Spec.Template.Spec.Containers, "expected at least one container")

	assert.Nil(t, sts.Spec.Template.Spec.Containers[0].LivenessProbe,
		"no liveness probe expected when health is disabled")
}
