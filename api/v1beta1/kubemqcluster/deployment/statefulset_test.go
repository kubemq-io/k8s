package deployment

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var statefulSetConfigData = `
spec:
  replicas: 5
  selector:
    matchLabels:
      app: kubemq-cluster
  serviceName: kubemq-cluster
  template:
    metadata:
      annotations:
        prometheus.io/path: /metrics
        prometheus.io/port: "8080"
        prometheus.io/scrape: "true"
      creationTimestamp: null
      labels:
        app: kubemq-cluster
    spec:
      containers:
      - env:
        - name: CLUSTER_NAME
          value: kubemq-cluster
        - name: CLUSTER_ROUTES
          value: kubemq-cluster:5228
        - name: CLUSTER_ENABLE
          value: "true"
        - name: CHECKSUM
        envFrom:
        - secretRef:
            name: kubemq-cluster
        - configMapRef:
            name: kubemq-cluster
        imagePullPolicy: Always
        name: kubemq-cluster
        ports:
        - containerPort: 50000
          name: grpc-port
          protocol: TCP
        - containerPort: 8080
          name: api-port
          protocol: TCP
        - containerPort: 9090
          name: rest-port
          protocol: TCP
        - containerPort: 5228
          name: cluster-port
          protocol: TCP
        resources: {}
      restartPolicy: Always
      securityContext:
        fsGroup: 200
  updateStrategy:
    type: RollingUpdate
`

func TestStatefulSetConfig_Spec(t *testing.T) {

	tests := []struct {
		name    string
		cfg     *StatefulSetConfig
		wantErr bool
	}{
		{
			name: "full",
			cfg: &StatefulSetConfig{
				Id:                    "",
				Name:                  "kubemq-cluster",
				Namespace:             "kubemq",
				ImagePullPolicy:       "Always",
				Replicas:              5,
				Volume:                "",
				StorageClass:          "",
				statefulset:           nil,
				Health:                "",
				Resources:             "",
				NodeSelectors:         "",
				Image:                 "",
				ServiceAccount:        "",
				ConfigCheckSum:        "",
				Standalone:            false,
				StatefulSetConfigData: "",
			},
			wantErr: false,
		},
		{
			name: "with_template",
			cfg: &StatefulSetConfig{
				Id:                    "",
				Name:                  "kubemq-cluster",
				Namespace:             "kubemq",
				ImagePullPolicy:       "Always",
				Replicas:              5,
				Volume:                "",
				StorageClass:          "",
				statefulset:           nil,
				Health:                "",
				Resources:             "",
				NodeSelectors:         "",
				Image:                 "",
				ServiceAccount:        "",
				ConfigCheckSum:        "",
				Standalone:            false,
				StatefulSetConfigData: statefulSetConfigData,
			},
			wantErr: false,
		},
		{
			name: "stand_alone",
			cfg: &StatefulSetConfig{
				Id:                    "",
				Name:                  "kubemq-cluster",
				Namespace:             "kubemq",
				ImagePullPolicy:       "Always",
				Replicas:              5,
				Volume:                "",
				StorageClass:          "",
				statefulset:           nil,
				Health:                "",
				Resources:             "",
				NodeSelectors:         "",
				Image:                 "",
				ServiceAccount:        "",
				ConfigCheckSum:        "",
				Standalone:            true,
				StatefulSetConfigData: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sts, err := tt.cfg.Get()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.EqualValues(t, tt.cfg.Name, sts.Name)
				data, _ := yaml.Marshal(sts)
				fmt.Println(string(data))
			}
		})
	}
}

// connectorContainerPorts is the static, always-present set of connector
// containerPorts that the default broker StatefulSet template must render
// regardless of any connector config or the Standalone flag.
var connectorContainerPorts = map[string]int32{
	"mqtt":      1883,
	"mqtt-tls":  8883,
	"mqtt-ws":   8083,
	"amqp":      5672,
	"amqp-tls":  5671,
	"stomp":     61613,
	"stomp-tls": 61614,
	"aws-http":  4566,
	"gcp-grpc":  8085,
	"kafka":     9092,
	"kafka-tls": 9093,
}

