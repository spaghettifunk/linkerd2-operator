package templates

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

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
