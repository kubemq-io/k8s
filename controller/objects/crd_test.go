package objects

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	extbeta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"testing"
	"time"
)

func init() {
	utilruntime.Must(extbeta1.AddToScheme(scheme))
}

func TestCrd_Get(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	d := NewCrd(testConfig)
	found, err := d.Get(ctx, "kubemqclusters.core.k8s.kubemq.io")
	require.NoError(t, err)
	fmt.Println(found)
}
