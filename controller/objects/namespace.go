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

type Namespace struct {
	*config.Configuration
	Log logr.Logger
}

func NewNamespace(cfg *config.Configuration) *Namespace {
	return &Namespace{
		Configuration: cfg,
		Log:           cfg.Log.WithValues("api-version", "v1", "kind", "Namespace"),
	}
}
func (n *Namespace) Apply(ctx context.Context, manifest string) error {
	parsed := &corev1.Namespace{}
	found := &corev1.Namespace{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = n.Namespace
	err := n.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name, Namespace: parsed.Namespace}, found)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			parsed.Namespace = n.Namespace
			err = n.Client.Create(ctx, parsed)
			if err != nil {
				n.Log.Error(err, "create object error", "name", parsed.Name, "namespace", n.Namespace)
				return fmt.Errorf("create namespace error, %w", err)
			}
			n.Log.Info("object created", "name", parsed.Name, "namespace", n.Namespace)
			return nil
		} else {
			return err
		}
	} else {
		if !subset.SubsetEqual(parsed, found) {
			parsed.ResourceVersion = found.ResourceVersion
			err = n.Client.Update(ctx, parsed)
			if err != nil {
				n.Log.Error(err, "update object error", "name", parsed.Name, "namespace", n.Namespace)
				return fmt.Errorf("update namespace error, %w", err)
			}
			n.Log.Info("object configured", "name", parsed.Name, "namespace", n.Namespace)
			return nil
		} else {
			n.Log.Info("object unchanged", "name", parsed.Name, "namespace", n.Namespace)
			return nil
		}
	}
}
func (n *Namespace) Delete(ctx context.Context, manifest string) error {
	parsed := &corev1.Namespace{}
	found := &corev1.Namespace{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = n.Namespace
	err := n.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name}, found)
	if err != nil && client.IgnoreNotFound(err) == nil {
		n.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", n.Namespace)
		return fmt.Errorf("delete service account error, %w", err)

	} else {
		err := n.Client.Delete(ctx, found)
		if err != nil {
			n.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", n.Namespace)
			return fmt.Errorf("delete service account error, %w", err)
		}
		n.Log.Info("object deleted", "name", parsed.Name, "namespace", n.Namespace)
		return nil
	}
}
