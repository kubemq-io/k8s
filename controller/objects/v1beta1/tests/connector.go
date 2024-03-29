package tests

const (
	DefaultConnectorManifest = `
apiVersion: core.k8s.kubemq.io/v1beta1
kind: KubemqConnector
metadata:
  name: kubemq-bridges-test
spec:
  type: bridges
  replicas: 1
  config: |-
    bindings:
      - name: bridges-example-binding
        properties:
          log_level: "info"
        sources:
          kind: source.events
          name: cluster-sources
          connections:
            - address: "kubemq-cluster-grpc:50000"
              client_id: "cluster-events-source"
              channel: "events.source"
              group:   ""
              concurrency: "1"
              auto_reconnect: "true"
              reconnect_interval_seconds: "1"
              max_reconnects: "0"
        targets:
          kind: target.events
          name: cluster-targets
          connections:
            - address: "kubemq-cluster-grpc:50000"
              client_id: "cluster-events-target"
              channels: "events.target"
`
)
