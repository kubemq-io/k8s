package config

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The server reads Authorization.Url raw and hands it to http.Get (validateURL
// rejects a base64 blob). Authorization.Policy, in contrast, is base64-decoded by
// the server. These tests pin AUTHORIZATION_URL to SetStringVariable (raw) and
// AUTHORIZATION_POLICY_DATA to SetDataVariable (base64) — a regression that swaps
// the setter on either field would ship a value the server can't parse.

func TestAuthorizationConfig_Url_EmittedRaw(t *testing.T) {
	cfg := newTestConfig()
	url := "https://authz.example.com/acl"
	(&AuthorizationConfig{Url: url}).SetConfig(cfg)

	v := vars(cfg)
	assert.Equal(t, "true", v["AUTHORIZATION_ENABLE"])
	assert.Equal(t, url, v["AUTHORIZATION_URL"], "AUTHORIZATION_URL must be raw, not base64")
}

func TestAuthorizationConfig_Policy_EmittedBase64(t *testing.T) {
	cfg := newTestConfig()
	const policy = "some-policy-json"
	(&AuthorizationConfig{Policy: policy}).SetConfig(cfg)

	v := vars(cfg)
	assert.Equal(t, "true", v["AUTHORIZATION_ENABLE"])
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(policy)), v["AUTHORIZATION_POLICY_DATA"],
		"AUTHORIZATION_POLICY_DATA must be base64-encoded")
	assert.NotEqual(t, policy, v["AUTHORIZATION_POLICY_DATA"], "must not be the raw policy string")
}

func TestAuthorizationConfig_Empty_EmitsNothing(t *testing.T) {
	cfg := newTestConfig()
	(&AuthorizationConfig{}).SetConfig(cfg)

	v := vars(cfg)
	_, hasEnable := v["AUTHORIZATION_ENABLE"]
	_, hasURL := v["AUTHORIZATION_URL"]
	_, hasPolicy := v["AUTHORIZATION_POLICY_DATA"]
	assert.False(t, hasEnable, "empty authorization must not emit AUTHORIZATION_ENABLE")
	assert.False(t, hasURL)
	assert.False(t, hasPolicy)
}

func TestAuthorizationConfig_AutoReloadZero_Omitted(t *testing.T) {
	cfg := newTestConfig()
	(&AuthorizationConfig{Url: "https://authz.example.com/acl"}).SetConfig(cfg)

	v := vars(cfg)
	_, hasReload := v["AUTHORIZATION_AUTO_RELOAD"]
	assert.False(t, hasReload, "AUTO_RELOAD must be omitted when zero")
}
