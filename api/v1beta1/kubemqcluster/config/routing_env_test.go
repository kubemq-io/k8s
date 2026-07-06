package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The server binds Routing.Data / Routing.Url and reads them RAW (no base64
// decode). These tests pin the operator to emitting raw values — a regression to
// SetDataVariable (base64) would ship a routing table / URL the server can't parse.

func TestRoutingConfig_DataAndUrl_EmittedRaw(t *testing.T) {
	cfg := newTestConfig()
	data := `{"routes":[{"channel":"orders.*","routes":"events"}]}`
	url := "https://routing.example.com/table.json"
	(&RoutingConfig{Data: data, Url: url, AutoReload: 30}).SetConfig(cfg)

	v := vars(cfg)
	assert.Equal(t, "true", v["ROUTING_ENABLE"])
	assert.Equal(t, data, v["ROUTING_DATA"], "ROUTING_DATA must be raw, not base64")
	assert.Equal(t, url, v["ROUTING_URL"], "ROUTING_URL must be raw, not base64")
	assert.Equal(t, "30", v["ROUTING_AUTO_RELOAD"])
}

func TestRoutingConfig_Empty_EmitsNothing(t *testing.T) {
	cfg := newTestConfig()
	(&RoutingConfig{}).SetConfig(cfg)

	v := vars(cfg)
	_, hasEnable := v["ROUTING_ENABLE"]
	_, hasData := v["ROUTING_DATA"]
	_, hasURL := v["ROUTING_URL"]
	assert.False(t, hasEnable, "empty routing must not emit ROUTING_ENABLE")
	assert.False(t, hasData)
	assert.False(t, hasURL)
}

func TestRoutingConfig_AutoReloadZero_Omitted(t *testing.T) {
	cfg := newTestConfig()
	(&RoutingConfig{Data: "x"}).SetConfig(cfg)

	v := vars(cfg)
	_, hasReload := v["ROUTING_AUTO_RELOAD"]
	assert.False(t, hasReload, "AUTO_RELOAD must be omitted when zero")
}
