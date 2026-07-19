package v1alpha1

import (
	"testing"

	"github.com/kubemq-io/k8s/api/v1alpha1/kubemqcluster/config"
	"github.com/stretchr/testify/require"
)

// TestV1alpha1DeepCopy_Env guards the hand-maintained deepcopy of the WP-4-full typed-sink
// fields (F10.5): the v1alpha1 KubemqClusterSpec Env/EnvFromSecrets and the v1alpha1
// StoreConfig.Engine. A shallow-copy regression here would let a deepcopy alias the parent's
// map/slice/pointer and silently corrupt reconcile state.
func TestV1alpha1DeepCopy_Env(t *testing.T) {
	engine := "next"
	in := &KubemqClusterSpec{
		Env:            map[string]string{"FOO": "bar"},
		EnvFromSecrets: []string{"sec-a", "sec-b"},
		Store:          &config.StoreConfig{Engine: &engine},
	}
	out := in.DeepCopy()

	require.Equal(t, in.Env, out.Env)
	require.Equal(t, in.EnvFromSecrets, out.EnvFromSecrets)
	require.NotNil(t, out.Store)
	require.NotNil(t, out.Store.Engine)
	require.Equal(t, "next", *out.Store.Engine)

	// Mutating the copy must not touch the original (deep, not shallow).
	out.Env["FOO"] = "changed"
	out.EnvFromSecrets[0] = "changed"
	*out.Store.Engine = "legacy"
	require.Equal(t, "bar", in.Env["FOO"], "Env map aliased — deepcopy is shallow")
	require.Equal(t, "sec-a", in.EnvFromSecrets[0], "EnvFromSecrets slice aliased — deepcopy is shallow")
	require.Equal(t, "next", *in.Store.Engine, "Store.Engine pointer aliased — deepcopy is shallow")
}
