package tap

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kube-aggregator/pkg/apis/apiregistration"

	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

func (r *Reconciler) serviceAccount() runtime.Object {
	return &apiv1.ServiceAccount{
		ObjectMeta: templates.ObjectMeta(serviceAccountName, r.labels(), r.Config),
	}
}

func (r *Reconciler) roleBindingAuthReader() runtime.Object {
	return &rbacv1.RoleBinding{
		ObjectMeta: templates.ObjectMeta(roleBindingNameAuthReader, r.labels(), r.Config),
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     "extension-apiserver-authentication-reader",
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

func (r *Reconciler) clusterRole() runtime.Object {
	return &rbacv1.ClusterRole{
		ObjectMeta: templates.ObjectMetaClusterScope(clusterRoleName, r.labels(), r.Config),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "services", "replicationcontrollers", "namespaces", "nodes"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"extensions", "apps"},
				Resources: []string{"daemonsets", "deployments", "replicasets", "statefulsets"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"extensions", "batch"},
				Resources: []string{"cronjobs", "jobs"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	}
}

func (r *Reconciler) clusterRoleAdmin() runtime.Object {
	return &rbacv1.ClusterRole{
		ObjectMeta: templates.ObjectMetaClusterScope(clusterRoleNameAdmin, r.labels(), r.Config),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"tap.linkerd.io"},
				Resources: []string{"*"},
				Verbs:     []string{"watch"},
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

func (r *Reconciler) clusterRoleBindingAuthDelegator() runtime.Object {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: templates.ObjectMetaClusterScope(clusterRoleBindingNameAuthDelegator, r.labels(), r.Config),
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "system:auth-delegator",
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

func (r *Reconciler) apiService() runtime.Object {
	return &apiregistration.APIService{
		ObjectMeta: templates.ObjectMeta(apiServiceName, r.labels(), r.Config),
		Spec: apiregistration.APIServiceSpec{
			Group:                "tap.linkerd.io",
			Version:              "v1alpha1",
			GroupPriorityMinimum: int32(1000),
			VersionPriority:      int32(100),
			Service: &apiregistration.ServiceReference{
				Name:      "linkerd-tap",
				Namespace: r.Config.Namespace,
			},
			CABundle: []byte(r.Config.Spec.SelfSignedCertificates.TrustAnchorsPEM),
		},
	}
}
