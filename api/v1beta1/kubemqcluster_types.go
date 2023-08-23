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
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KubemqClusterSpec defines the desired state of KubemqCluster
type KubemqClusterSpec struct {
	// +optional
	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty" yaml:"replicas,omitempty"`

	// +optional
	ConfigData string `json:"configData,omitempty" yaml:"configData,omitempty"`

	// +optional
	License string `json:"license,omitempty" yaml:"license,omitempty"`

	// +optional
	Key string `json:"key,omitempty" yaml:"key,omitempty"`

	// +optional
	Standalone bool `json:"standalone,omitempty" yaml:"standalone,omitempty"`

	// +optional
	Volume *config.VolumeConfig `json:"volume,omitempty" yaml:"volume,omitempty"`

	// +optional
	Image *config.ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	// +optional
	Api *config.ApiConfig `json:"api,omitempty" yaml:"api,omitempty"`

	// +optional
	Rest *config.RestConfig `json:"rest,omitempty" yaml:"rest,omitempty"`

	// +optional
	Grpc *config.GrpcConfig `json:"grpc,omitempty" yaml:"grpc,omitempty"`

	// +optional
	Tls *config.TlsConfig `json:"tls,omitempty" yaml:"tls,omitempty"`

	// +optional
	Resources *config.ResourceConfig `json:"resources,omitempty" yaml:"resources,omitempty"`

	// +optional
	NodeSelectors *config.NodeSelectorConfig `json:"nodeSelectors,omitempty" yaml:"nodeSelectors,omitempty"`

	// +optional
	Authentication *config.AuthenticationConfig `json:"authentication,omitempty" yaml:"authentication,omitempty"`

	// +optional
	Authorization *config.AuthorizationConfig `json:"authorization,omitempty" yaml:"authorization,omitempty"`

	// +optional
	Health *config.HealthConfig `json:"health,omitempty" yaml:"health,omitempty"`

	// +optional
	Routing *config.RoutingConfig `json:"routing,omitempty" yaml:"routing,omitempty"`

	// +optional
	Log *config.LogConfig `json:"log,omitempty" yaml:"log,omitempty"`

	// +optional
	Notification *config.NotificationConfig `json:"notification,omitempty" yaml:"notification,omitempty"`

	// +optional
	Store *config.StoreConfig `json:"store,omitempty" yaml:"store,omitempty"`

	// +optional
	Queue *config.QueueConfig `json:"queue,omitempty" yaml:"queue,omitempty"`

	// +optional
	StatefulSetConfigData string `json:"statefulsetConfigData,omitempty" yaml:"statefulSetConfigData,omitempty"`
}

// KubemqClusterStatus defines the observed state of KubemqCluster
type KubemqClusterStatus struct {
	Replicas *int32 `json:"replicas" yaml:"replicas"`

	Version string `json:"version" yaml:"version"`

	Ready int32 `json:"ready" yaml:"ready"`

	Grpc string `json:"grpc" yaml:"grpc"`

	Rest string `json:"rest" yaml:"rest"`

	Api string `json:"api" yaml:"api"`

	Selector string `json:"selector" yaml:"selector"`

	LicenseType string `json:"license_type" yaml:"licenseType"`

	LicenseTo string `json:"license_to" yaml:"licenseTo"`

	LicenseExpire string `json:"license_expire" yaml:"licenseExpire"`

	Status string `json:"status" yaml:"status"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubemqclusters,scope=Namespaced
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:printcolumn:JSONPath=".status.version",name=Version,type=string
// +kubebuilder:printcolumn:JSONPath=".status.status",name=Status,type=string
// +kubebuilder:printcolumn:JSONPath=".status.replicas",name=Replicas,type=string
// +kubebuilder:printcolumn:JSONPath=".status.ready",name=Ready,type=string
// +kubebuilder:printcolumn:JSONPath=".status.grpc",name=gRPC,type=string
// +kubebuilder:printcolumn:JSONPath=".status.rest",name=Rest,type=string
// +kubebuilder:printcolumn:JSONPath=".status.api",name=API,type=string
// +kubebuilder:printcolumn:JSONPath=".status.license_type",name=License-type,type=string
// +kubebuilder:printcolumn:JSONPath=".status.license_to",name=License-To,type=string
// +kubebuilder:printcolumn:JSONPath=".status.license_expire",name=License-Expire,type=string

type KubemqCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubemqClusterSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status KubemqClusterStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubemqClusterList contains a list of KubemqCluster
type KubemqClusterList struct {
	metav1.TypeMeta `json:",inline" yaml:"inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Items           []KubemqCluster `json:"items" yaml:"items"`
}

func init() {
	SchemeBuilder.Register(&KubemqCluster{}, &KubemqClusterList{})
}
