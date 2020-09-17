package objects

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	deploymentTest = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: deployment-test
  template:
    metadata:
      labels:
        app: deployment-test
    spec:
      containers:
        - name: kubemq-operator
          image: kubemq/kubemq-operator:latest
          command:
            - kubemq-operator
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
            - containerPort: 8090
          env:
            - name: WATCH_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: DEPLOY_ON_START
              value: false
            - name: SOURCE
              value: "gcp"
            - name: DEBUG_MODE
              value: "false"
`
)

func TestDeployment_ApplyDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	d := &Deployment{
		Configuration: testConfig,
	}
	err := d.Apply(ctx, deploymentTest)
	require.NoError(t, err)
	time.Sleep(5 * time.Second)
	err = d.Delete(ctx, deploymentTest)
	require.NoError(t, err)
}
