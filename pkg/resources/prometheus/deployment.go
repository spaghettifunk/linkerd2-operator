package prometheus

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
							Name: "data",
						},
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
						{
							Name: "prometheus-config",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: configmapName,
									},
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

func (r *Reconciler) containers() []apiv1.Container {
	prometheusConfig := r.Config.Spec.Prometheus
	containers := []apiv1.Container{
		templates.DefaultProxyContainer(r.Config.Spec),
		{
			Name:            "prometheus",
			Image:           *prometheusConfig.Image,
			ImagePullPolicy: r.Config.Spec.ImagePullPolicy,
			Args: []string{
				"--storage.tsdb.path=/data",
				"--storage.tsdb.retention.time=6h",
				"--config.file=/etc/prometheus/prometheus.yml",
				"--log.level=info",
			},
			// LivenessProbe:  templates.DefaultLivenessProbe("/-/healthy", "9090", 30, 30),
			// ReadinessProbe: templates.DefaultReadinessProbe("/-/ready", "9090", 30, 30),
			Resources: *prometheusConfig.Resources,
			Ports: []apiv1.ContainerPort{
				templates.DefaultContainerPort("admin-http", 9090),
			},
			SecurityContext: &apiv1.SecurityContext{
				RunAsUser: util.Int64Pointer(65534),
			},
			VolumeMounts: []apiv1.VolumeMount{
				{
					MountPath: "/data",
					Name:      "data",
				},
				{
					Name:      "prometheus-config",
					ReadOnly:  true,
					MountPath: "/etc/prometheus/prometheus.yml",
					SubPath:   "prometheus.yml",
				},
			},
			TerminationMessagePath:   apiv1.TerminationMessagePathDefault,
			TerminationMessagePolicy: apiv1.TerminationMessageReadFile,
		},
	}

	return containers
}
