package destination

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Reconciler) deployment() runtime.Object {
	labels := util.MergeStringMaps(r.labels(), r.deploymentLabels())
	return &appsv1.Deployment{
		ObjectMeta: templates.ObjectMetaWithAnnotations(
			deploymentName,
			util.MergeMultipleStringMaps(r.deploymentLabels(), r.labels()),
			templates.DefaultAnnotations(string(r.Config.Spec.Version)),
			r.Config,
		),
		Spec: appsv1.DeploymentSpec{
			Strategy: appsv1.DeploymentStrategy{
				// TODO: enable only when podAntiAffinity is true
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{IntVal: 1},
				},
			},
			Replicas: r.Config.Spec.Controller.ReplicaCount,
			Selector: &v1.LabelSelector{
				MatchLabels: r.labels(),
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: templates.DefaultAnnotations(string(r.Config.Spec.Version)),
				},
				Spec: apiv1.PodSpec{
					ServiceAccountName: serviceAccountName,
					Containers:         r.containers(),
					InitContainers:     templates.ProxyInitContainer(),
					Volumes: []apiv1.Volume{
						{
							Name: "config",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: "linkerd-config",
									},
								},
							},
						},
						// TODO: add Identity volume here
						// TODO: add Tracing labels here
					},
				},
			},
		},
	}
}

func (r *Reconciler) containers() []apiv1.Container {

	args := []string{
		"destination",
		"-addr=:8086",
		"-controller-namespace=" + r.Config.Namespace,
		"-enable-h2-upgrade=true",
		"-log-level=info",
	}

	destinationConfig := r.Config.Spec.Destination
	containers := []apiv1.Container{
		{
			Name:            "destination",
			Image:           *destinationConfig.Image,
			ImagePullPolicy: r.Config.Spec.ImagePullPolicy,
			Args:            args,
			LivenessProbe:   templates.DefaultLivenessProbe("/ping", "9996", 10, 30),
			ReadinessProbe:  templates.DefaultReadinessProbe("/ready", "9996", 7, 30),
			Resources:       *destinationConfig.Resources,
			Ports: []apiv1.ContainerPort{
				templates.DefaultContainerPort("grpc", 8086),
				templates.DefaultContainerPort("admin-http", 9996),
			},
			SecurityContext: &apiv1.SecurityContext{
				RunAsUser: util.Int64Pointer(2103),
			},
			VolumeMounts: []apiv1.VolumeMount{
				{
					MountPath: "/var/run/linkerd/config",
					Name:      "config",
				},
			},
			TerminationMessagePath:   apiv1.TerminationMessagePathDefault,
			TerminationMessagePolicy: apiv1.TerminationMessageReadFile,
		},
	}

	return containers
}
