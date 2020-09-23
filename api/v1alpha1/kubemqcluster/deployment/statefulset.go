package deployment

import (
	"github.com/ghodss/yaml"
	"github.com/kubemq-io/k8s/pkg/template"
	appsv1 "k8s.io/api/apps/v1"
)

var defaultKubeMQStatefulSetTemplate = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.Name}}
spec:
  selector:
    matchLabels:
      app: {{.Name}}
  replicas: {{.Replicas}}
  serviceName: {{.Name}}
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: {{.Name}}
      annotations:
        prometheus.io/scrape: 'true'
        prometheus.io/port: '8080'
        prometheus.io/path: '/metrics'
    spec:
{{ .NodeSelectors }}
      serviceAccountName: {{.ServiceAccount}}
      securityContext:
        fsGroup: 200
      containers:
        - env:
            - name: CLUSTER_NAME
              value: {{.Name}}
            - name: CLUSTER_ENABLE
              value: 'true'
            - name: CLUSTER_ROUTES
              value: '{{.Name}}:5228'
          envFrom:
            - secretRef:
                name: {{.Name}}
            - configMapRef:
                name: {{.Name}}
          image: {{.Image}}
          imagePullPolicy: {{.ImagePullPolicy}}
          name: {{.Name}}
{{ .Health }}
{{ .Resources }}
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
{{if .Volume  }}
          volumeMounts:
            - name: {{.Name}}-vol
              mountPath: './kubemq/store'
{{end}}
      restartPolicy: Always
{{if  .Volume  }}  
  volumeClaimTemplates:
    - metadata:
        name: {{.Name}}-vol
      spec:
        accessModes: [ "ReadWriteOnce" ]
        storageClassName: {{.StorageClass}}
        resources:
          requests:
            storage: {{.Volume}}
{{end}}
`

type StatefulSetConfig struct {
	Id              string
	Name            string
	Namespace       string
	ImagePullPolicy string
	Replicas        int
	Volume          string
	StorageClass    string
	statefulset     *appsv1.StatefulSet
	Health          string
	Resources       string
	NodeSelectors   string
	Image           string
	ServiceAccount  string
}

func DefaultStatefulSetConfig(id, name, namespace string) *StatefulSetConfig {
	return &StatefulSetConfig{
		Id:              id,
		Name:            name,
		Namespace:       namespace,
		ImagePullPolicy: "Always",
		Replicas:        3,
		Volume:          "",
		StorageClass:    "",
		statefulset:     nil,
		Health:          "",
		Resources:       "",
		NodeSelectors:   "",
		Image:           "",
		ServiceAccount:  "kubemq-cluster",
	}
}

func (sc *StatefulSetConfig) SetReplicas(value int) *StatefulSetConfig {
	if value == 0 {
		sc.Replicas = 3
	} else {
		sc.Replicas = value
	}

	return sc
}

func (sc *StatefulSetConfig) SetVolume(value string) *StatefulSetConfig {
	sc.Volume = value
	return sc
}
func (sc *StatefulSetConfig) SetStorageClass(value string) *StatefulSetConfig {
	sc.StorageClass = value
	return sc
}

func (sc *StatefulSetConfig) SetImagePullPolicy(value string) *StatefulSetConfig {
	if value == "" {
		sc.ImagePullPolicy = "Always"

	} else {
		sc.ImagePullPolicy = value
	}

	return sc
}
func (sc *StatefulSetConfig) SetImageName(value string) *StatefulSetConfig {
	sc.Image = value
	return sc
}
func (sc *StatefulSetConfig) SetHealthProbe(value string) *StatefulSetConfig {
	sc.Health = value
	return sc
}
func (sc *StatefulSetConfig) SetResources(value string) *StatefulSetConfig {
	sc.Resources = value
	return sc
}
func (sc *StatefulSetConfig) SetNodeSelectors(value string) *StatefulSetConfig {
	sc.NodeSelectors = value
	return sc
}
func (sc *StatefulSetConfig) SetServiceAccount(value string) *StatefulSetConfig {
	sc.ServiceAccount = value
	return sc
}

func (sc *StatefulSetConfig) Spec() ([]byte, error) {
	if sc.statefulset == nil {
		t := template.NewTemplate(defaultKubeMQStatefulSetTemplate, sc)
		data, err := t.Get()
		return data, err
	}

	return yaml.Marshal(sc.statefulset)
}
func (sc *StatefulSetConfig) Set(value *appsv1.StatefulSet) *StatefulSetConfig {
	sc.statefulset = value
	return sc
}

func (sc *StatefulSetConfig) Get() (*appsv1.StatefulSet, error) {
	if sc.statefulset != nil {
		return sc.statefulset, nil
	}
	sts := &appsv1.StatefulSet{}
	data, err := sc.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, sts)
	if err != nil {
		return nil, err
	}

	return sts, nil
}
