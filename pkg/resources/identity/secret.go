package identity

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/certs"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"

	apiv1 "k8s.io/api/core/v1"
)

func (r *Reconciler) secret() runtime.Object {
	return &apiv1.Secret{
		ObjectMeta: templates.ObjectMetaWithAnnotations(
			secretName,
			r.labels(),
			map[string]string{
				"linkerd.io/identity-issuer-expiry": certs.DefaultLifetime.String(),
			},
			r.Config,
		),
		Data: map[string][]byte{
			"crt.pem": []byte(r.Config.Spec.SelfSignedCertificates.CrtPEM),
			"key.pem": []byte(r.Config.Spec.SelfSignedCertificates.KeyPEM),
		},
	}
}
