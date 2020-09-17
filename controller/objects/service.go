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

type Service struct {
	*config.Configuration
	Log logr.Logger
}

func NewService(cfg *config.Configuration) *Service {
	return &Service{
		Configuration: cfg,
		Log:           cfg.Log.WithValues("api-version", "v1", "kind", "Service"),
	}
}
func (s *Service) Apply(ctx context.Context, manifest string) error {
	parsed := &corev1.Service{}
	found := &corev1.Service{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = s.Namespace
	err := s.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name, Namespace: parsed.Namespace}, found)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			parsed.Namespace = s.Namespace
			err = s.Client.Create(ctx, parsed)
			if err != nil {
				s.Log.Error(err, "create object error", "name", parsed.Name, "namespace", s.Namespace)
				return fmt.Errorf("create service error, %w", err)
			}
			s.Log.Info("object created", "name", parsed.Name, "namespace", s.Namespace)
			return nil
		} else {
			return err
		}
	} else {
		if !subset.SubsetEqual(parsed.Spec, found.Spec) {
			parsed.ResourceVersion = found.ResourceVersion
			parsed.Spec.ClusterIP = found.Spec.ClusterIP
			err = s.Client.Update(ctx, parsed)
			if err != nil {
				s.Log.Error(err, "update object error", "name", parsed.Name, "namespace", s.Namespace)
				return fmt.Errorf("update service error, %w", err)
			}
			s.Log.Info("object configured", "name", parsed.Name, "namespace", s.Namespace)
			return nil
		} else {
			s.Log.Info("object unchanged", "name", parsed.Name, "namespace", s.Namespace)
			return nil
		}
	}
}
func (s *Service) Delete(ctx context.Context, manifest string) error {
	parsed := &corev1.Service{}
	found := &corev1.Service{}
	if err := yaml.Unmarshal([]byte(manifest), parsed); err != nil {
		return fmt.Errorf("parsing manifest error, %w", err)
	}
	parsed.Namespace = s.Namespace
	err := s.Reader.Get(ctx, types.NamespacedName{Name: parsed.Name, Namespace: parsed.Namespace}, found)
	if err != nil && client.IgnoreNotFound(err) == nil {
		s.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", s.Namespace)
		return fmt.Errorf("delete service error, %w", err)

	} else {
		err := s.Client.Delete(ctx, found)
		if err != nil {
			s.Log.Error(err, "delete object error", "name", parsed.Name, "namespace", s.Namespace)
			return fmt.Errorf("delete service error, %w", err)
		}
		s.Log.Info("object deleted", "name", parsed.Name, "namespace", s.Namespace)
		return nil
	}
}
