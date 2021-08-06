package v1beta1

import (
	"context"
	tests2 "github.com/kubemq-io/k8s/controller/objects/v1beta1/tests"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConnector_Apply(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	c := NewConnector(testConfig)
	err := c.Apply(ctx, tests2.DefaultConnectorManifest)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)
	err = c.Delete(ctx, tests2.DefaultConnectorManifest)
	require.NoError(t, err)
}
