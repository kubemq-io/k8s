package v1beta1

import (
	"context"
	tests2 "github.com/kubemq-io/k8s/controller/objects/v1beta1/tests"
	"github.com/stretchr/testify/require"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"testing"
	"time"
)

func TestDashboard_Apply(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	// Guard: skip when the KubemqDashboard CRD is not installed in the target cluster.
	// The dashboard CRD may not be present in all environments; a missing CRD returns
	// "no matches for kind" which is an environment pre-condition, not a production code defect.
	d := NewDashboard(testConfig)
	err := d.Apply(ctx, tests2.DefaultDashboardManifest)
	if apimeta.IsNoMatchError(err) {
		t.Skipf("skipping TestDashboard_Apply: KubemqDashboard CRD not installed in cluster (%v)", err)
	}
	require.NoError(t, err)
	time.Sleep(2 * time.Second)
	err = d.Delete(ctx, tests2.DefaultDashboardManifest)
	require.NoError(t, err)
}