func TestStatefulSetConfig_ConnectorContainerPorts(t *testing.T) {
	tests := []struct {
		name       string
		standalone bool
	}{
		{name: "cluster", standalone: false},
		{name: "standalone", standalone: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultStatefulSetConfig("", "kubemq-cluster", "kubemq").
				SetStandalone(tt.standalone)

			sts, err := cfg.Get()
			require.NoError(t, err)
			require.Len(t, sts.Spec.Template.Spec.Containers, 1)

			got := map[string]int32{}
			for _, p := range sts.Spec.Template.Spec.Containers[0].Ports {
				got[p.Name] = p.ContainerPort
			}

			for name, port := range connectorContainerPorts {
				gotPort, ok := got[name]
				require.Truef(t, ok, "connector containerPort %q (%d) missing from rendered StatefulSet (standalone=%v)", name, port, tt.standalone)
				assert.Equalf(t, port, gotPort, "connector containerPort %q has wrong port", name)
			}
		})
	}
}

// TestStatefulSetConfig_NextClustered asserts the clustered-next surface renders
// only when engine=next AND not standalone: the raft containerPort, POD_NAME
// downward-API env, CLUSTER_REPLICATION_PEERS, podManagementPolicy=Parallel, and the
// /ready readinessProbe. Legacy and standalone-next render none of it.
func TestStatefulSetConfig_NextClustered(t *testing.T) {
	const peers = "1@kubemq-cluster-0.kubemq-cluster.kubemq.svc.cluster.local:5229," +
		"2@kubemq-cluster-1.kubemq-cluster.kubemq.svc.cluster.local:5229," +
		"3@kubemq-cluster-2.kubemq-cluster.kubemq.svc.cluster.local:5229"

	tests := []struct {
		name        string
		engine      string
		standalone  bool
		wantSurface bool
	}{
		{name: "clustered_next", engine: "next", standalone: false, wantSurface: true},
		{name: "standalone_next", engine: "next", standalone: true, wantSurface: false},
		{name: "legacy", engine: "", standalone: false, wantSurface: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultStatefulSetConfig("", "kubemq-cluster", "kubemq").
				SetStandalone(tt.standalone).SetEngine(tt.engine).SetReplicationPeers(peers)
			sts, err := cfg.Get()
			require.NoError(t, err)
			require.Len(t, sts.Spec.Template.Spec.Containers, 1)
			c := sts.Spec.Template.Spec.Containers[0]

			ports := map[string]int32{}
			for _, p := range c.Ports {
				ports[p.Name] = p.ContainerPort
			}
			env := map[string]string{}
			var podName *string
			for _, e := range c.Env {
				env[e.Name] = e.Value
				if e.Name == "POD_NAME" && e.ValueFrom != nil && e.ValueFrom.FieldRef != nil {
					fp := e.ValueFrom.FieldRef.FieldPath
					podName = &fp
				}
			}

			if tt.wantSurface {
				assert.Equal(t, int32(5229), ports["raft-port"], "raft-port must render")
				require.NotNil(t, podName, "POD_NAME downward-API env must render")
				assert.Equal(t, "metadata.name", *podName)
				assert.Equal(t, peers, env["CLUSTER_REPLICATION_PEERS"])
				assert.Equal(t, "Parallel", string(sts.Spec.PodManagementPolicy))
				require.NotNil(t, c.ReadinessProbe, "readinessProbe must render")
				require.NotNil(t, c.ReadinessProbe.HTTPGet)
				assert.Equal(t, "/ready", c.ReadinessProbe.HTTPGet.Path)
				assert.Equal(t, 8080, c.ReadinessProbe.HTTPGet.Port.IntValue())
			} else {
				_, hasRaft := ports["raft-port"]
				assert.False(t, hasRaft, "raft-port must NOT render")
				assert.Nil(t, podName, "POD_NAME must NOT render")
				_, hasPeers := env["CLUSTER_REPLICATION_PEERS"]
				assert.False(t, hasPeers, "CLUSTER_REPLICATION_PEERS must NOT render")
				assert.NotEqual(t, "Parallel", string(sts.Spec.PodManagementPolicy))
				assert.Nil(t, c.ReadinessProbe, "readinessProbe must NOT render")
			}
			// The operator never injects the server-derived replica id.
			_, hasReplicaID := env["CLUSTER_REPLICATION_REPLICA_ID"]
			assert.False(t, hasReplicaID, "CLUSTER_REPLICATION_REPLICA_ID must never render")
		})
	}
}
