package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStoreConfig_DeepCopy_Engine guards the hand-written DeepCopy against a missed
// Engine copy (F4): a shallow copy would alias the original's *string, so a later
// mutation of the source would silently flip a "next" cluster back to legacy. The
// copy must own an independent allocation and nil must round-trip to nil.
func TestStoreConfig_DeepCopy_Engine(t *testing.T) {
	orig := &StoreConfig{Engine: strptr("next")}
	cp := orig.DeepCopy()

	require.NotNil(t, cp.Engine, "DeepCopy must copy a non-nil Engine")
	assert.Equal(t, "next", *cp.Engine)

	// Independent allocation: the copy must not alias the original's pointer.
	assert.NotSame(t, orig.Engine, cp.Engine, "DeepCopy must allocate a new *string for Engine")

	// Mutating the original must not bleed into the copy.
	*orig.Engine = "legacy"
	assert.Equal(t, "next", *cp.Engine, "copy Engine must be independent of the original")

	// A nil Engine must round-trip to nil (no spurious allocation).
	nilCopy := (&StoreConfig{}).DeepCopy()
	assert.Nil(t, nilCopy.Engine, "nil Engine must copy to nil")
}

// TestStoreConfig_SetConfig_Engine pins the STORE_ENGINE wire contract (ADD-004,
// load-bearing): only engine=next emits the key, and it emits the literal value
// "next". Unset and "legacy" emit nothing so an existing legacy CR's ConfigMap is
// byte-identical (stable CHECKSUM, no roll). The literal "STORE_ENGINE" key here is
// the exact env-var name the server binds — a rename must break this test.
func TestStoreConfig_SetConfig_Engine(t *testing.T) {
	tests := []struct {
		name    string
		engine  *string
		wantKey bool
	}{
		{name: "nil engine emits nothing", engine: nil, wantKey: false},
		{name: "legacy emits nothing", engine: strptr("legacy"), wantKey: false},
		{name: "next emits STORE_ENGINE", engine: strptr("next"), wantKey: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newTestConfig()
			(&StoreConfig{Engine: tt.engine}).SetConfig(cfg)

			got, ok := vars(cfg)["STORE_ENGINE"]
			if tt.wantKey {
				require.True(t, ok, "STORE_ENGINE must be emitted for engine=next")
				assert.Equal(t, "next", got, "STORE_ENGINE wire-contract value")
			} else {
				assert.False(t, ok, "STORE_ENGINE must not be emitted for %s", tt.name)
			}
		})
	}
}
