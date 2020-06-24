package web

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	apiv1 "k8s.io/api/core/v1"
)

func (r *Reconciler) service() runtime.Object {
	return &apiv1.Service{
		ObjectMeta: templates.ObjectMetaWithAnnotations(serviceName, r.labels(), r.annotations(), r.Config),
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			// TODO: remove hardcoded values
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       int32(8084),
					TargetPort: intstr.FromString("8084"),
				},
				{
					Name:       "admin-http",
					Port:       int32(9994),
					TargetPort: intstr.FromString("9994"),
				},
			},
		},
	}
}
