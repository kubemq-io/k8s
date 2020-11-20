package manifests

const (
	DefaultConnectorBridgesManifest = `
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqConnector
metadata:
  name: kubemq-bridges
spec:
  type: bridges
  replicas: 1
  config: |-
    bindings: null
`
	DefaultConnectorTargetsManifest = `
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqConnector
metadata:
  name: kubemq-targets
spec:
  type: targets
  replicas: 1
  config: |-
    bindings: null
`
	DefaultConnectorSourcesManifest = `
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqConnector
metadata:
  name: kubemq-sources
spec:
  type: sources
  replicas: 1
  config: |-
    bindings: null
`
)
