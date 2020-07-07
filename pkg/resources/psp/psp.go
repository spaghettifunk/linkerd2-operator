package psp

import (
	"github.com/go-logr/logr"
	"github.com/goph/emperror"

	linkerdv1alpha1 "github.com/spaghettifunk/linkerd2-operator/api/v1alpha1"
	"github.com/spaghettifunk/linkerd2-operator/pkg/k8sutil"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	componentName   = "linkerd-control-plane"
	roleName        = "linkerd-psp"
	roleBindingName = "linkerd-psp"
)

// Reconciler .
type Reconciler struct {
	resources.Reconciler
}

// New .
func New(client client.Client, config *linkerdv1alpha1.Linkerd) *Reconciler {
	return &Reconciler{
		Reconciler: resources.Reconciler{
			Client: client,
			Config: config,
		},
	}
}

// Reconcile .
func (r *Reconciler) Reconcile(log logr.Logger) error {
	log = log.WithValues("component", componentName)

	desiredState := k8sutil.DesiredStatePresent

	log.Info("Reconciling")

	for _, res := range []resources.ResourceWithDesiredState{
		{Resource: r.role, DesiredState: desiredState},
		{Resource: r.roleBinding, DesiredState: desiredState},
	} {
		o := res.Resource()
		err := k8sutil.Reconcile(log, r.Client, o, res.DesiredState)
		if err != nil {
			return emperror.WrapWith(err, "failed to reconcile resource", "resource", o.GetObjectKind().GroupVersionKind())
		}
	}

	log.Info("Reconciled")

	return nil
}

func (r *Reconciler) labels() map[string]string {
	return map[string]string{
		"linkerd.io/control-plane-ns": r.Config.Namespace,
	}
}
