package prometheus

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) configmap() runtime.Object {
	return &apiv1.ConfigMap{
		ObjectMeta: templates.ObjectMetaWithAnnotations(serviceName, r.labels(), templates.DefaultAnnotations(string(r.Config.Spec.Version)), r.Config),
		Data: map[string]string{
			// TODO: load Prometheus config file here
			"prometheus.yaml": "",
		},
	}
}
