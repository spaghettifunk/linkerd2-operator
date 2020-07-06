package web

import (
	"fmt"

	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Reconciler) deployment() runtime.Object {
	labels := util.MergeStringMaps(r.labels(), r.deploymentLabels())
	return &appsv1.Deployment{
		ObjectMeta: templates.ObjectMeta(deploymentName, labels, r.Config),
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: util.IntstrPointer(1),
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: templates.DefaultAnnotations(string(r.Config.Spec.Version)),
				},
				Spec: apiv1.PodSpec{
					DNSPolicy:     apiv1.DNSClusterFirst,
					RestartPolicy: apiv1.RestartPolicyAlways,
					NodeSelector: map[string]string{
						"beta.kubernetes.io/os": "linux",
					},
					Tolerations: []apiv1.Toleration{
						{
							Operator: apiv1.TolerationOpExists,
							Effect:   apiv1.TaintEffectNoSchedule,
						},
						{
							Operator: apiv1.TolerationOpExists,
							Effect:   apiv1.TaintEffectNoExecute,
						},
					},
					TerminationGracePeriodSeconds: util.Int64Pointer(5),
					ServiceAccountName:            serviceAccountName,
					Containers:                    r.container(),
					InitContainers:                templates.ProxyInitContainer(),
					Volumes: []apiv1.Volume{
						{
							Name: "config",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: "linkerd-config",
									},
									DefaultMode: util.IntPointer(420),
								},
							},
						},
						{
							Name: "linkerd-identity-end-entity",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{
									Medium: apiv1.StorageMediumMemory,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *Reconciler) container() []apiv1.Container {

	apiAddr := fmt.Sprintf("-api-addr=linkerd-controller-api.%s.svc.cluster.local:8085", r.Config.Namespace)
	grafanaAddr := fmt.Sprintf("-grafana-addr=linkerd-grafana.%s.svc.cluster.local:3000", r.Config.Namespace)
	controllerNamespace := fmt.Sprintf("-controller-namespace=%s", r.Config.Namespace)
	enforcedHost := fmt.Sprintf("-enforced-host=^(localhost|127\\.0\\.0\\.1|linkerd-web\\.%s\\.svc\\.cluster\\.local|linkerd-web\\.%s\\.svc|\\[::1\\])(:\\d+)?$", r.Config.Namespace, r.Config.Namespace)

	args := []string{
		apiAddr,
		grafanaAddr,
		controllerNamespace,
		enforcedHost,
		"-log-level=info",
	}

	webConfig := r.Config.Spec.Web
	containers := []apiv1.Container{
		{
			Name:            "web",
			Image:           *webConfig.Image,
			ImagePullPolicy: r.Config.Spec.ImagePullPolicy,
			Args:            args,
			LivenessProbe:   templates.DefaultLivenessProbe("/ping", "9994", 10, 1),
			ReadinessProbe:  templates.DefaultReadinessProbe("/ready", "9994", 7, 1),
			Resources:       *webConfig.Resources,
			Ports: []apiv1.ContainerPort{
				templates.DefaultContainerPort("http", 8084),
				templates.DefaultContainerPort("admin-http", 9994),
			},
			VolumeMounts: []apiv1.VolumeMount{
				{
					Name:      "config",
					MountPath: "/var/run/linkerd/config",
				},
			},
			TerminationMessagePath:   apiv1.TerminationMessagePathDefault,
			TerminationMessagePolicy: apiv1.TerminationMessageReadFile,
		},
	}

	return containers
}
