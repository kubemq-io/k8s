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

func TestDefaultServiceConfig_ConnectorServices(t *testing.T) {
	list := DefaultServiceConfig("some-id", "kubemq-namespace", "svc")

	// the 4 connector services must be present (default-on)
	for _, key := range []string{"mqtt", "amqp", "stomp", "aws"} {
		_, ok := list[key]
		assert.Truef(t, ok, "expected connector service %q to be present", key)
	}

	tests := []struct {
		key      string
		wantName string
		ports    []ServicePort
	}{
		{
			key:      "mqtt",
			wantName: "svc-mqtt",
			ports: []ServicePort{
				{Name: "mqtt", Port: 1883, TargetPort: 1883},
				{Name: "mqtt-tls", Port: 8883, TargetPort: 8883},
				{Name: "mqtt-ws", Port: 8083, TargetPort: 8083},
			},
		},
		{
			key:      "amqp",
			wantName: "svc-amqp",
			ports: []ServicePort{
				{Name: "amqp", Port: 5672, TargetPort: 5672},
				{Name: "amqp-tls", Port: 5671, TargetPort: 5671},
			},
		},
		{
			key:      "stomp",
			wantName: "svc-stomp",
			ports: []ServicePort{
				{Name: "stomp", Port: 61613, TargetPort: 61613},
				{Name: "stomp-tls", Port: 61614, TargetPort: 61614},
			},
		},
		{
			key:      "aws",
			wantName: "svc-aws",
			ports: []ServicePort{
				{Name: "aws-http", Port: 4566, TargetPort: 4566},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			cfg := list[tt.key]
			require.NotNil(t, cfg)

			svc, err := cfg.Get()
			require.NoError(t, err)

			assert.EqualValues(t, tt.wantName, svc.Name)
			assert.EqualValues(t, "kubemq-namespace", svc.Namespace)
			assert.EqualValues(t, "ClusterIP", string(svc.Spec.Type))

			// exactly the expected ports, in the fixed order
			require.Len(t, svc.Spec.Ports, len(tt.ports))
			for i, want := range tt.ports {
				got := svc.Spec.Ports[i]
				assert.EqualValues(t, want.Name, got.Name)
				assert.EqualValues(t, want.Port, got.Port)
				assert.EqualValues(t, "TCP", string(got.Protocol))
				assert.EqualValues(t, want.TargetPort, got.TargetPort.IntValue())
				// targetPort must equal port and no nodePort must be set
				assert.EqualValues(t, got.Port, got.TargetPort.IntValue())
				assert.EqualValues(t, int32(0), got.NodePort)
			}

			data, _ := yaml.Marshal(svc)
			fmt.Println(string(data))
		})
	}
}
