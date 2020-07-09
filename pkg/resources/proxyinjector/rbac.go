package proxyinjector

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

func (r *Reconciler) serviceAccount() runtime.Object {
	return &apiv1.ServiceAccount{
		ObjectMeta: templates.ObjectMeta(serviceAccountName, r.labels(), r.Config),
	}
}

func (r *Reconciler) clusterRole() runtime.Object {
	return &rbacv1.ClusterRole{
		ObjectMeta: templates.ObjectMetaClusterScope(clusterRoleName, r.labels(), r.Config),
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"create", "patch"},
				APIGroups: []string{""},
				Resources: []string{"events"},
			},
			{
				Verbs:     []string{"get", "list", "watch"},
				APIGroups: []string{""},
				Resources: []string{"namespaces", "replicationcontrollers"},
			},
			{
				Verbs:     []string{"list", "watch"},
				APIGroups: []string{""},
				Resources: []string{"pods"},
			},
			{
				Verbs:     []string{"list", "get", "watch"},
				APIGroups: []string{"extensions", "apps"},
				Resources: []string{"deployments", "replicasets", "daemonsets", "statefulsets"},
			},
			{
				Verbs:     []string{"list", "get", "watch"},
				APIGroups: []string{"extensions", "batch"},
				Resources: []string{"cronjobs", "jobs"},
			},
		},
	}
}

func (r *Reconciler) clusterRoleBinding() runtime.Object {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: templates.ObjectMetaClusterScope(clusterRoleBindingName, r.labels(), r.Config),
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: r.Config.Namespace,
			},
		},
	}
}

func (r *Reconciler) mutatingWebhookConfiguration() runtime.Object {
	ignore := admissionregistrationv1beta1.Ignore
	none := admissionregistrationv1beta1.SideEffectClassNone
	return &admissionregistrationv1beta1.MutatingWebhookConfiguration{
		ObjectMeta: templates.ObjectMetaClusterScope(mutatingWebhookConfiguration, r.labels(), r.Config),
		Webhooks: []admissionregistrationv1beta1.MutatingWebhook{
			{
				AdmissionReviewVersions: []string{"v1beta1"},
				Name:                    "linkerd-proxy-injector.linkerd.io",
				NamespaceSelector: &v1.LabelSelector{
					MatchExpressions: []v1.LabelSelectorRequirement{
						{
							Key:      "config.linkerd.io/admission-webhooks",
							Operator: "NotIn",
							Values:   []string{"disabled"},
						},
					},
				},
				ClientConfig: admissionregistrationv1beta1.WebhookClientConfig{
					Service: &admissionregistrationv1beta1.ServiceReference{
						Name:      "linkerd-proxy-injector",
						Namespace: r.Config.Namespace,
						Path:      util.StrPointer("/"),
					},
					CABundle: []byte(r.Config.Spec.SelfSignedCertificates.TrustAnchorsPEM),
				},
				FailurePolicy: &ignore,
				SideEffects:   &none,
				Rules: []admissionregistrationv1beta1.RuleWithOperations{
					{
						Operations: []admissionregistrationv1beta1.OperationType{"CREATE"},
						Rule: admissionregistrationv1beta1.Rule{
							APIGroups:   []string{""},
							APIVersions: []string{"v1"},
							Resources:   []string{"pods"},
						},
					},
				},
			},
		},
	}
}
