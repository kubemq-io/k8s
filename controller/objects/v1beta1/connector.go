package v1beta1

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/kubemq-io/k8s/api/v1beta1"
	"github.com/kubemq-io/k8s/controller/config"
	"github.com/kubemq-io/k8s/pkg/subset"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Connector struct {
	*config.Configuration
	Log logr.Logger
}

func NewConnector(cfg *config.Configuration) *Connector {
	return &Connector{
		Configuration: cfg,
		Log:           cfg.Log.WithValues("api-version", "apiVersion: core.k8s.kubemq.io/v1beta1", "kind", "Connector"),
	}
}
func (c *Connector) Apply(ctx context.Context, manifest string) error {
	parsed := &v1beta1.KubemqConnector{}
	found := &v1beta1.KubemqConnector{}
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
				return fmt.Errorf("create connector error, %w", err)
			}
			c.Log.Info("object created", "name", parsed.Name, "namespace", c.Namespace)
			return nil
		} else {
			return err
		}
	} else {
		if !subset.SubsetEqual(parsed.Spec, found.Spec) {
			parsed.ResourceVersion = found.ResourceVersion
			err = c.Client.Update(ctx, parsed)
			if err != nil {
				c.Log.Error(err, "update object error", "name", parsed.Name, "namespace", c.Namespace)
				return fmt.Errorf("update connector error, %w", err)
			}
			c.Log.Info("object configured", "name", parsed.Name, "namespace", c.Namespace)
			return nil
		} else {
			c.Log.Info("object unchanged", "name", parsed.Name, "namespace", c.Namespace)
			return nil
		}
	}
}
func (c *Connector) Delete(ctx context.Context, manifest string) error {
	parsed := &v1beta1.KubemqConnector{}
	found := &v1beta1.KubemqConnector{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = c.Namespace
	err := c.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name, Namespace: parsed.Namespace}, found)
	if err != nil && client.IgnoreNotFound(err) == nil {
		c.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", c.Namespace)
		return fmt.Errorf("delete connector error, %w", err)

	} else {
		err := c.Client.Delete(ctx, found)
		if err != nil {
			c.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", c.Namespace)
			return fmt.Errorf("delete connector error, %w", err)
		}
		c.Log.Info("object deleted", "name", parsed.Name, "namespace", c.Namespace)
		return nil
	}
}
