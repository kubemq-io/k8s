package controller

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-io/k8s/controller/config"
	"github.com/kubemq-io/k8s/controller/objects"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deployer interface {
	Apply(ctx context.Context, manifest string) error
	Delete(ctx context.Context, manifest string) error
}

type DeployerService struct {
	cfg *config.Configuration
}

func NewDeployer(cfg *config.Configuration) *DeployerService {
	return &DeployerService{
		cfg: cfg,
	}
}

func (d *DeployerService) Apply(ctx context.Context, manifest string) error {
	typeMeta := metav1.TypeMeta{}
	err := yaml.Unmarshal([]byte(manifest), &typeMeta)
	if err != nil {
		return fmt.Errorf("invalid object, cann't parse object type metadata, %w", err)
	}
	key := fmt.Sprintf("%s/%s", typeMeta.APIVersion, typeMeta.Kind)
	switch key {
	case "core.k8s.kubemq.io/v1alpha1/KubemqCluster":
		return objects.NewCluster(d.cfg).Apply(ctx, manifest)
	case "core.k8s.kubemq.io/v1alpha1/KubemqDashboard":
		return objects.NewDashboard(d.cfg).Apply(ctx, manifest)
	case "core.k8s.kubemq.io/v1alpha1/KubemqConnector":
		return objects.NewConnector(d.cfg).Apply(ctx, manifest)
	case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition":
		return objects.NewCrd(d.cfg).Apply(ctx, manifest)
	case "v1/ServiceAccount":
		return objects.NewServiceAccount(d.cfg).Apply(ctx, manifest)
	case "v1/ConfigMap":
		return objects.NewConfigMap(d.cfg).Apply(ctx, manifest)
	case "v1/Service":
		return objects.NewConfigMap(d.cfg).Apply(ctx, manifest)
	default:
		return fmt.Errorf("unknown object to apply")
	}
}

func (d *DeployerService) Delete(ctx context.Context, manifest string) error {
	typeMeta := metav1.TypeMeta{}
	err := yaml.Unmarshal([]byte(manifest), &typeMeta)
	if err != nil {
		return fmt.Errorf("invalid object, cann't parse object type metadata, %w", err)
	}
	key := fmt.Sprintf("%s/%s", typeMeta.APIVersion, typeMeta.Kind)
	switch key {
	case "core.k8s.kubemq.io/v1alpha1/KubemqCluster":
		return objects.NewCluster(d.cfg).Delete(ctx, manifest)
	case "core.k8s.kubemq.io/v1alpha1/KubemqDashboard":
		return objects.NewDashboard(d.cfg).Delete(ctx, manifest)
	case "core.k8s.kubemq.io/v1alpha1/KubemqConnector":
		return objects.NewConnector(d.cfg).Delete(ctx, manifest)
	case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition":
		return objects.NewCrd(d.cfg).Delete(ctx, manifest)
	case "v1/ServiceAccount":
		return objects.NewServiceAccount(d.cfg).Delete(ctx, manifest)
	case "v1/ConfigMap":
		return objects.NewConfigMap(d.cfg).Delete(ctx, manifest)
	case "v1/Service":
		return objects.NewConfigMap(d.cfg).Delete(ctx, manifest)
	default:
		return fmt.Errorf("unknown object to apply")
	}
}
