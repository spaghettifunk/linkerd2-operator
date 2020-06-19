package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LinkerdSpec defines the desired state of Linkerd
type LinkerdSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Size is the size of the linkerd deployment
	Size     int32  `json:"size"`
	LogLevel string `json:"logLevel"`
}

// LinkerdStatus defines the observed state of Linkerd
type LinkerdStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Linkerd is the Schema for the linkerds API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=linkerds,scope=Namespaced
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
