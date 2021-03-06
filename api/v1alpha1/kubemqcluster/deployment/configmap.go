package deployment

import (
	"encoding/base64"
	"github.com/ghodss/yaml"
	"github.com/kubemq-io/k8s/pkg/template"
	apiv1 "k8s.io/api/core/v1"
	"reflect"
	"strings"
)

var defaultKubeMQConfigMapTemplate = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.Name}}
data:
{{ range $key, $value := .Variables}}
  {{$key}}: "{{$value}}"
{{end}}
`

type ConfigMapConfig struct {
	Id        string
	Name      string
	Namespace string
	Variables map[string]string
	configMap *apiv1.ConfigMap
}

func DefaultConfigMap(id, name, namespace string) map[string]*ConfigMapConfig {
	cm := make(map[string]*ConfigMapConfig)
	cm[name] = &ConfigMapConfig{
		Id:        id,
		Name:      name,
		Namespace: namespace,
		Variables: map[string]string{},
	}
	return cm
}

func (c *ConfigMapConfig) SetStringVariable(key, value string) *ConfigMapConfig {
	if value != "" {
		c.Variables[strings.ToUpper(key)] = value
	}
	return c
}
func (c *ConfigMapConfig) SetDataVariable(key, value string) *ConfigMapConfig {
	if value != "" {
		c.Variables[strings.ToUpper(key)] = base64.StdEncoding.EncodeToString([]byte(value))
	}
	return c
}
func (c *ConfigMapConfig) Spec() ([]byte, error) {
	t := template.NewTemplate(defaultKubeMQConfigMapTemplate, c)
	return t.Get()
}
func (c *ConfigMapConfig) Set(value *apiv1.ConfigMap) *ConfigMapConfig {
	c.configMap = value
	return c
}
func (c *ConfigMapConfig) Get() (*apiv1.ConfigMap, error) {
	if c.configMap != nil {
		return c.configMap, nil
	}
	cm := &apiv1.ConfigMap{}
	data, err := c.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, cm)
	if err != nil {
		return nil, err
	}
	c.configMap = cm
	return cm, nil
}

func (c *ConfigMapConfig) Data() map[string]string {
	if c.configMap == nil {
		return nil
	}
	out := make(map[string]string)
	for key, value := range c.configMap.Data {
		if !strings.Contains(value, " ") {
			sDec, err := base64.StdEncoding.DecodeString(value)
			if err == nil {
				out[key] = string(sDec)
				continue
			}

		}
		out[key] = value
	}

	return out
}
func (c *ConfigMapConfig) EqualConfigMap(target *apiv1.ConfigMap) bool {
	if c.configMap == nil {
		return false
	}
	return reflect.DeepEqual(c.configMap.Data, target.Data) && reflect.DeepEqual(c.configMap.BinaryData, target.BinaryData)
}
