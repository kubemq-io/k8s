package objects

import (
	"context"
	tests2 "github.com/kubemq-io/k8s/controller/objects/tests"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCluster_Apply(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	c := &Cluster{
		Configuration: testConfig,
	}
	err := c.Apply(ctx, tests2.DefaultClusterManifest)
	require.NoError(t, err)
	time.Sleep(10 * time.Second)
	err = c.Apply(ctx, tests2.UpdatedClusterManifest)
	require.NoError(t, err)
	time.Sleep(10 * time.Second)
	err = c.Delete(ctx, tests2.DefaultClusterManifest)
	require.NoError(t, err)
}
