package identity

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

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
				// RollingUpdate: &appsv1.RollingUpdateDeployment{
				// 	MaxUnavailable: &intstr.IntOrString{IntVal: 1},
				// },
			},
			Replicas: r.Config.Spec.Identity.ReplicaCount,
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
						{
							Name: "identity-issuer",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: secretName,
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
					NodeSelector: map[string]string{
						"beta.kubernetes.io/os": "linux",
					},
					// Affinity: &apiv1.Affinity{
					// 	PodAntiAffinity: &apiv1.PodAntiAffinity{
					// 		RequiredDuringSchedulingIgnoredDuringExecution: []apiv1.PodAffinityTerm{
					// 			{
					// 				LabelSelector: &metav1.LabelSelector{
					// 					MatchExpressions: []metav1.LabelSelectorRequirement{
					// 						{
					// 							Key:      "linkerd.io/control-plane-component",
					// 							Operator: "In",
					// 							Values:   []string{"identity"},
					// 						},
					// 					},
					// 				},
					// 				TopologyKey: "kubernetes.io/hostname",
					// 			},
					// 		},
					// 		PreferredDuringSchedulingIgnoredDuringExecution: []apiv1.WeightedPodAffinityTerm{
					// 			{
					// 				Weight: 100,
					// 				PodAffinityTerm: apiv1.PodAffinityTerm{
					// 					LabelSelector: &metav1.LabelSelector{
					// 						MatchExpressions: []metav1.LabelSelectorRequirement{
					// 							{
					// 								Key:      "linkerd.io/control-plane-component",
					// 								Operator: "In",
					// 								Values:   []string{"identity"},
					// 							},
					// 						},
					// 					},
					// 					TopologyKey: "failure-domain.beta.kubernetes.io/zone",
					// 				},
					// 			},
					// 		},
					// 	},
					// },
				},
			},
		},
	}
}

func (r *Reconciler) containers() []apiv1.Container {
	identityConfig := r.Config.Spec.Identity
	containers := []apiv1.Container{
		templates.DefaultProxyContainer(r.Config.Spec),
		{
			Name:            "identity",
			Image:           *identityConfig.Image,
			ImagePullPolicy: r.Config.Spec.ImagePullPolicy,
			Args: []string{
				"identity",
				"-log-level=info",
			},
			LivenessProbe:  templates.DefaultLivenessProbe("/ping", 9990, 10, 30),
			ReadinessProbe: templates.DefaultReadinessProbe("/ready", 9990, 7, 30),
			Resources:      *identityConfig.Resources,
			Ports: []apiv1.ContainerPort{
				templates.DefaultContainerPort("grpc", 8080),
				templates.DefaultContainerPort("admin-http", 9990),
			},
			SecurityContext: &apiv1.SecurityContext{
				RunAsUser: util.Int64Pointer(2103),
			},
			VolumeMounts: []apiv1.VolumeMount{
				{
					MountPath: "/var/run/linkerd/config",
					Name:      "config",
				},
				{
					MountPath: "/var/run/linkerd/identity/issuer",
					Name:      "identity-issuer",
				},
			},
			TerminationMessagePath:   apiv1.TerminationMessagePathDefault,
			TerminationMessagePolicy: apiv1.TerminationMessageReadFile,
		},
	}
	return containers
}
