package objects

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-io/k8s/api/v1alpha1"
	"github.com/kubemq-io/k8s/controller/config"
	"github.com/kubemq-io/k8s/pkg/subset"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Dashboard struct {
	*config.Configuration
}

func NewDashboard(cfg *config.Configuration) *Dashboard {
	return &Dashboard{
		Configuration: cfg,
	}
}
func (d *Dashboard) Apply(ctx context.Context, manifest string) error {
	parsed := &v1alpha1.KubemqDashboard{}
	found := &v1alpha1.KubemqDashboard{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = d.Namespace
	err := d.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name, Namespace: parsed.Namespace}, found)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			parsed.Namespace = d.Namespace
			err = d.Client.Create(ctx, parsed)
			if err != nil {
				d.Log.Error(err, "create object error", "name", parsed.Name, "namespace", d.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
				return fmt.Errorf("create dashboard error, %w", err)
			}
			d.Log.Info("object created", "name", parsed.Name, "namespace", d.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
			return nil
		} else {
			return err
		}
	} else {
		if !subset.SubsetEqual(parsed.Spec, found.Spec) {
			parsed.ResourceVersion = found.ResourceVersion
			err = d.Client.Update(ctx, parsed)
			if err != nil {
				d.Log.Error(err, "update object error", "name", parsed.Name, "namespace", d.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
				return fmt.Errorf("update dashboard error, %w", err)
			}
			d.Log.Info("object configured", "name", parsed.Name, "namespace", d.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
			return nil
		} else {
			d.Log.Info("object unchanged", "name", parsed.Name, "namespace", d.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
			return nil
		}
	}
}
func (d *Dashboard) Delete(ctx context.Context, manifest string) error {
	parsed := &v1alpha1.KubemqDashboard{}
	found := &v1alpha1.KubemqDashboard{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = d.Namespace
	err := d.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name, Namespace: parsed.Namespace}, found)
	if err != nil && client.IgnoreNotFound(err) == nil {
		d.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", d.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
		return fmt.Errorf("delete dashboard error, %w", err)

	} else {
		err := d.Client.Delete(ctx, found)
		if err != nil {
			d.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", d.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
			return fmt.Errorf("delte dashboard error, %w", err)
		}
		d.Log.Info("object deleted", "name", parsed.Name, "namespace", d.Namespace, "api-version", parsed.APIVersion, "kind", parsed.Kind)
		return nil
	}
}
