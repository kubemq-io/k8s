package manifests

const KubemqRole = `
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: kubemq-cluster
rules:
  - apiGroups:
      - security.openshift.io
    resources:
      - securitycontextconstraints
    verbs:
      - use
      - delete
      - get
      - list
      - patch
      - update
      - watch
    resourceNames:
      - privileged
`
