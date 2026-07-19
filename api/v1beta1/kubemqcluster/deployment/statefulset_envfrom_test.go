package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStatefulSet_EnvFrom_ExtraSecretRefs asserts SetExtraSecretRefs appends
// additional secretRef envFrom entries AFTER the base secretRef+configMapRef
// pair, in the order supplied — so user-provided secret keys win over
// ConfigMap-sourced defaults. An empty/nil slice is parse-neutral: only the
// base secretRef+configMapRef pair renders.
func TestStatefulSet_EnvFrom_ExtraSecretRefs(t *testing.T) {
	tests := []struct {
		name            string
		extraSecretRefs []string
		wantExtra       []string // expected trailing secretRef names, after the base pair; nil == none
	}{
		{
			name:            "two_extra_secret_refs_in_order",
			extraSecretRefs: []string{"s1", "s2"},
			wantExtra:       []string{"s1", "s2"},
		},
		{
			name:            "empty_slice_is_parse_neutral",
			extraSecretRefs: []string{},
			wantExtra:       nil,
		},
		{
			name:            "nil_slice_is_parse_neutral",
			extraSecretRefs: nil,
			wantExtra:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultStatefulSetConfig("", "kubemq-cluster", "kubemq").
				SetExtraSecretRefs(tt.extraSecretRefs)

			sts, err := cfg.Get()
			require.NoError(t, err)
			require.Len(t, sts.Spec.Template.Spec.Containers, 1)

			c := sts.Spec.Template.Spec.Containers[0]
			require.Len(t, c.EnvFrom, len(tt.wantExtra)+2, "expected base secretRef+configMapRef plus %d extra secretRef(s)", len(tt.wantExtra))

			// First entry is always the base secretRef.
			require.NotNil(t, c.EnvFrom[0].SecretRef, "first envFrom entry must be the base secretRef")
			assert.Equal(t, "kubemq-cluster", c.EnvFrom[0].SecretRef.Name)

			// Second entry is always the base configMapRef.
			require.NotNil(t, c.EnvFrom[1].ConfigMapRef, "second envFrom entry must be the base configMapRef")
			assert.Equal(t, "kubemq-cluster", c.EnvFrom[1].ConfigMapRef.Name)

			// Remaining entries are the extra secretRefs, in order, after the base pair.
			var gotExtra []string
			for _, ef := range c.EnvFrom[2:] {
				require.NotNil(t, ef.SecretRef, "extra envFrom entries must be secretRefs")
				gotExtra = append(gotExtra, ef.SecretRef.Name)
			}
			assert.Equal(t, tt.wantExtra, gotExtra)
		})
	}
}
