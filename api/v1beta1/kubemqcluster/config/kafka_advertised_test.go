package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestKafkaConfig_SetConfig_AdvertisedHostDefault covers F4: when the Kafka
// connector is enabled and AdvertisedHost is left unset, SetConfig must
// auto-emit the in-cluster short-form Service DNS name
// (<cluster>-kafka.<namespace>.svc) so external/in-cluster clients don't
// hang on reconnect (the M-23 connect-then-hang footgun). An explicit
// AdvertisedHost always wins, and a disabled/opt-in-off connector must never
// emit the key at all.
func TestKafkaConfig_SetConfig_AdvertisedHostDefault(t *testing.T) {
	tests := []struct {
		name           string
		enabled        *bool
		advertisedHost *string
		wantPresent    bool
		wantValue      string
	}{
		{
			name:        "unset advertisedHost + enabled -> .svc default",
			enabled:     boolptr(true),
			wantPresent: true,
			wantValue:   testClusterName + "-kafka.kubemq.svc",
		},
		{
			name:           "explicit advertisedHost wins over the default",
			enabled:        boolptr(true),
			advertisedHost: strptr("kafka.example.com"),
			wantPresent:    true,
			wantValue:      "kafka.example.com",
		},
		{
			name:        "explicitly disabled -> key absent",
			enabled:     boolptr(false),
			wantPresent: false,
		},
		{
			name:        "nil Enabled (opt-in off) -> key absent",
			enabled:     nil,
			wantPresent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newTestConfig()
			(&KafkaConfig{Enabled: tt.enabled, AdvertisedHost: tt.advertisedHost}).SetConfig(cfg)

			v := vars(cfg)
			got, present := v["CONNECTORS_KAFKA_ADVERTISED_HOST"]
			assert.Equal(t, tt.wantPresent, present, "CONNECTORS_KAFKA_ADVERTISED_HOST presence")
			if tt.wantPresent {
				assert.Equal(t, tt.wantValue, got)
			}
		})
	}
}
