package manifests

const (
	KubemqOperatorServiceAccount = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubemq-operator
`
	KubemqClusterServiceAccount = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubemq-cluster
`
)
