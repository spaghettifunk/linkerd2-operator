package templates

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/spaghettifunk/linkerd2-operator/api/v1alpha1"
	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
)

// DefaultAnnotations are the default annotations for deployments
func DefaultAnnotations(version string) map[string]string {
	return map[string]string{
		"linkerd.io/created-by": version,
	}
}

// GetResourcesRequirementsOrDefault sets the new resources constraints or use the defaults
func GetResourcesRequirementsOrDefault(requirements *apiv1.ResourceRequirements, defaults *apiv1.ResourceRequirements) apiv1.ResourceRequirements {
	if requirements != nil {
		return *requirements
	}
	return *defaults
}

// DefaultRollingUpdateStrategy defines the default rolling update strategy
func DefaultRollingUpdateStrategy() appsv1.DeploymentStrategy {
	return appsv1.DeploymentStrategy{
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxSurge:       util.IntstrPointer(1),
			MaxUnavailable: util.IntstrPointer(0),
		},
	}
}

// DefaultReadinessProbe returns the default readiness probe values
func DefaultReadinessProbe(path, port string, failureThreshold, timeoutSeconds int) *apiv1.Probe {
	return &apiv1.Probe{
		TimeoutSeconds:   int32(timeoutSeconds),
		FailureThreshold: int32(failureThreshold),
		Handler: apiv1.Handler{
			HTTPGet: &apiv1.HTTPGetAction{
				Path: path,
				Port: intstr.FromString(port),
			},
		},
	}
}

// DefaultLivenessProbe returns the default liveness probe values
func DefaultLivenessProbe(path, port string, initialDelaySeconds, timeoutSeconds int) *apiv1.Probe {
	return &apiv1.Probe{
		TimeoutSeconds:      int32(timeoutSeconds),
		InitialDelaySeconds: int32(initialDelaySeconds),
		Handler: apiv1.Handler{
			HTTPGet: &apiv1.HTTPGetAction{
				Path: path,
				Port: intstr.FromString(port),
			},
		},
	}
}

// DefaultServicePort returns the default values for the servicePort
func DefaultServicePort(name string, port, targetPort int) apiv1.ServicePort {
	return apiv1.ServicePort{
		Name:       name,
		Port:       int32(port),
		TargetPort: intstr.FromInt(targetPort),
	}
}

// DefaultContainerPort returns the default values for the containerPort
func DefaultContainerPort(name string, containerPort int) apiv1.ContainerPort {
	return apiv1.ContainerPort{
		Name:          name,
		ContainerPort: int32(containerPort),
	}
}

// ProxyInitContainer returns the proxy init container definition
func ProxyInitContainer() []apiv1.Container {
	initContainers := []apiv1.Container{
		{
			Name:            "linkerd-init",
			Image:           "gcr.io/linkerd-io/proxy-init:v1.3.3",
			ImagePullPolicy: apiv1.PullIfNotPresent,
			Args: []string{
				"--incoming-proxy-port",
				"4143",
				"--outgoing-proxy-port",
				"4140",
				"--proxy-uid",
				"2102",
				"--inbound-ports-to-ignore",
				"4190,4191",
				"--outbound-ports-to-ignore",
				"443",
			},
			Resources: apiv1.ResourceRequirements{
				Limits: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("100m"),
					apiv1.ResourceMemory: resource.MustParse("50Mi"),
				},
				Requests: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("10m"),
					apiv1.ResourceMemory: resource.MustParse("10Mi"),
				},
			},
			SecurityContext: &apiv1.SecurityContext{
				AllowPrivilegeEscalation: util.BoolPointer(false),
				Capabilities: &apiv1.Capabilities{
					Add: []apiv1.Capability{
						"NET_ADMIN",
						"NET_RAW",
					},
				},
				Privileged:             util.BoolPointer(false),
				ReadOnlyRootFilesystem: util.BoolPointer(false),
				RunAsNonRoot:           util.BoolPointer(false),
				RunAsUser:              util.Int64Pointer(0),
			},
			TerminationMessagePath:   apiv1.TerminationMessagePathDefault,
			TerminationMessagePolicy: apiv1.TerminationMessageReadFile,
		},
	}
	return initContainers
}

// DefaultProxyContainer returns the Proxy container definition
func DefaultProxyContainer(config v1alpha1.LinkerdSpec) apiv1.Container {
	return apiv1.Container{
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
				apiv1.ResourceMemory: resource.MustParse("50Mi"),
			},
		},
		VolumeMounts: []apiv1.VolumeMount{
			{
				Name:      "linkerd-identity-end-entity",
				MountPath: "/var/run/linkerd/identity/end-entity",
			},
		},
		// LivenessProbe: &apiv1.Probe{
		// 	Handler: apiv1.Handler{
		// 		HTTPGet: &apiv1.HTTPGetAction{
		// 			Path: "/live",
		// 			Port: intstr.FromString("4191"),
		// 		},
		// 	},
		// 	InitialDelaySeconds: int32(10),
		// },
		// ReadinessProbe: &apiv1.Probe{
		// 	Handler: apiv1.Handler{
		// 		HTTPGet: &apiv1.HTTPGetAction{
		// 			Path: "/ready",
		// 			Port: intstr.FromString("4191"),
		// 		},
		// 	},
		// 	InitialDelaySeconds: int32(2),
		// },
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
				Name:  "LINKERD2_PROXY_IDENTITY_TRUST_ANCHORS",
				Value: config.SelfSignedCertificates.TrustAnchorsPEM,
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
	}
}
