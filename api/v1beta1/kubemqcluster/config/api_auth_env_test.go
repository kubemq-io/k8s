package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestApiAuthConfig_EnvNameContract pins the exact API_AUTH_* env-var names the
// server binds (config/env.go convertEnvFormat) plus the os.Getenv-read admin
// username. These names are a hard contract with kubemq-server; changing them
// here silently breaks auth provisioning, so this test must fail loudly.
func TestApiAuthConfig_EnvNameContract(t *testing.T) {
	cfg := newTestConfig()
	(&ApiConfig{
		Auth: &ApiAuthConfig{
			Enable:               boolptr(true),
			SessionIdleMinutes:   ptr32(45),
			SessionAbsoluteHours: ptr32(12),
			StorePath:            strptr("/data/auth"),
			TrustedTLSProxy:      boolptr(true),
			AdminUsername:        strptr("root"),
		},
	}).SetConfig(cfg)

	v := vars(cfg)
	assert.Equal(t, "true", v["API_AUTH_ENABLE"])
	assert.Equal(t, "45", v["API_AUTH_SESSION_IDLE_MINUTES"])
	assert.Equal(t, "12", v["API_AUTH_SESSION_ABSOLUTE_HOURS"])
	assert.Equal(t, "/data/auth", v["API_AUTH_STORE_PATH"])
	assert.Equal(t, "true", v["API_AUTH_TRUSTED_TLS_PROXY"])
	assert.Equal(t, "root", v["KUBEMQ_API_ADMIN_USERNAME"])

	// The admin PASSWORD must never be emitted into the ConfigMap (it is a
	// Secret injected into the pod by the operator via valueFrom).
	for k := range v {
		assert.NotContains(t, k, "ADMIN_PASSWORD", "admin password must not leak into the ConfigMap")
	}
}

// TestApiAuthConfig_Disabled confirms Enable=false writes the explicit "false"
// signal (so a disabled block is unambiguous on the wire).
func TestApiAuthConfig_Disabled(t *testing.T) {
	cfg := newTestConfig()
	(&ApiConfig{Auth: &ApiAuthConfig{Enable: boolptr(false)}}).SetConfig(cfg)

	v := vars(cfg)
	assert.Equal(t, "false", v["API_AUTH_ENABLE"])
}

// TestApiAuthConfig_OmittedFields ensures unset fields emit no env keys (the
// server applies its own defaults).
func TestApiAuthConfig_OmittedFields(t *testing.T) {
	cfg := newTestConfig()
	(&ApiConfig{Auth: &ApiAuthConfig{Enable: boolptr(true)}}).SetConfig(cfg)

	v := vars(cfg)
	assert.Equal(t, "true", v["API_AUTH_ENABLE"])
	_, hasIdle := v["API_AUTH_SESSION_IDLE_MINUTES"]
	_, hasAbs := v["API_AUTH_SESSION_ABSOLUTE_HOURS"]
	_, hasStore := v["API_AUTH_STORE_PATH"]
	_, hasProxy := v["API_AUTH_TRUSTED_TLS_PROXY"]
	_, hasUser := v["KUBEMQ_API_ADMIN_USERNAME"]
	assert.False(t, hasIdle)
	assert.False(t, hasAbs)
	assert.False(t, hasStore)
	assert.False(t, hasProxy)
	assert.False(t, hasUser)
}

// TestApiAuthConfig_NilAuth confirms an Api block with no Auth sub-block emits
// no auth env at all (backward-compatible with clusters that never set it).
func TestApiAuthConfig_NilAuth(t *testing.T) {
	cfg := newTestConfig()
	(&ApiConfig{Port: ptr32(8080)}).SetConfig(cfg)

	v := vars(cfg)
	_, hasEnable := v["API_AUTH_ENABLE"]
	assert.False(t, hasEnable)
}

// TestApiAuthConfig_DeepCopy_Independent verifies the deep copy is fully
// independent (mutating the source leaves the copy untouched).
func TestApiAuthConfig_DeepCopy_Independent(t *testing.T) {
	src := &ApiAuthConfig{
		Enable:               boolptr(true),
		SessionIdleMinutes:   ptr32(45),
		SessionAbsoluteHours: ptr32(12),
		StorePath:            strptr("/data/auth"),
		TrustedTLSProxy:      boolptr(true),
		AdminUsername:        strptr("root"),
		AdminSecretRef:       strptr("my-secret"),
		AdminSecretKey:       strptr("password"),
	}
	cp := src.DeepCopy()

	*src.Enable = false
	*src.SessionIdleMinutes = 0
	*src.SessionAbsoluteHours = 0
	*src.StorePath = "mutated"
	*src.TrustedTLSProxy = false
	*src.AdminUsername = "mutated"
	*src.AdminSecretRef = "mutated"
	*src.AdminSecretKey = "mutated"

	require.NotNil(t, cp.Enable)
	assert.Equal(t, true, *cp.Enable)
	assert.Equal(t, int32(45), *cp.SessionIdleMinutes)
	assert.Equal(t, int32(12), *cp.SessionAbsoluteHours)
	assert.Equal(t, "/data/auth", *cp.StorePath)
	assert.Equal(t, true, *cp.TrustedTLSProxy)
	assert.Equal(t, "root", *cp.AdminUsername)
	assert.Equal(t, "my-secret", *cp.AdminSecretRef)
	assert.Equal(t, "password", *cp.AdminSecretKey)
}

// TestApiConfig_DeepCopy_IncludesAuth verifies the parent ApiConfig deep copy
// carries the nested Auth block (regression guard for the wiring).
func TestApiConfig_DeepCopy_IncludesAuth(t *testing.T) {
	src := &ApiConfig{
		Port: ptr32(8080),
		Auth: &ApiAuthConfig{Enable: boolptr(true), AdminUsername: strptr("root")},
	}
	cp := src.DeepCopy()
	*src.Auth.Enable = false
	*src.Auth.AdminUsername = "mutated"

	require.NotNil(t, cp.Auth)
	assert.Equal(t, true, *cp.Auth.Enable)
	assert.Equal(t, "root", *cp.Auth.AdminUsername)
}
