package v1beta1

import (
	"context"
	tests2 "github.com/kubemq-io/k8s/controller/objects/v1beta1/tests"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDashboard_Apply(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	d := NewDashboard(testConfig)
	err := d.Apply(ctx, tests2.DefaultDashboardManifest)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)
	err = d.Delete(ctx, tests2.DefaultDashboardManifest)
	require.NoError(t, err)
}
