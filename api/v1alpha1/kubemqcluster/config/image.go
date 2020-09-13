package config

import (
	"github.com/kubemq-io/k8s/api/v1alpha1/kubemqcluster/deployment"
	"os"
)

const fallbackImage = "docker.io/kubemq/kubemq:latest"

type ImageConfig struct {
	// +optional
	Image string `json:"image,omitempty"`
	// +optional
	// +kubebuilder:validation:Pattern=(IfNotPresent|Always|Never)
	PullPolicy string `json:"pullPolicy,omitempty"`
}

func (c *ImageConfig) getFromEnv(currentValue string, envKey string, def string) string {
	if currentValue != "" {
		return currentValue
	}
	fromEnv := os.Getenv(envKey)
	if fromEnv != "" {
		return fromEnv
	}
	return def
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
