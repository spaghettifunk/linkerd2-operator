// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the linkerd v1alpha1 API group
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/spaghettifunk/linkerd2-operator/pkg/apis/linkerd
// +k8s:defaulter-gen=TypeMeta
// +groupName=linkerd.io
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "linkerd.io", Version: "v1alpha1"}
)

// Resource is required by pkg/client/listers/...
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
