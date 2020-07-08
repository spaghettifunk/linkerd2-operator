package heartbeat

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
		},
	}
}

func (r *Reconciler) roleBinding() runtime.Object {
	return &rbacv1.RoleBinding{
		ObjectMeta: templates.ObjectMeta(roleBindingName, r.labels(), r.Config),
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     roleBindingName,
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
