package proxyinjector

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
