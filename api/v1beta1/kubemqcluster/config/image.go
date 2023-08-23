package config

import (
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
	"os"
)

const fallbackImage = "docker.io/kubemq/kubemq:latest"

type ImageConfig struct {
	// +optional
	Image string `json:"image,omitempty" yaml:"image,omitempty"`
	// +optional
	// +kubebuilder:validation:Pattern=(IfNotPresent|Always|Never)
	PullPolicy string `json:"pullPolicy,omitempty" yaml:"pullPolicy,omitempty"`
}

func (c *ImageConfig) GetImage() string {
	if c.Image == "" {
		imageFromEnv := os.Getenv("RELATED_IMAGE_KUBEMQ_CLUSTER")
		if imageFromEnv != "" {
			return imageFromEnv
		} else {
			return fallbackImage
		}
	} else {
		return c.Image
	}

}

func (c *ImageConfig) SetConfig(config *deployment.Config) *ImageConfig {
	config.StatefulSet.SetImageName(c.GetImage())
	if c.PullPolicy == "" {
		c.PullPolicy = "Always"
	}
	config.StatefulSet.SetImagePullPolicy(c.PullPolicy)
	return c
}
