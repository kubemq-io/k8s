package config

import (
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

func ptr32(v int32) *int32 { return &v }
func ptr64(v int64) *int64 { return &v }

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
}
