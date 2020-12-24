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
	LimitsEphemeralStorage string `json:"limitsEphemeralStorage,omitempty"`

	// +optional
	RequestsCpu string `json:"requestsCpu,omitempty"`
	// +optional
	RequestsMemory string `json:"requestsMemory,omitempty"`

	// +optional
	RequestsEphemeralStorage string `json:"requestsEphemeralStorage,omitempty"`
}

func (o *ResourceConfig) SetConfig(config *deployment.Config) *ResourceConfig {
	tmpl := `          resources:
            limits:	
              cpu: %s
              memory: %s
              ephemeral-storage: %s
            requests:
              cpu: %s
              memory: %s
              ephemeral-storage: %s
`

	resources := fmt.Sprintf(tmpl,
		o.LimitsCpu,
		o.LimitsMemory,
		o.LimitsEphemeralStorage,
		o.RequestsCpu,
		o.RequestsMemory,
		o.RequestsEphemeralStorage)
	config.StatefulSet.SetResources(resources)
	return o
}
