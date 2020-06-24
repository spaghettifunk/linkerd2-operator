package psp

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Reconciler) deployment() runtime.Object {
	return &policyv1.PodSecurityPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "linkerd-control-plane",
			Labels: r.labels(),
		},
		Spec: policyv1.PodSecurityPolicySpec{
			AllowPrivilegeEscalation: util.BoolPointer(false),
			ReadOnlyRootFilesystem:   true,
			// TODO: fix - if cniEnabled do the below
			AllowedCapabilities: []v1.Capability{
				"NET_ADMIN",
				"NET_RAW",
			},
			RequiredDropCapabilities: []v1.Capability{"ALL"},
			HostNetwork:              false,
			HostIPC:                  false,
			HostPID:                  false,
			SELinux: policyv1.SELinuxStrategyOptions{
				Rule: policyv1.SELinuxStrategyRunAsAny,
			},
			RunAsUser: policyv1.RunAsUserStrategyOptions{
				// TODO: fix - if cniEnabled do the below
				// Rule: policyv1.RunAsUserStrategyMustRunAsNonRoot,
				Rule: policyv1.RunAsUserStrategyRunAsAny,
			},
			SupplementalGroups: policyv1.SupplementalGroupsStrategyOptions{
				Rule: policyv1.SupplementalGroupsStrategyMustRunAs,
				Ranges: []policyv1.IDRange{
					// TODO: fix - if cniEnabled do the below
					// {
					// 	Min: int64(10001),
					// 	Max: int64(65535),
					// },
					{
						Min: int64(1),
						Max: int64(65535),
					},
				},
			},
			FSGroup: policyv1.FSGroupStrategyOptions{
				Rule: policyv1.FSGroupStrategyMustRunAs,
				Ranges: []policyv1.IDRange{
					// TODO: fix - if cniEnabled do the below
					// {
					// 	Min: int64(10001),
					// 	Max: int64(65535),
					// },
					{
						Min: int64(1),
						Max: int64(65535),
					},
				},
			},
			Volumes: []policyv1.FSType{
				policyv1.ConfigMap,
				policyv1.EmptyDir,
				policyv1.Secret,
				policyv1.Projected,
				policyv1.DownwardAPI,
				policyv1.PersistentVolumeClaim,
			},
		},
	}
}
