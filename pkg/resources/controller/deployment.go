package controller

import (
	"fmt"

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
		"-public-api",
		fmt.Sprintf("-prometheus-url=http://linkerd-prometheus.%s.svc.%s:9090", r.Config.Namespace, "cluster.local"),
		fmt.Sprintf("-destination-addr=linkerd-dst.%s.svc.%s:8086", r.Config.Namespace, "cluster.local"),
		"-controller-namespace=" + r.Config.Namespace,
		"-log-level=info",
	}

	controllerConfig := r.Config.Spec.Controller
	containers := []apiv1.Container{
		{
			Name:            "public-api",
			Image:           *controllerConfig.Image,
			ImagePullPolicy: r.Config.Spec.ImagePullPolicy,
			Args:            args,
			LivenessProbe: &apiv1.Probe{
				InitialDelaySeconds: int32(10),
				Handler: apiv1.Handler{
					HTTPGet: &apiv1.HTTPGetAction{
						Path: "/ping",
						Port: intstr.FromString("9995"),
					},
				},
			},
			ReadinessProbe: &apiv1.Probe{
				FailureThreshold: int32(7),
				Handler: apiv1.Handler{
					HTTPGet: &apiv1.HTTPGetAction{
						Path: "/ready",
						Port: intstr.FromString("9995"),
					},
				},
			},
			Resources: templates.GetResourcesRequirementsOrDefault(nil, nil),
			Ports: []apiv1.ContainerPort{
				{
					Name:          "http",
					ContainerPort: int32(8085),
				},
				{
					Name:          "admin-http",
					ContainerPort: int32(9995),
				},
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
