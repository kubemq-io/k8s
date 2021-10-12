package objects

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCrd_Get(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	d := NewCrd(testConfig)
	found,err := d.GetBeta1(ctx,"kubemqclusters.core.k8s.kubemq.io")
	require.NoError(t, err)
	fmt.Println(found)
}
