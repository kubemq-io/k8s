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

// EstablishedEngineAnnotation records the persistence engine ("legacy" or "next")
// the operator has established for this cluster's on-disk store. Consulted by the
// operator's fail-closed derivation guard on every reconcile; once set it is
// authoritative over spec.store.engine for guard purposes.
const EstablishedEngineAnnotation = "core.k8s.kubemq.io/established-engine"

// OperatorComputedEnvKeyPrefix is the reserved prefix for env var names the
// operator computes and injects for next-engine cluster replication. Keys under
// this prefix (and OperatorComputedEnvKeys) may not be set via spec.env.
const OperatorComputedEnvKeyPrefix = "CLUSTER_REPLICATION_"

// OperatorComputedEnvKeys lists the individual (non-prefixed) env var names the
// operator computes and injects. Keys in this list may not be set via spec.env.
// KUBEMQ_TERMINATION_GRACE_PERIOD_SECONDS is denied UNCONDITIONALLY: it describes an
// external constraint the pod cannot read for itself (Kubernetes does not expose
// terminationGracePeriodSeconds through the downward API), so setting it from spec.env
// — with or without the matching pod-spec field — can only make the server believe a
// shutdown budget it does not have. spec.terminationGracePeriodSeconds writes both.
//
// The Kafka addressing keys are NOT here: the operator computes them only for a
// clustered cluster with the connector enabled, and denying them everywhere would fail
// the reconcile of an existing cluster that reaches them through spec.env. The operator
// rejects them contextually instead, on the path where it computes them.
var OperatorComputedEnvKeys = []string{
	"STORE_ENGINE", "CLUSTER_ENABLE", "CLUSTER_NAME", "CLUSTER_ROUTES",
	"API_BIND_ADDRESS", "CHECKSUM", "POD_NAME",
	"KUBEMQ_TERMINATION_GRACE_PERIOD_SECONDS",
}

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
	// KeySecretRef names an existing Secret holding the license key, so the raw key
	// need not be a literal in the CR / Helm values. Mutually exclusive with Key.
	KeySecretRef *string `json:"keySecretRef,omitempty" yaml:"keySecretRef,omitempty"`

	// +optional
	// KeySecretKey is the data key within KeySecretRef holding the license key.
	// Defaults to "key".
	KeySecretKey *string `json:"keySecretKey,omitempty" yaml:"keySecretKey,omitempty"`

	// +optional
	// LicenseSecretRef names an existing Secret holding the license data, so the raw
	// license need not be a literal in the CR / Helm values. Mutually exclusive with License.
	LicenseSecretRef *string `json:"licenseSecretRef,omitempty" yaml:"licenseSecretRef,omitempty"`

	// +optional
	// LicenseSecretKey is the data key within LicenseSecretRef holding the license data.
	// Defaults to "license".
	LicenseSecretKey *string `json:"licenseSecretKey,omitempty" yaml:"licenseSecretKey,omitempty"`

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
	Mcp *config.McpConfig `json:"mcp,omitempty" yaml:"mcp,omitempty"`

	// +optional
	Agents *config.AgentsConfig `json:"agents,omitempty" yaml:"agents,omitempty"`

	// +optional
	Ce *config.CeConfig `json:"ce,omitempty" yaml:"ce,omitempty"`

	// +optional
	Mqtt *config.MqttConfig `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`

	// +optional
	Amqp *config.AmqpConfig `json:"amqp,omitempty" yaml:"amqp,omitempty"`

	// +optional
	Amqp10 *config.Amqp10Config `json:"amqp10,omitempty" yaml:"amqp10,omitempty"`

	// +optional
	Stomp *config.StompConfig `json:"stomp,omitempty" yaml:"stomp,omitempty"`

	// +optional
	Aws *config.AwsConfig `json:"aws,omitempty" yaml:"aws,omitempty"`

	// +optional
	Gcp *config.GcpConfig `json:"gcp,omitempty" yaml:"gcp,omitempty"`

	// +optional
	Kafka *config.KafkaConfig `json:"kafka,omitempty" yaml:"kafka,omitempty"`

	// +optional
	Telemetry *config.TelemetryConfig `json:"telemetry,omitempty" yaml:"telemetry,omitempty"`

	// +optional
	Audit *config.AuditConfig `json:"audit,omitempty" yaml:"audit,omitempty"`

	// +optional
	Http *config.HttpConfig `json:"http,omitempty" yaml:"http,omitempty"`

	// TerminationGracePeriodSeconds is the pod's shutdown budget. The server splits it
	// between connector teardown (which requeues in-flight messages) and store
	// shutdown, and it cannot read the value itself — Kubernetes does not expose
	// terminationGracePeriodSeconds through the downward API. The operator therefore
	// writes it to BOTH the pod spec and the KUBEMQ_TERMINATION_GRACE_PERIOD_SECONDS
	// env key. Unset leaves the Kubernetes default (30s), which is what the server
	// assumes when it is not told.
	// +optional
	// +kubebuilder:validation:Minimum=1
	TerminationGracePeriodSeconds *int32 `json:"terminationGracePeriodSeconds,omitempty" yaml:"terminationGracePeriodSeconds,omitempty"`

	// +optional
	PodDisruptionBudget *config.PodDisruptionBudgetConfig `json:"podDisruptionBudget,omitempty" yaml:"podDisruptionBudget,omitempty"`

	// +optional
	PodAntiAffinity *config.PodAntiAffinityConfig `json:"podAntiAffinity,omitempty" yaml:"podAntiAffinity,omitempty"`

	// +optional
	StatefulSetConfigData string `json:"statefulsetConfigData,omitempty" yaml:"statefulsetConfigData,omitempty"`

	// Env sets additional environment variables on the StatefulSet's ConfigMap
	// overlay. Keys under OperatorComputedEnvKeyPrefix or in OperatorComputedEnvKeys
	// are rejected by the operator (reserved for engine/cluster wiring).
	// +optional
	Env map[string]string `json:"env,omitempty" yaml:"env,omitempty"`

	// EnvFromSecrets lists additional Secret names whose keys are projected as
	// environment variables on the StatefulSet pods, alongside the operator-owned
	// ConfigMap envFrom source.
	// +optional
	EnvFromSecrets []string `json:"envFromSecrets,omitempty" yaml:"envFromSecrets,omitempty"`
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

	// Engine reports the persistence engine ("legacy" or "next") the operator has
	// established for this cluster. Display-only; the operator derives and stamps
	// the effective engine via EstablishedEngineAnnotation, not this field.
	// +optional
	Engine string `json:"engine,omitempty" yaml:"engine,omitempty"`

	// Conditions represent the latest available observations of the cluster's state.
	// Known .status.conditions.type are: "ReconcileError" and "EphemeralNextStore".
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" yaml:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
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

// GetConditions returns the cluster's status conditions. Satisfies the operator's
// apis.ConditionsAware interface.
func (in *KubemqCluster) GetConditions() []metav1.Condition {
	return in.Status.Conditions
}

// SetConditions replaces the cluster's status conditions. Satisfies the operator's
// apis.ConditionsAware interface.
func (in *KubemqCluster) SetConditions(conditions []metav1.Condition) {
	in.Status.Conditions = conditions
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
