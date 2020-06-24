package tap

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"

	apiv1 "k8s.io/api/core/v1"
)

func (r *Reconciler) secret() runtime.Object {
	return &apiv1.Secret{
		ObjectMeta: templates.ObjectMetaWithAnnotations(secretName, r.labels(), templates.DefaultAnnotations(string(r.Config.Spec.Version)), r.Config),
		// TODO: fix the certificates
		Data: map[string][]byte{
			"crt.pem": nil,
			"key.pem": nil,
		},
	}
}
