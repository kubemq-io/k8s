package deployment

import (
	"fmt"
	"strings"
)

type Config struct {
	Id          string
	Name        string
	Namespace   string
	StatefulSet *StatefulSetConfig
	Services    map[string]*ServiceConfig
	ConfigMaps  map[string]*ConfigMapConfig
	Secrets     map[string]*SecretConfig
}

func NewKubeMQManifestConfig(id, name, namespace string) *Config {
	return &Config{
		Id:          id,
		Name:        name,
		Namespace:   namespace,
		StatefulSet: nil,
		Services:    make(map[string]*ServiceConfig),
		ConfigMaps:  make(map[string]*ConfigMapConfig),
		Secrets:     make(map[string]*SecretConfig),
	}
}

func DefaultKubeMQManifestConfig(id, name, namespace string, standalone bool) *Config {
	if standalone {
		return &Config{
			Id:          id,
			Name:        name,
			Namespace:   namespace,
			StatefulSet: DefaultStatefulSetConfig(id, name, namespace),
			Services:    DefaultServiceConfigWithHeadless(id, namespace, name),
			ConfigMaps:  DefaultConfigMap(id, name, namespace),
			Secrets:     DefaultSecretConfig(id, name, namespace),
		}
	} else {
		return &Config{
			Id:          id,
			Name:        name,
			Namespace:   namespace,
			StatefulSet: DefaultStatefulSetConfig(id, name, namespace),
			Services:    DefaultServiceConfig(id, namespace, name),
			ConfigMaps:  DefaultConfigMap(id, name, namespace),
			Secrets:     DefaultSecretConfig(id, name, namespace),
		}
	}

}

func (c *Config) SetSecretStringValues(secName, key, value string) {
	sec, ok := c.Secrets[secName]
	if ok {
		sec.SetStringVariable(key, value)
	}
}
func (c *Config) SetSecretDataValues(secName, key, value string) {
	sec, ok := c.Secrets[secName]
	if ok {
		sec.SetDataVariable(key, value)
	}
}
func (c *Config) SetConfigMapStringValues(cmName, key, value string) {
	cm, ok := c.ConfigMaps[cmName]
	if ok {
		cm.SetStringVariable(key, value)
	}
}
func (c *Config) SetConfigMapDataValues(cmName, key, value string) {
	cm, ok := c.ConfigMaps[cmName]
	if ok {
		cm.SetDataVariable(key, value)
	}
}
func (c *Config) Spec() ([]byte, error) {
	var manifest []string

	if c.StatefulSet != nil {
		stsSpec, err := c.StatefulSet.Spec()
		if err != nil {
			return nil, fmt.Errorf("error on statefull spec rendring: %s", err.Error())
		}
		manifest = append(manifest, string(stsSpec))
	}

	for name, svc := range c.Services {
		svcSpec, err := svc.Spec()
		if err != nil {
			return nil, fmt.Errorf("error on service %s spec rendring: %s", name, err.Error())
		}
		manifest = append(manifest, string(svcSpec))
	}
	for _, cm := range c.ConfigMaps {
		configMapSpec, err := cm.Spec()
		if err != nil {
			return nil, fmt.Errorf("error on config map spec rendring: %s", err.Error())
		}
		manifest = append(manifest, string(configMapSpec))
	}

	for _, sec := range c.Secrets {
		secretSpec, err := sec.Spec()
		if err != nil {
			return nil, fmt.Errorf("error on secret spec rendring: %s", err.Error())
		}
		manifest = append(manifest, string(secretSpec))
	}

	return []byte(strings.Join(manifest, "\n---\n")), nil
}
