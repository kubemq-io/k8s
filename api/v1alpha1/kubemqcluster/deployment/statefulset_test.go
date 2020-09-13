package deployment

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStatefulSetConfig_Spec(t *testing.T) {

	tests := []struct {
		name    string
		cfg     *StatefulSetConfig
		wantErr bool
	}{
		{
			name: "full",
			cfg: &StatefulSetConfig{
				Id:              "",
				Name:            "kubemq-cluster",
				Namespace:       "kubemq",
				ImagePullPolicy: "Always",
				Replicas:        5,
				Volume:          "",
				StorageClass:    "",
				statefulset:     nil,
				Health:          "",
				Resources:       "",
				NodeSelectors:   "",
				Image:           "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sts, err := tt.cfg.Get()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.EqualValues(t, tt.cfg.Name, sts.Name)
			}
		})
	}
}
