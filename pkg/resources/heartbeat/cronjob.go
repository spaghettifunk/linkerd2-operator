package heartbeat

import (
	"fmt"

	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"
)

func (r *Reconciler) cronjob() runtime.Object {
	labels := util.MergeStringMaps(r.labels(), r.deploymentLabels())
	return &batchv1.CronJob{
		ObjectMeta: templates.ObjectMetaWithAnnotations(
			cronjobName,
			util.MergeMultipleStringMaps(r.deploymentLabels(), r.labels()),
			templates.DefaultAnnotations(string(r.Config.Spec.Version)),
			r.Config,
		),
		Spec: batchv1.CronJobSpec{
			Schedule:                   "",
			SuccessfulJobsHistoryLimit: util.IntPointer(0),
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: templates.DefaultAnnotations(string(r.Config.Spec.Version)),
				},
				Spec: batchv1.JobSpec{
					Template: core.PodTemplateSpec{
						Spec: core.PodSpec{
							ServiceAccountName: serviceAccountName,
							RestartPolicy:      core.RestartPolicyNever,
							Containers: []core.Container{
								{
									Name:            componentName,
									Image:           *r.Config.Spec.Heartbeat.Image,
									ImagePullPolicy: core.PullAlways,
									Args: []string{
										"heartbeat",
										fmt.Sprintf("-prometheus-url=http://linkerd-prometheus.%s.svc.%s:9090", r.Config.Namespace, "cluster.local"),
										"-controller-namespace=" + r.Config.Namespace,
										"-log-level=info",
									},
									SecurityContext: &core.SecurityContext{
										RunAsUser: util.Int64Pointer(2103),
									},
									// TODO: add default resources here
									Resources: core.ResourceRequirements{},
								},
							},
						},
					},
				},
			},
		},
	}
}
