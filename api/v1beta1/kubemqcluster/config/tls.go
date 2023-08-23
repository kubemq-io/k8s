package config

import "github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"

type TlsConfig struct {
	// +optional
	Cert string `json:"cert,omitempty" yaml:"cert,omitempty"`

	// +optional
	Key string `json:"key,omitempty" yaml:"key,omitempty"`

	// +optional
	Ca string `json:"ca,omitempty" yaml:"ca"`
}

func (o *TlsConfig) SetConfig(config *deployment.Config) *TlsConfig {
	secConfig, ok := config.Secrets[config.Name]
	if ok {
		if o.Cert != "" {
			secConfig.SetDataVariable("SECURITY_CERT_DATA", o.Cert)
		}
		if o.Key != "" {
			secConfig.SetDataVariable("SECURITY_KEY_DATA", o.Key)
		}
		if o.Ca != "" {
			secConfig.SetDataVariable("SECURITY_CA_DATA", o.Ca)
		}
	}
	return o
}
