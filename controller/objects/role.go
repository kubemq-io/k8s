package objects

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/kubemq-io/k8s/controller/config"
	"github.com/kubemq-io/k8s/pkg/subset"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Role struct {
	*config.Configuration
	Log logr.Logger
}

func NewRole(cfg *config.Configuration) *Role {
	return &Role{
		Configuration: cfg,
		Log:           cfg.Log.WithValues("api-version", "rbac.authorization.k8s.io/v1", "kind", "Role"),
	}
}
func (r *Role) Apply(ctx context.Context, manifest string) error {
	parsed := &rbac.Role{}
	found := &rbac.Role{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = r.Namespace
	err := r.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name, Namespace: parsed.Namespace}, found)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			parsed.Namespace = r.Namespace
			err = r.Client.Create(ctx, parsed)
			if err != nil {
				r.Log.Error(err, "create object error", "name", parsed.Name, "namespace", r.Namespace)
				return fmt.Errorf("create role error, %w", err)
			}
			r.Log.Info("object created", "name", parsed.Name, "namespace", r.Namespace)
			return nil
		} else {
			return err
		}
	} else {
		if !subset.SubsetEqual(parsed, found) {
			parsed.ResourceVersion = found.ResourceVersion
			err = r.Client.Update(ctx, parsed)
			if err != nil {
				r.Log.Error(err, "update object error", "name", parsed.Name, "namespace", r.Namespace)
				return fmt.Errorf("update role error, %w", err)
			}
			r.Log.Info("object configured", "name", parsed.Name, "namespace", r.Namespace)
			return nil
		} else {
			r.Log.Info("object unchanged", "name", parsed.Name, "namespace", r.Namespace)
			return nil
		}
	}
}
func (r *Role) Delete(ctx context.Context, manifest string) error {
	parsed := &rbac.Role{}
	found := &rbac.Role{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = r.Namespace
	err := r.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name, Namespace: parsed.Namespace}, found)
	if err != nil && client.IgnoreNotFound(err) == nil {
		r.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", r.Namespace)
		return fmt.Errorf("delete role error, %w", err)

	} else {
		err := r.Client.Delete(ctx, found)
		if err != nil {
			r.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", r.Namespace)
			return fmt.Errorf("delete role error, %w", err)
		}
		r.Log.Info("object deleted", "name", parsed.Name, "namespace", r.Namespace)
		return nil
	}
}
