package deployment

import (
	"github.com/ghodss/yaml"
	"github.com/kubemq-io/k8s/pkg/template"
	apiv1 "k8s.io/api/core/v1"
)

var defaultKubeMQServiceTemplate = `
apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.AppName}}
spec:
  ports:
    - name: {{.PortName}}
      port: {{.ContainerPort}}
      protocol: TCP
      targetPort: {{.TargetPort}}
      nodePort: {{.NodePort}}
  sessionAffinity: None
{{ if not .Headless }}
  type: {{.Expose}}
{{else}}
  clusterIP: None
{{end}}
  selector:
    app: {{.AppName}}
`
var defaultKubeMQServiceHeadlessTemplate = `
apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.AppName}}
spec:
  ports:
    - name: grpc-port
      port: 50000
      protocol: TCP
      targetPort: 50000
    - name: api-port
      port: 8080
      protocol: TCP
      targetPort: 8080
    - name: rest-port
      port: 9090
      protocol: TCP
      targetPort: 9090
    - name: cluster-port
      port: 5228
      protocol: TCP
      targetPort: 5228
  sessionAffinity: None
  clusterIP: None
  selector:
    app: {{.AppName}}
`

type ServiceConfig struct {
	Id            string
	Name          string
	Namespace     string
	AppName       string
	Expose        string
	ContainerPort int32
	TargetPort    int32
	PortName      string
	NodePort      int32
	Headless      bool
	service       *apiv1.Service
}

func ImportServiceConfig(spec []byte) (*ServiceConfig, error) {
	svc := &apiv1.Service{}
	err := yaml.Unmarshal(spec, svc)
	if err != nil {
		return nil, err
	}
	return &ServiceConfig{
		Id:            "",
		Name:          svc.Name,
		Namespace:     svc.Namespace,
		AppName:       "",
		Expose:        "",
		ContainerPort: 0,
		TargetPort:    0,
		PortName:      "",
		NodePort:      0,
		Headless:      false,
		service:       svc,
	}, nil
}

func NewServiceConfig(id, name, namespace, appName string) *ServiceConfig {
	return &ServiceConfig{
		Id:            id,
		Name:          name,
		Namespace:     namespace,
		AppName:       appName,
		Expose:        "",
		ContainerPort: 0,
		TargetPort:    0,
		PortName:      "",
		NodePort:      0,
		service:       nil,
		Headless:      false,
	}
}

func DefaultServiceConfig(id, namespace, appName string) map[string]*ServiceConfig {
	list := map[string]*ServiceConfig{}
	list["grpc"] = &ServiceConfig{
		Id:            id,
		Name:          appName + "-grpc",
		Namespace:     namespace,
		AppName:       appName,
		Expose:        "ClusterIP",
		ContainerPort: 50000,
		TargetPort:    50000,
		PortName:      "grpc-port",
		NodePort:      0,
		service:       nil,
		Headless:      false,
	}
	list["rest"] = &ServiceConfig{
		Id:            id,
		Name:          appName + "-rest",
		Namespace:     namespace,
		AppName:       appName,
		Expose:        "ClusterIP",
		ContainerPort: 9090,
		TargetPort:    9090,
		PortName:      "rest-port",
		NodePort:      0,
		service:       nil,
		Headless:      false,
	}
	list["api"] = &ServiceConfig{
		Id:            id,
		Name:          appName + "-api",
		Namespace:     namespace,
		AppName:       appName,
		Expose:        "ClusterIP",
		ContainerPort: 8080,
		TargetPort:    8080,
		PortName:      "api-port",
		NodePort:      0,
		service:       nil,
		Headless:      false,
	}
	list["internal"] = &ServiceConfig{
		Id:        id,
		Name:      appName,
		Namespace: namespace,
		AppName:   appName,
		Headless:  true,
	}
	return list
}

func (s *ServiceConfig) SetExpose(value string) *ServiceConfig {
	s.Expose = value
	return s
}
func (s *ServiceConfig) SetContainerPort(value int32) *ServiceConfig {
	s.ContainerPort = value
	return s
}
func (s *ServiceConfig) SetTargetPort(value int32) *ServiceConfig {
	s.TargetPort = value
	return s
}
func (s *ServiceConfig) SetNodePort(value int32) *ServiceConfig {
	s.NodePort = value
	return s
}
func (s *ServiceConfig) SetPortName(value string) *ServiceConfig {
	s.PortName = value
	return s
}
func (s *ServiceConfig) SetHeadless(value bool) *ServiceConfig {
	s.Headless = value
	return s
}

func (s *ServiceConfig) Spec() ([]byte, error) {

	if s.service == nil {
		tmpl := defaultKubeMQServiceTemplate
		if s.Headless {
			tmpl = defaultKubeMQServiceHeadlessTemplate
		}
		t := template.NewTemplate(tmpl, s)
		return t.Get()
	}
	return yaml.Marshal(s.service)
}
func (s *ServiceConfig) Set(value *apiv1.Service) *ServiceConfig {
	s.service = value
	return s
}
func (s *ServiceConfig) Get() (*apiv1.Service, error) {
	if s.service != nil {
		return s.service, nil
	}
	svc := &apiv1.Service{}
	data, err := s.Spec()
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, svc)
	if err != nil {
		return nil, err
	}
	s.service = svc
	return svc, nil
}
