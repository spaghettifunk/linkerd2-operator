package resources

import (
	"github.com/go-logr/logr"
	"github.com/spaghettifunk/linkerd2-operator/pkg/k8sutil"
	"k8s.io/apimachinery/pkg/runtime"

	linkerdv1alpha1 "github.com/spaghettifunk/linkerd2-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceWithDesiredState defines the desidered state based on the resources
type ResourceWithDesiredState struct {
	Resource     Resource
	DesiredState k8sutil.DesiredState
}

// ResourceVariationWithDesiredState defines the desidered state based on the variation of the resources
type ResourceVariationWithDesiredState struct {
	ResourceVariation ResourceVariation
	DesiredState      k8sutil.DesiredState
}

// Reconciler is the object holding the client and the configuration of the operator
type Reconciler struct {
	client.Client
	Config *linkerdv1alpha1.Linkerd
}

// ComponentReconciler is the interface that is used for each sub-component to reconcile with the config
type ComponentReconciler interface {
	Reconcile(log logr.Logger) error
}

// Resource defines a runtime.Object type
type Resource func() runtime.Object

// ResourceVariation defines a runtime.Object type
type ResourceVariation func(t string) runtime.Object

// ResolveVariations takes the desired state and try to match it
func ResolveVariations(t string, v []ResourceVariationWithDesiredState, desiredState k8sutil.DesiredState) []ResourceWithDesiredState {
	var state k8sutil.DesiredState
	resources := make([]ResourceWithDesiredState, 0)
	for i := range v {
		i := i
		if v[i].DesiredState == k8sutil.DesiredStateAbsent || desiredState == k8sutil.DesiredStateAbsent {
			state = k8sutil.DesiredStateAbsent
		} else {
			state = k8sutil.DesiredStatePresent
		}
		resource := ResourceWithDesiredState{
			func() runtime.Object {
				return v[i].ResourceVariation(t)
			},
			state,
		}
		resources = append(resources, resource)
	}
	return resources
}
