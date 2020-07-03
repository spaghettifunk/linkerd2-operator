package prometheus

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/resource"
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
					},
				},
			},
		},
	}
}

func (r *Reconciler) containers() []apiv1.Container {
	containers := []apiv1.Container{
		{
			Name:            "prometheus",
			Image:           "prom/prometheus:v2.15.2",
			ImagePullPolicy: r.Config.Spec.ImagePullPolicy,
			Args: []string{
				"--storage.tsdb.path=/data",
				"--storage.tsdb.retention.time=6h",
				"--config.file=/etc/prometheus/prometheus.yml",
				"--log.level=info",
			},
			LivenessProbe: &apiv1.Probe{
				InitialDelaySeconds: int32(30),
				TimeoutSeconds:      int32(30),
				Handler: apiv1.Handler{
					HTTPGet: &apiv1.HTTPGetAction{
						Path: "/-/healthy",
						Port: intstr.FromString("9090"),
					},
				},
			},
			ReadinessProbe: &apiv1.Probe{
				InitialDelaySeconds: int32(30),
				TimeoutSeconds:      int32(30),
				Handler: apiv1.Handler{
					HTTPGet: &apiv1.HTTPGetAction{
						Path: "/-/ready",
						Port: intstr.FromString("9090"),
					},
				},
			},
			Resources: apiv1.ResourceRequirements{
				Limits: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("4"),
					apiv1.ResourceMemory: resource.MustParse("8Gi"),
				},
				Requests: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("300m"),
					apiv1.ResourceMemory: resource.MustParse("300Mi"),
				},
			},
			Ports: []apiv1.ContainerPort{
				{
					Name:          "admin-http",
					ContainerPort: int32(9090),
				},
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
		// TODO: export this container among the defaults
		{
			Name:            "linkerd-proxy",
			Image:           "gcr.io/linkerd-io/proxy:stable-2.8.1",
			ImagePullPolicy: apiv1.PullIfNotPresent,
			Resources: apiv1.ResourceRequirements{
				Limits: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("1"),
					apiv1.ResourceMemory: resource.MustParse("250Mi"),
				},
				Requests: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("100m"),
					apiv1.ResourceMemory: resource.MustParse("20Mi"),
				},
			},
			VolumeMounts: []apiv1.VolumeMount{
				{
					Name:      "linkerd-identity-end-entity",
					MountPath: "/var/run/linkerd/identity/end-entity",
				},
			},
			LivenessProbe: &apiv1.Probe{
				Handler: apiv1.Handler{
					HTTPGet: &apiv1.HTTPGetAction{
						Path: "/live",
						Port: intstr.FromString("4191"),
					},
				},
				InitialDelaySeconds: int32(10),
			},
			ReadinessProbe: &apiv1.Probe{
				Handler: apiv1.Handler{
					HTTPGet: &apiv1.HTTPGetAction{
						Path: "/ready",
						Port: intstr.FromString("4191"),
					},
				},
				InitialDelaySeconds: int32(2),
			},
			SecurityContext: &apiv1.SecurityContext{
				RunAsUser:              util.Int64Pointer(2102),
				ReadOnlyRootFilesystem: util.BoolPointer(true),
			},
			Env: []apiv1.EnvVar{
				{
					Name:  "LINKERD2_PROXY_LOG",
					Value: "warn,linkerd=info",
				},
				{
					Name:  "LINKERD2_PROXY_DESTINATION_SVC_ADDR",
					Value: "linkerd-dst.linkerd.svc.cluster.local:8086",
				},
				{
					Name:  "LINKERD2_PROXY_DESTINATION_GET_NETWORKS",
					Value: "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16",
				},
				{
					Name:  "LINKERD2_PROXY_CONTROL_LISTEN_ADDR",
					Value: "0.0.0.0:4190",
				},
				{
					Name:  "LINKERD2_PROXY_ADMIN_LISTEN_ADDR",
					Value: "0.0.0.0:4191",
				},
				{
					Name:  "LINKERD2_PROXY_OUTBOUND_LISTEN_ADDR",
					Value: "127.0.0.1:4140",
				},
				{
					Name:  "LINKERD2_PROXY_INBOUND_LISTEN_ADDR",
					Value: "0.0.0.0:4143",
				},
				{
					Name:  "LINKERD2_PROXY_DESTINATION_GET_SUFFIXES",
					Value: "svc.cluster.local.",
				},
				{
					Name:  "LINKERD2_PROXY_DESTINATION_PROFILE_SUFFIXES",
					Value: "svc.cluster.local.",
				},
				{
					Name:  "LINKERD2_PROXY_INBOUND_ACCEPT_KEEPALIVE",
					Value: "10000ms",
				},
				{
					Name:  "LINKERD2_PROXY_OUTBOUND_CONNECT_KEEPALIVE",
					Value: "10000ms",
				},
				{
					Name: "_pod_ns",
					ValueFrom: &apiv1.EnvVarSource{
						FieldRef: &apiv1.ObjectFieldSelector{
							FieldPath: "metadata.namespace",
						},
					},
				},
				{
					Name:  "LINKERD2_PROXY_DESTINATION_CONTEXT",
					Value: "ns:$(_pod_ns)",
				},
				{
					Name:  "LINKERD2_PROXY_IDENTITY_DIR",
					Value: "/var/run/linkerd/identity/end-entity",
				},
				{
					// TODO: pass the correct .cert file
					Name:  "LINKERD2_PROXY_IDENTITY_TRUST_ANCHORS",
					Value: "",
				},
				{
					Name:  "LINKERD2_PROXY_IDENTITY_TOKEN_FILE",
					Value: "/var/run/secrets/kubernetes.io/serviceaccount/token",
				},
				{
					Name:  "LINKERD2_PROXY_IDENTITY_SVC_ADDR",
					Value: "localhost.:8080",
				},
				{
					Name: "_pod_sa",
					ValueFrom: &apiv1.EnvVarSource{
						FieldRef: &apiv1.ObjectFieldSelector{
							FieldPath: "spec.serviceAccountName",
						},
					},
				},
				{
					Name:  "_l5d_ns",
					Value: "linkerd",
				},
				{
					Name:  "_l5d_trustdomain",
					Value: "cluster.local",
				},
				{
					Name:  "LINKERD2_PROXY_IDENTITY_LOCAL_NAME",
					Value: "$(_pod_sa).$(_pod_ns).serviceaccount.identity.$(_l5d_ns).$(_l5d_trustdomain)",
				},
				{
					Name:  "LINKERD2_PROXY_IDENTITY_SVC_NAME",
					Value: "linkerd-identity.$(_l5d_ns).serviceaccount.identity.$(_l5d_ns).$(_l5d_trustdomain)",
				},
				{
					Name:  "LINKERD2_PROXY_DESTINATION_SVC_NAME",
					Value: "linkerd-destination.$(_l5d_ns).serviceaccount.identity.$(_l5d_ns).$(_l5d_trustdomain)",
				},
				{
					Name:  "LINKERD2_PROXY_TAP_SVC_NAME",
					Value: "linkerd-tap.$(_l5d_ns).serviceaccount.identity.$(_l5d_ns).$(_l5d_trustdomain)",
				},
			},
		},
	}

	return containers
}