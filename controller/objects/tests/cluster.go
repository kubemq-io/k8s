package tests

const (
	DefaultClusterManifest = `
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster-test
spec:
  replicas: 3
`
	UpdatedClusterManifest = `
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster
spec:
  replicas: 1
`
)
