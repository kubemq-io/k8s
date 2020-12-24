package deployment

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceConfig_Spec(t *testing.T) {
	type fields struct {
		Id            string
		Name          string
		Namespace     string
		AppName       string
		Type          string
		ContainerPort int32
		TargetPort    int32
		PortName      string
		Headless      bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "full",
			fields: fields{
				Id:            "some-id",
				Name:          "kubemq",
				Namespace:     "kubemq-namespace",
				AppName:       "svc",
				Type:          "NodePort",
				ContainerPort: 5000,
				TargetPort:    6000,
				PortName:      "kube-rest",
				Headless:      false,
			},
			wantErr: false,
		},
		{
			name: "full",
			fields: fields{
				Id:            "some-id",
				Name:          "kubemq",
				Namespace:     "kubemq-namespace",
				AppName:       "svc",
				Type:          "NodePort",
				ContainerPort: 5000,
				TargetPort:    6000,
				PortName:      "kube-rest",
				Headless:      true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ServiceConfig{
				Id:            tt.fields.Id,
				Name:          tt.fields.Name,
				Namespace:     tt.fields.Namespace,
				AppName:       tt.fields.AppName,
				Expose:        tt.fields.Type,
				ContainerPort: tt.fields.ContainerPort,
				TargetPort:    tt.fields.TargetPort,
				PortName:      tt.fields.PortName,
				NodePort:      0,
				Headless:      tt.fields.Headless,
				service:       nil,
			}
			svc, err := s.Get()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.EqualValues(t, tt.fields.Name, svc.Name)
				data, _ := yaml.Marshal(svc)
				fmt.Println(string(data))
			}
		})
	}
}
