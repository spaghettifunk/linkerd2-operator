package prometheus

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"

	apiv1 "k8s.io/api/core/v1"
)

func (r *Reconciler) service() runtime.Object {
	return &apiv1.Service{
		ObjectMeta: templates.ObjectMetaWithAnnotations(serviceName, r.labels(), templates.DefaultAnnotations(string(r.Config.Spec.Version)), r.Config),
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			// TODO: fix hardcoded values
			Selector: map[string]string{
				"linkerd.io/control-plane-component": componentName,
			},
			Ports: []apiv1.ServicePort{
				templates.DefaultServicePort("admin-http", 9090, 9090),
			},
		},
	}
}
