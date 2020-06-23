package v1alpha1

import (
	"regexp"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const supportedLinkerdMinorVersionRegex = "^2.*"

// BaseK8sResourceConfiguration defines basic K8s resource spec configurations
type BaseK8sResourceConfiguration struct {
	Image           *string                      `json:"image,omitempty"`
	Resources       *corev1.ResourceRequirements `json:"resources,omitempty"`
	NodeSelector    map[string]string            `json:"nodeSelector,omitempty"`
	Affinity        *corev1.Affinity             `json:"affinity,omitempty"`
	Tolerations     []corev1.Toleration          `json:"tolerations,omitempty"`
	PodAnnotations  map[string]string            `json:"podAnnotations,omitempty"`
	SecurityContext *corev1.SecurityContext      `json:"securityContext,omitempty"`
	ReplicaCount    *int32                       `json:"replicaCount,omitempty"`
}

// ControllerConfiguration defines the k8s spec configuration for the linkerd controller
type ControllerConfiguration struct {
	BaseK8sResourceConfiguration `json:",inline"`
}

// DestinationConfiguration defines the k8s spec configuration for the linkerd destination
type DestinationConfiguration struct {
	BaseK8sResourceConfiguration `json:",inline"`
}

// HeartbeatConfiguration defines the k8s spec configuration for the linkerd heartbeat
// This is a CronJob type
type HeartbeatConfiguration struct {
	BaseK8sResourceConfiguration `json:",inline"`
}

// IdentityConfiguration defines the k8s spec configuration for the linkerd identity
type IdentityConfiguration struct {
	BaseK8sResourceConfiguration `json:",inline"`
}

// PrometheusConfiguration defines the k8s spec configuration for the prometheus deployment
type PrometheusConfiguration struct {
	BaseK8sResourceConfiguration `json:",inline"`
}

// ProxyInitConfiguration defines the k8s spec configuration for the proxy init container
type ProxyInitConfiguration struct {
	BaseK8sResourceConfiguration `json:",inline"`
}

// ProxyInjectorConfiguration defines the k8s spec configuration for the proxy injector
type ProxyInjectorConfiguration struct {
	BaseK8sResourceConfiguration `json:",inline"`
}

// TapConfiguration defines the k8s spec configuration for the linkerd tap
type TapConfiguration struct {
	BaseK8sResourceConfiguration `json:",inline"`
}

// WebConfiguration defines the k8s spec configuration for the linkerd web
type WebConfiguration struct {
	BaseK8sResourceConfiguration `json:",inline"`
}

// LinkerdVersion stores the intended Linkerd version
type LinkerdVersion string

// LinkerdSpec defines the desired state of Linkerd
type LinkerdSpec struct {
	// Contains the intended Linkerd version
	Version LinkerdVersion `json:"version"`
	// Size is the size of the linkerd deployment
	Size int32 `json:"size"`
	// LogLevel is the log level for the linkerd controller
	LogLevel string `json:"logLevel"`
	// List of namespaces to label with sidecar auto injection enabled
	AutoInjectionNamespaces []string `json:"autoInjectionNamespaces,omitempty"`
	// ImagePullPolicy describes a policy for if/when to pull a container image
	// +kubebuilder:validation:Enum=Always;Never;IfNotPresent
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Controller configuration options
	Controller ControllerConfiguration `json:"controller,omitempty"`

	// Destination configuration options
	Destination DestinationConfiguration `json:"destination,omitempty"`

	// Heartbeat configuration options
	Heartbeat HeartbeatConfiguration `json:"heartbeat,omitempty"`

	// Identity configuration options
	Identity IdentityConfiguration `json:"identity,omitempty"`

	// Prometheus configuration options
	Prometheus PrometheusConfiguration `json:"prometheus,omitempty"`

	// ProxyInit configuration options
	ProxyInit ProxyInitConfiguration `json:"proxyInit,omitempty"`

	// ProxyInjector configuration options
	ProxyInjector ProxyInjectorConfiguration `json:"proxyInjector,omitempty"`

	// Tap configuration options
	Tap TapConfiguration `json:"tap,omitempty"`

	// Web configuration options
	Web WebConfiguration `json:"web,omitempty"`
}

// LinkerdStatus defines the observed state of Linkerd
type LinkerdStatus struct {
	Status       ConfigState `json:"Status,omitempty"`
	ErrorMessage string      `json:"ErrorMessage,omitempty"`
}

// IsSupported checks if the version of Linkerd is complied with the supported one by the operator
func (v LinkerdVersion) IsSupported() bool {
	re, _ := regexp.Compile(supportedLinkerdMinorVersionRegex)

	return re.Match([]byte(v))
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Linkerd is the Schema for the linkerds API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.Status",description="Status of the resource"
// +kubebuilder:printcolumn:name="Error",type="string",JSONPath=".status.ErrorMessage",description="Error message"
// +kubebuilder:resource:path=linkerds,scope=Namespaced
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Linkerd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LinkerdSpec   `json:"spec,omitempty"`
	Status LinkerdStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LinkerdList contains a list of Linkerd
type LinkerdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Linkerd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Linkerd{}, &LinkerdList{})
}
