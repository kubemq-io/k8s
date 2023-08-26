/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KubemqConnectorSpec defines the desired state of KubemqConnector
type KubemqConnectorSpec struct {
	// +optional
	// +kubebuilder:validation:Minimum=0

	Replicas *int32 `json:"replicas,omitempty" yaml:"replicas,omitempty"`

	Type string `json:"type" yaml:"type"`

	// +optional
	Image string `json:"image,omitempty" yaml:"image,omitempty"`

	Config string `json:"config" yaml:"config"`

	// +optional

	NodePort int32 `json:"nodePort,omitempty" yaml:"node_port"`

	// +optional

	ServiceType string `json:"serviceType" yaml:"service_type"`
}

// KubemqConnectorStatus defines the observed state of KubemqConnector
type KubemqConnectorStatus struct {
	Replicas int32 `json:"replicas" yaml:"replicas"`

	Type string `json:"type" yaml:"type"`

	Image string `json:"image" yaml:"image"`

	Api string `json:"api" yaml:"api"`

	Status string `json:"status" yaml:"status"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubemqconnectors,scope=Namespaced
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:printcolumn:JSONPath=".status.replicas",name=Replicas,type=string
// +kubebuilder:printcolumn:JSONPath=".status.type",name=Type,type=string
// +kubebuilder:printcolumn:JSONPath=".status.image",name=Image,type=string
// +kubebuilder:printcolumn:JSONPath=".status.api",name=API,type=string
// +kubebuilder:printcolumn:JSONPath=".status.status",name=Status,type=string

// KubemqConnector is the Schema for the kubemqconnectors API
type KubemqConnector struct {
	metav1.TypeMeta   `json:",inline" yaml:"inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec   KubemqConnectorSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status KubemqConnectorStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubemqConnectorList contains a list of KubemqConnector
type KubemqConnectorList struct {
	metav1.TypeMeta `json:",inline" yaml:"metav1.inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Items           []KubemqConnector `json:"items" yaml:"items"`
}

func init() {
	SchemeBuilder.Register(&KubemqConnector{}, &KubemqConnectorList{})
}
