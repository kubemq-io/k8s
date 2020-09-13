package objects

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-io/k8s/controller/config"
	"github.com/kubemq-io/k8s/pkg/subset"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceAccount struct {
	*config.Configuration
}

func NewServiceAccount(cfg *config.Configuration) *ServiceAccount {
	return &ServiceAccount{
		Configuration: cfg,
	}
}
func (c *ServiceAccount) Apply(ctx context.Context, manifest string) error {
	parsed := &corev1.ServiceAccount{}
	found := &corev1.ServiceAccount{}
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
				c.Log.Error(err, "create object error", "name", parsed.Name, "namespace", c.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
				return fmt.Errorf("create service account error, %w", err)
			}
			c.Log.Info("object created", "name", parsed.Name, "namespace", c.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
			return nil
		} else {
			return err
		}
	} else {
		if !subset.SubsetEqual(parsed, found) {
			parsed.ResourceVersion = found.ResourceVersion
			err = c.Client.Update(ctx, parsed)
			if err != nil {
				c.Log.Error(err, "update object error", "name", parsed.Name, "namespace", c.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
				return fmt.Errorf("update service account error, %w", err)
			}
			c.Log.Info("object configured", "name", parsed.Name, "namespace", c.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
			return nil
		} else {
			c.Log.Info("object unchanged", "name", parsed.Name, "namespace", c.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
			return nil
		}
	}
}
func (c *ServiceAccount) Delete(ctx context.Context, manifest string) error {
	parsed := &corev1.ServiceAccount{}
	found := &corev1.ServiceAccount{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = c.Namespace
	err := c.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name}, found)
	if err != nil && client.IgnoreNotFound(err) == nil {
		c.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", c.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
		return fmt.Errorf("delete service account error, %w", err)

	} else {
		err := c.Client.Delete(ctx, found)
		if err != nil {
			c.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", c.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
			return fmt.Errorf("delte service account error, %w", err)
		}
		c.Log.Info("object deleted", "name", parsed.Name, "namespace", c.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
		return nil
	}
}
