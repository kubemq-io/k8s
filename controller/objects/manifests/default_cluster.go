package manifests

const (
	DefaultClusterManifest = `
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster
spec:
  replicas: 3
`
)
