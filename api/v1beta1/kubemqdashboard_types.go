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

// KubemqDashboardSpec defines the desired state of KubemqDashboard
type KubemqDashboardSpec struct {
	// +optional
	Port int32 `json:"port,omitempty"`

	// +optional
	Prometheus *PrometheusConfig `json:"prometheus,omitempty"`

	// +optional
	Grafana *GrafanaConfig `json:"grafana,omitempty"`
}

// KubemqDashboardStatus defines the observed state of KubemqDashboard
type KubemqDashboardStatus struct {
	Status            string `json:"status"`
	Address           string `json:"address"`
	PrometheusVersion string `json:"prometheus_version"`
	GrafanaVersion    string `json:"grafana_version"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubemqdashboards,scope=Namespaced
// +kubebuilder:printcolumn:JSONPath=".status.status",name=Status,type=string
// +kubebuilder:printcolumn:JSONPath=".status.address",name=Address,type=string
// +kubebuilder:printcolumn:JSONPath=".status.prometheus_version",name=Prometheus-Version,type=string
// +kubebuilder:printcolumn:JSONPath=".status.grafana_version",name=Grafana-Version,type=string

// KubemqDashboard is the Schema for the kubemqdashboards API
type KubemqDashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubemqDashboardSpec   `json:"spec,omitempty"`
	Status KubemqDashboardStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type KubemqDashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubemqDashboard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubemqDashboard{}, &KubemqDashboardList{})
}

type PrometheusConfig struct {
	// +optional
	NodePort int32 `json:"nodePort,omitempty"`
	// +optional
	Image string `json:"image,omitempty"`
}

type GrafanaConfig struct {
	// +optional
	DashboardUrl string `json:"dashboardUrl,omitempty"`
	// +optional
	Image string `json:"image,omitempty"`
}
