package web

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"

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
		ObjectMeta: templates.ObjectMeta(clusterRoleName, r.labels(), r.Config),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"clusterroles", "clusterrolebindings"},
				Verbs:     []string{"list"},
			},
			{
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"list"},
			},
			{
				APIGroups: []string{"admissionregistration.k8s.io"},
				Resources: []string{"mutatingwebhookconfigurations", "validatingwebhookconfigurations"},
				Verbs:     []string{"list"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"podsecuritypolicies"},
				Verbs:     []string{"list"},
			},
			{
				APIGroups: []string{"linkerd.io"},
				Resources: []string{"serviceprofiles"},
				Verbs:     []string{"list"},
			},
			{
				APIGroups: []string{"apiregistration.k8s.io"},
				Resources: []string{"apiservices"},
				Verbs:     []string{"get"},
			},
		},
	}
}

func (r *Reconciler) clusterRoleBindingWebAdmin() runtime.Object {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: templates.ObjectMeta(clusterRoleBindingNameWebAdmin, r.labels(), r.Config),
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
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

func (r *Reconciler) clusterRoleBindingWebCheck() runtime.Object {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: templates.ObjectMeta(clusterRoleBindingNameWebCheck, r.labels(), r.Config),
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
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

func (r *Reconciler) role() runtime.Object {
	return &rbacv1.Role{
		ObjectMeta: templates.ObjectMeta(roleName, r.labels(), r.Config),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{""},
				Resources:     []string{"configmaps"},
				Verbs:         []string{"get"},
				ResourceNames: []string{"linkerd-config"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces", "configmaps"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"serviceaccounts", "pods"},
				Verbs:     []string{"list"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"replicasets"},
				Verbs:     []string{"list"},
			},
		},
	}
}
