package config

import (
	"encoding/json"
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
)

type OIDCConfig struct {
	Issuer                     string `json:"issuer" yaml:"issuer"`
	ClientID                   string `json:"clientID" yaml:"clientID"`
	SkipClientIDCheck          bool   `json:"skipClientIDCheck" yaml:"skipClientIDCheck"`
	SkipExpiryCheck            bool   `json:"skipExpiryCheck" yaml:"skipExpiryCheck"`
	SkipIssuerCheck            bool   `json:"skipIssuerCheck" yaml:"skipIssuerCheck"`
	InsecureSkipSignatureCheck bool   `json:"insecureSkipSignatureCheck" yaml:"insecureSkipSignatureCheck"`
}

func (o *OIDCConfig) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}
func (o *OIDCConfig) Validate() error {
	if o.Issuer == "" {
		return fmt.Errorf("issuer is required")
	}
	if o.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	return nil
}

type AdditionalDataConfiguration struct {
	Oidc *OIDCConfig `json:"oidc,omitempty" yaml:"oidc"`
}

func NewAdditionalDataConfiguration() *AdditionalDataConfiguration {
	return &AdditionalDataConfiguration{}
}

func (adc *AdditionalDataConfiguration) Unmarshal(configData string) error {
	if err := json.Unmarshal([]byte(configData), adc); err != nil {
		return fmt.Errorf("error unmarshaling config data, %s", err.Error())
	}
	if adc.Oidc != nil {
		if err := adc.Oidc.Validate(); err != nil {
			return fmt.Errorf("error validating oidc config, %s", err.Error())
		}
	}
	return nil
}

func (adc *AdditionalDataConfiguration) SetConfig(config *deployment.Config) *AdditionalDataConfiguration {
	if adc.Oidc != nil {
		secConfig, ok := config.ConfigMaps[config.Name]
		if ok {
			if adc.Oidc != nil {
				secConfig.SetStringVariable("AUTHENTICATION_ENABLE", "true").
					SetStringVariable("AUTHENTICATION_TYPE", "oidc").
					SetDataVariable("AUTHENTICATION_CONFIG", adc.Oidc.String())
			}
		}
	}
	return adc
}
