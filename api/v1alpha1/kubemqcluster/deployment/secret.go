package deployment

import (
	"encoding/base64"
	"github.com/ghodss/yaml"
	"github.com/kubemq-io/k8s/pkg/template"
	apiv1 "k8s.io/api/core/v1"
	"reflect"
	"strings"
)

var defaultKubeMQSecretTemplate = `
apiVersion: v1
kind: Secret
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.Name}}
type: Opaque
data:
{{ range $key, $value := .DataVariables}}
  {{$key}}: "{{$value}}"
{{end}}
stringData:
{{ range $key, $value := .StringVariables}}
  {{$key}}: "{{$value}}"
{{end}}
`

type SecretConfig struct {
	Id              string
	Name            string
	Namespace       string
	DataVariables   map[string]string
	StringVariables map[string]string
	secret          *apiv1.Secret
}

func ImportSecret(spec []byte) (*SecretConfig, error) {
	sec := &apiv1.Secret{}
	err := yaml.Unmarshal(spec, sec)
	if err != nil {
		return nil, err
	}
	return &SecretConfig{
		Id:              "",
		Name:            sec.Name,
		Namespace:       sec.Namespace,
		DataVariables:   nil,
		StringVariables: nil,
		secret:          sec,
	}, nil
}
func NewSecretConfig(id, name, namespace string) *SecretConfig {
	return &SecretConfig{
		Id:              id,
		Name:            name,
		Namespace:       namespace,
		DataVariables:   map[string]string{},
		StringVariables: map[string]string{},
	}
}
func DefaultSecretConfig(id, name, namespace string) map[string]*SecretConfig {
	secs := make(map[string]*SecretConfig)
	secs[name] = &SecretConfig{
		Id:              id,
		Name:            name,
		Namespace:       namespace,
		DataVariables:   map[string]string{},
		StringVariables: map[string]string{},
	}
	return secs
}
func (s *SecretConfig) SetDataVariable(key, value string) *SecretConfig {
	if value != "" {
		s.DataVariables[strings.ToUpper(key)] = base64.StdEncoding.EncodeToString([]byte(value))
	}

	return s
}
func (s *SecretConfig) SetStringVariable(key, value string) *SecretConfig {
	if value != "" {
		s.StringVariables[strings.ToUpper(key)] = value
	}
	return s
}
func (s *SecretConfig) Spec() ([]byte, error) {
	if s.secret == nil {

		t := template.NewTemplate(defaultKubeMQSecretTemplate, s)
		return t.Get()
	}
	return yaml.Marshal(s.secret)
}

func (s *SecretConfig) Set(value *apiv1.Secret) *SecretConfig {
	s.secret = value
	return s
}
func (s *SecretConfig) Get() (*apiv1.Secret, error) {
	if s.secret != nil {
		return s.secret, nil
	}
	sec := &apiv1.Secret{}
	data, err := s.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, sec)
	if err != nil {
		return nil, err
	}
	s.secret = sec
	return sec, nil
}

func (s *SecretConfig) StringData() map[string]string {
	if s.secret == nil {
		return nil
	}
	out := make(map[string]string)
	for key, value := range s.secret.StringData {
		out[key] = value
	}

	return out
}
func (s *SecretConfig) Data() map[string][]byte {
	if s.secret == nil {
		return nil
	}
	out := make(map[string][]byte)
	for key, value := range s.secret.Data {
		out[key] = value
	}
	return out
}

func (s *SecretConfig) Equal(target *apiv1.Secret) bool {
	if s.secret == nil {
		return false
	}
	return reflect.DeepEqual(s.secret.Data, target.Data) && reflect.DeepEqual(s.secret.StringData, target.StringData)
}
