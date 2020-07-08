package heartbeat

import (
	"fmt"

	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"

	batchv1 "k8s.io/api/batch/v1"
	v1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Reconciler) cronjob() runtime.Object {
	return &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cronjobName,
			Namespace: r.Config.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":             "heartbeat",
				"app.kubernetes.io/part-of":          "Linkerd",
				"app.kubernetes.io/version":          "stable-2.8.1",
				"linkerd.io/control-plane-component": "heartbeat",
				"linkerd.io/control-plane-ns":        "linkerd",
			},
			Annotations: map[string]string{
				"linkerd.io/created-by": "linkerd/cli stable-2.8.1",
			},
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: "16 8 * * *",
			JobTemplate: v1beta1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"linkerd.io/control-plane-component": "heartbeat",
						"linkerd.io/workload-ns":             "linkerd",
					},
					Annotations: map[string]string{
						"linkerd.io/created-by": "linkerd/cli stable-2.8.1",
					},
				},
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							ServiceAccountName: serviceAccountName,
							RestartPolicy:      v1.RestartPolicyNever,
							Containers: []v1.Container{
								{
									Name:            componentName,
									Image:           *r.Config.Spec.Controller.Image,
									ImagePullPolicy: v1.PullIfNotPresent,
									Args: []string{
										"heartbeat",
										fmt.Sprintf("-prometheus-url=http://linkerd-prometheus.%s.svc.%s:9090", r.Config.Namespace, "cluster.local"),
										"-controller-namespace=" + r.Config.Namespace,
										"-log-level=info",
									},
									SecurityContext: &v1.SecurityContext{
										RunAsUser: util.Int64Pointer(2103),
									},
									Resources: v1.ResourceRequirements{
										Limits: v1.ResourceList{
											v1.ResourceCPU:    resource.MustParse("1"),
											v1.ResourceMemory: resource.MustParse("250Mi"),
										},
										Requests: v1.ResourceList{
											v1.ResourceCPU:    resource.MustParse("100m"),
											v1.ResourceMemory: resource.MustParse("50Mi"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
