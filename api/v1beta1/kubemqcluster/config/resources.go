package config

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
	"github.com/kubemq-io/k8s/pkg/template"
)

const resourceTmpl = `          resources:
            limits:	
{{ if .LimitsCpu }}
              cpu: {{.LimitsCpu}}
{{end}}
{{ if .LimitsMemory }}
              memory: {{.LimitsMemory}}
{{end}}
{{ if .LimitsEphemeralStorage }}
              ephemeral-storage: {{.LimitsEphemeralStorage}}
{{end}}
            requests:
{{ if .RequestsCpu }}
              cpu: {{.RequestsCpu}}
{{end}}
{{ if .RequestsMemory }}
              memory: {{.RequestsMemory}}
{{end}}
{{ if .RequestsEphemeralStorage }}
              ephemeral-storage: {{.RequestsEphemeralStorage}}
{{end}}
`

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
	t := template.NewTemplate(resourceTmpl, o)
	data, err := t.Get()
	if err != nil {
		fmt.Println(err.Error())
		return o
	}
	config.StatefulSet.SetResources(string(data))
	return o
}
