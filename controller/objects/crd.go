package objects

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/kubemq-io/k8s/controller/config"
	"github.com/kubemq-io/k8s/pkg/subset"
	ext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Crd struct {
	*config.Configuration
	Log logr.Logger
}

func NewCrd(cfg *config.Configuration) *Crd {
	return &Crd{
		Configuration: cfg,
		Log:           cfg.Log.WithValues("api-version", "apiextensions.k8s.io/v1", "kind", "CustomResourceDefinition"),
	}
}
func (c *Crd) Apply(ctx context.Context, manifest string) error {
	parsed := &ext.CustomResourceDefinition{}
	found := &ext.CustomResourceDefinition{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	err := c.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name}, found)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			parsed.Namespace = c.Namespace
			err = c.Client.Create(ctx, parsed)
			if err != nil {
				c.Log.Error(err, "create object error", "name", parsed.Name, "namespace", c.Namespace)
				return fmt.Errorf("create crd error, %w", err)
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
				return fmt.Errorf("update crd error, %w", err)
			}
			c.Log.Info("object configured", "name", parsed.Name, "namespace", c.Namespace)
			return nil
		} else {
			c.Log.Info("object unchanged", "name", found.Name, "namespace", c.Namespace)
			return nil
		}
	}
}
func (c *Crd) Delete(ctx context.Context, manifest string) error {
	parsed := &ext.CustomResourceDefinition{}
	found := &ext.CustomResourceDefinition{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = c.Namespace
	err := c.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name}, found)
	if err != nil && client.IgnoreNotFound(err) == nil {
		c.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", c.Namespace)
		return fmt.Errorf("delete crd error, %w", err)

	} else {
		err := c.Client.Delete(ctx, found)
		if err != nil {
			c.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", c.Namespace)
			return fmt.Errorf("delete crd error, %w", err)
		}
		c.Log.Info("object deleted", "name", parsed.Name, "namespace", c.Namespace)
		return nil
	}
}
