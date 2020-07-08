package heartbeat

import (
	"fmt"

	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"
)

func (r *Reconciler) cronjob() runtime.Object {
	return &batchv1.CronJob{
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
		Spec: batchv1.CronJobSpec{
			Schedule: "16 8 * * *",
			JobTemplate: batchv1.JobTemplateSpec{
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
					Template: core.PodTemplateSpec{
						Spec: core.PodSpec{
							ServiceAccountName: serviceAccountName,
							RestartPolicy:      core.RestartPolicyNever,
							Containers: []core.Container{
								{
									Name:            componentName,
									Image:           *r.Config.Spec.Controller.Image,
									ImagePullPolicy: core.PullIfNotPresent,
									Args: []string{
										"heartbeat",
										fmt.Sprintf("-prometheus-url=http://linkerd-prometheus.%s.svc.%s:9090", r.Config.Namespace, "cluster.local"),
										"-controller-namespace=" + r.Config.Namespace,
										"-log-level=info",
									},
									SecurityContext: &core.SecurityContext{
										RunAsUser: util.Int64Pointer(2103),
									},
									Resources: core.ResourceRequirements{
										Limits: core.ResourceList{
											core.ResourceCPU:    resource.MustParse("1"),
											core.ResourceMemory: resource.MustParse("250Mi"),
										},
										Requests: core.ResourceList{
											core.ResourceCPU:    resource.MustParse("100m"),
											core.ResourceMemory: resource.MustParse("50Mi"),
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
