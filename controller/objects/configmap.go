package objects

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/kubemq-io/k8s/controller/config"
	"github.com/kubemq-io/k8s/pkg/subset"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigMap struct {
	*config.Configuration
	Log logr.Logger
}

func NewConfigMap(cfg *config.Configuration) *ConfigMap {
	return &ConfigMap{
		Configuration: cfg,
		Log:           cfg.Log.WithValues("api-version", "v1", "kind", "ConfigMap"),
	}
}
func (c *ConfigMap) Apply(ctx context.Context, manifest string) error {
	parsed := &corev1.ConfigMap{}
	found := &corev1.ConfigMap{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = c.Namespace
	err := c.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name, Namespace: parsed.Namespace}, found)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			parsed.Namespace = c.Namespace
			err = c.Client.Create(ctx, parsed)
			if err != nil {
				c.Log.Error(err, "create object error", "name", parsed.Name, "namespace", c.Namespace)
				return fmt.Errorf("create configmap error, %w", err)
			}
			c.Log.Info("object created", "name", parsed.Name, "namespace", c.Namespace)
			return nil
		} else {
			return err
		}
	} else {
		if !subset.SubsetEqual(parsed, found) {
			parsed.ResourceVersion = found.ResourceVersion
			err = c.Client.Update(ctx, parsed)
			if err != nil {
				c.Log.Error(err, "update object error", "name", parsed.Name, "namespace", c.Namespace)
				return fmt.Errorf("update configmap error, %w", err)
			}
			c.Log.Info("object configured", "name", parsed.Name, "namespace", c.Namespace)
			return nil
		} else {
			c.Log.Info("object unchanged", "name", parsed.Name, "namespace", c.Namespace)
			return nil
		}
	}
}
func (c *ConfigMap) Delete(ctx context.Context, manifest string) error {
	parsed := &corev1.ConfigMap{}
	found := &corev1.ConfigMap{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = c.Namespace
	err := c.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name, Namespace: parsed.Namespace}, found)
	if err != nil {
		c.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", c.Namespace)
		return fmt.Errorf("delete configmap error, %w", err)

	} else {
		err := c.Client.Delete(ctx, found)
		if err != nil {
			c.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", c.Namespace)
			return fmt.Errorf("delete configmap error, %w", err)
		}
		c.Log.Info("object deleted", "name", parsed.Name, "namespace", c.Namespace)
		return nil
	}
}
