package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1alpha1/kubemqcluster/deployment"
)

type ResourceConfig struct {
	// +optional
	LimitsCpu string `json:"limitsCpu,omitempty"`
	// +optional
	LimitsMemory string `json:"limitsMemory,omitempty"`
	// +optional
	RequestsCpu string `json:"requestsCpu,omitempty"`
	// +optional
	RequestsMemory string `json:"requestsMemory,omitempty"`
}

func (o *ResourceConfig) SetConfig(config *deployment.Config) *ResourceConfig {
	tmpl := `          resources:
            limits:	
              cpu: %s
              memory: %s
            requests:
              cpu: %s
              memory: %s
`

	resources := fmt.Sprintf(tmpl,
		o.LimitsCpu,
		o.LimitsMemory,
		o.RequestsCpu,
		o.RequestsMemory)
	config.StatefulSet.SetResources(resources)
	return o
}
