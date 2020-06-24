package psp

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"

	rbacv1 "k8s.io/api/rbac/v1"
)

func (r *Reconciler) roleBinding() runtime.Object {
	return &rbacv1.RoleBinding{
		ObjectMeta: templates.ObjectMetaClusterScope(roleBindingName, r.labels(), r.Config),
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     roleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "linkerd-controller",
				Namespace: r.Config.Namespace,
			},
			{
				Kind:      "ServiceAccount",
				Name:      "linkerd-destination",
				Namespace: r.Config.Namespace,
			},
			// TODO: fix with grafana from CRD
			{
				Kind:      "ServiceAccount",
				Name:      "grafana",
				Namespace: r.Config.Namespace,
			},
			// TODO: fix with heartbeat from CRD
			{
				Kind:      "ServiceAccount",
				Name:      "linkerd-heartbeat",
				Namespace: r.Config.Namespace,
			},
			{
				Kind:      "ServiceAccount",
				Name:      "linkerd-identity",
				Namespace: r.Config.Namespace,
			},
			{
				Kind:      "ServiceAccount",
				Name:      "linkerd-prometheus",
				Namespace: r.Config.Namespace,
			},
			{
				Kind:      "ServiceAccount",
				Name:      "linkerd-proxy-injector",
				Namespace: r.Config.Namespace,
			},
			{
				Kind:      "ServiceAccount",
				Name:      "linkerd-sp-validator",
				Namespace: r.Config.Namespace,
			},
			{
				Kind:      "ServiceAccount",
				Name:      "linkerd-tap",
				Namespace: r.Config.Namespace,
			},
			{
				Kind:      "ServiceAccount",
				Name:      "linkerd-web",
				Namespace: r.Config.Namespace,
			},
			// TODO: fix with smiMetrics from CRD
			{
				Kind:      "ServiceAccount",
				Name:      "linkerd-smi-metrics",
				Namespace: r.Config.Namespace,
			},
		},
	}
}

func (r *Reconciler) role() runtime.Object {
	return &rbacv1.Role{
		ObjectMeta: templates.ObjectMetaClusterScope(roleName, r.labels(), r.Config),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"policy", "extensions"},
				Resources:     []string{"podsecuritypolicies"},
				Verbs:         []string{"use"},
				ResourceNames: []string{"linkerd-control-plane"},
			},
		},
	}
}
