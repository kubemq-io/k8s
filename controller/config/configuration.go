package config

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Configuration struct {
	Client    client.Client
	Reader    client.Reader
	Namespace string
	Log       logr.Logger
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}

func (c *Configuration) SetClient(value client.Client) *Configuration {
	c.Client = value
	return c
}

func (c *Configuration) SetNamespace(value string) *Configuration {
	c.Namespace = value
	return c
}
func (c *Configuration) SetReader(value client.Reader) *Configuration {
	c.Reader = value
	return c
}
func (c *Configuration) SetLog(value logr.Logger) *Configuration {
	c.Log = value
	return c
}
