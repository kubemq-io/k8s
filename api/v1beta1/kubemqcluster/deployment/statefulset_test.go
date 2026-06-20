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
