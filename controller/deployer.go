package controller

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/kubemq-io/k8s/controller/config"
	"github.com/kubemq-io/k8s/controller/objects"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
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
	var result *multierror.Error
	for _, spec := range d.splitter(manifest) {
		typeMeta := metav1.TypeMeta{}
		err := yaml.Unmarshal([]byte(spec), &typeMeta)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("invalid object, can't parse object type metadata: %w\n", err))
			continue
		}
		key := fmt.Sprintf("%s/%s", typeMeta.APIVersion, typeMeta.Kind)
		switch key {
		case "core.k8s.kubemq.io/v1alpha1/KubemqCluster":
			result = multierror.Append(result, objects.NewCluster(d.cfg).Apply(ctx, manifest))
		case "core.k8s.kubemq.io/v1alpha1/KubemqDashboard":
			result = multierror.Append(result, objects.NewDashboard(d.cfg).Apply(ctx, manifest))
		case "core.k8s.kubemq.io/v1alpha1/KubemqConnector":
			result = multierror.Append(result, objects.NewConnector(d.cfg).Apply(ctx, manifest))
		case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition":
			result = multierror.Append(result, objects.NewCrd(d.cfg).Apply(ctx, manifest))
		case "v1/ServiceAccount":
			result = multierror.Append(result, objects.NewServiceAccount(d.cfg).Apply(ctx, manifest))
		case "v1/ConfigMap":
			result = multierror.Append(result, objects.NewConfigMap(d.cfg).Apply(ctx, manifest))
		case "v1/Service":
			result = multierror.Append(result, objects.NewService(d.cfg).Apply(ctx, manifest))
		case "apps/v1/Deployment":
			result = multierror.Append(result, objects.NewDeployment(d.cfg).Apply(ctx, manifest))
		}
	}
	if result == nil {
		return nil
	}
	return result.ErrorOrNil()
}

func (d *DeployerService) Delete(ctx context.Context, manifest string) error {
	var result *multierror.Error
	for _, spec := range d.splitter(manifest) {
		typeMeta := metav1.TypeMeta{}
		err := yaml.Unmarshal([]byte(spec), &typeMeta)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("invalid object, can't parse object type metadata: %w\n", err))
			continue
		}
		key := fmt.Sprintf("%s/%s", typeMeta.APIVersion, typeMeta.Kind)
		switch key {
		case "core.k8s.kubemq.io/v1alpha1/KubemqCluster":
			result = multierror.Append(result, objects.NewCluster(d.cfg).Delete(ctx, manifest))
		case "core.k8s.kubemq.io/v1alpha1/KubemqDashboard":
			result = multierror.Append(result, objects.NewDashboard(d.cfg).Delete(ctx, manifest))
		case "core.k8s.kubemq.io/v1alpha1/KubemqConnector":
			result = multierror.Append(result, objects.NewConnector(d.cfg).Delete(ctx, manifest))
		case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition":
			result = multierror.Append(result, objects.NewCrd(d.cfg).Delete(ctx, manifest))
		case "v1/ServiceAccount":
			result = multierror.Append(result, objects.NewServiceAccount(d.cfg).Delete(ctx, manifest))
		case "v1/ConfigMap":
			result = multierror.Append(result, objects.NewConfigMap(d.cfg).Delete(ctx, manifest))
		case "v1/Service":
			result = multierror.Append(result, objects.NewService(d.cfg).Delete(ctx, manifest))
		case "apps/v1/Deployment":
			result = multierror.Append(result, objects.NewDeployment(d.cfg).Delete(ctx, manifest))
		}
	}

	if result == nil {
		return nil
	}
	return result.ErrorOrNil()
}

func (d *DeployerService) splitter(manifest string) []string {
	var validSpecs []string
	for _, spec := range strings.Split(manifest, "---") {
		typeMeta := metav1.TypeMeta{}
		err := yaml.Unmarshal([]byte(spec), &typeMeta)
		fmt.Println(spec, err)
		if err == nil {
			validSpecs = append(validSpecs, spec)
		}
	}

	return validSpecs
}
