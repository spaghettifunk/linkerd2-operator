package trafficsplit

import (
	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	linkerdv1alpha1 "github.com/spaghettifunk/linkerd2-operator/api/v1alpha1"
	"github.com/spaghettifunk/linkerd2-operator/pkg/k8sutil"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"

	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	componentName = "trafficsplits.split.smi-spec.io"
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
		{Resource: r.crd, DesiredState: desiredState},
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

func (r *Reconciler) crd() runtime.Object {
	return &extensionv1.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name:        componentName,
			Namespace:   r.Config.Namespace,
			Annotations: templates.DefaultAnnotations(string(r.Config.Spec.Version)),
			Labels:      r.labels(),
		},
		Spec: extensionv1.CustomResourceDefinitionSpec{
			Group:   "split.smi-spec.io",
			Version: "v1alpha1",
			Scope:   extensionv1.NamespaceScoped,
			Names: extensionv1.CustomResourceDefinitionNames{
				Plural:     "trafficsplits",
				Singular:   "trafficsplit",
				Kind:       "TrafficSplit",
				ShortNames: []string{"ts"},
			},
			AdditionalPrinterColumns: []extensionv1.CustomResourceColumnDefinition{
				{
					Name:        "Service",
					Type:        "string",
					Description: "the apex service of this split.",
					JSONPath:    ".spec.service",
				},
			},
		},
	}
}
