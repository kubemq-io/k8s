package v1alpha1

import (
	"context"
	tests2 "github.com/kubemq-io/k8s/controller/objects/v1alpha1/tests"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	"testing"
	"time"
)

func TestCluster_Apply(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	c := NewCluster(testConfig)
	err := c.Apply(ctx, tests2.DefaultClusterManifest)
	if meta.IsNoMatchError(err) {
		t.Skipf("skipping: KubemqCluster v1alpha1 CRD not registered in cluster (%v)", err)
	}
	require.NoError(t, err)
	time.Sleep(2 * time.Second)
	err = c.Apply(ctx, tests2.UpdatedClusterManifest)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)
	err = c.Delete(ctx, tests2.DefaultClusterManifest)
	require.NoError(t, err)
}
