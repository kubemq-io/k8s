package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&AuditConfig{}).SetConfig(cfg)
	assert.Len(t, vars(cfg), 0)
}

func TestAuditConfig_SetConfig_DisabledFalse(t *testing.T) {
	cfg := newTestConfig()
	// CRITICAL false-emission case: audit defaults ON server-side, so an explicit
	// false MUST be written to disable it.
	(&AuditConfig{Enable: ptrBool(false)}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 1)
	assert.Equal(t, "false", v["AUDIT_ENABLE"])
}

func TestAuditConfig_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&AuditConfig{
		Enable:                 ptrBool(true),
		RetentionHours:         ptr32(720),
		CleanupIntervalMinutes: ptr32(60),
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 3)
	assert.Equal(t, "true", v["AUDIT_ENABLE"])
	assert.Equal(t, "720", v["AUDIT_RETENTION_HOURS"])
	assert.Equal(t, "60", v["AUDIT_CLEANUP_INTERVAL_MINUTES"])
}

func TestAuditConfig_DeepCopy_Independent(t *testing.T) {
	src := &AuditConfig{Enable: ptrBool(true), RetentionHours: ptr32(720), CleanupIntervalMinutes: ptr32(60)}
	cp := src.DeepCopy()
	*src.Enable = false
	*src.RetentionHours = 1
	*src.CleanupIntervalMinutes = 1
	assert.Equal(t, true, *cp.Enable)
	assert.Equal(t, int32(720), *cp.RetentionHours)
	assert.Equal(t, int32(60), *cp.CleanupIntervalMinutes)
}
