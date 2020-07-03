package prometheus

import (
	linkerdv1alpha1 "github.com/spaghettifunk/linkerd2-operator/pkg/apis/linkerd/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	"github.com/spaghettifunk/linkerd2-operator/pkg/k8sutil"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	componentName          = "prometheus"
	serviceAccountName     = "linkerd-prometheus"
	clusterRoleName        = "linkerd-prometheus"
	clusterRoleBindingName = "linkerd-prometheus"
	deploymentName         = "linkerd-prometheus"
	serviceName            = "linkerd-prometheus"
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
		{Resource: r.serviceAccount, DesiredState: desiredState},
		{Resource: r.clusterRole, DesiredState: desiredState},
		{Resource: r.clusterRoleBinding, DesiredState: desiredState},
		{Resource: r.configmap, DesiredState: desiredState},
		{Resource: r.deployment, DesiredState: desiredState},
		{Resource: r.service, DesiredState: desiredState},
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
		"linkerd.io/control-plane-component": componentName,
		"linkerd.io/control-plane-ns":        r.Config.Namespace,
	}
}

// deploymentLabels returns the labels used for the deployment of the web component
func (r *Reconciler) deploymentLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":    componentName,
		"app.kubernetes.io/part-of": "Linkerd",
		"app.kubernetes.io/version": string(r.Config.Spec.Version),
	}
}

func (r *Reconciler) annotations() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name": componentName,
	}
}
