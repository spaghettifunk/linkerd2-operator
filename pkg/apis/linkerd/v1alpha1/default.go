package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"

	apiv1 "k8s.io/api/core/v1"
)

const (
	linkerdImageHub       = "docker.io/davideberdin"
	linkerdImageVersion   = "v0.1.0"
	defaultImageHub       = "docker.io/linkerd2"
	defaultImageVersion   = "2.28.0"
	defaultReplicaCount   = 1
	defaultMinReplicas    = 1
	defaultMaxReplicas    = 5
	defaultProxyImage     = defaultImageHub + "/" + "proxyv2" + ":" + defaultImageVersion
	defaultProxyInitImage = defaultImageHub + "/" + "proxyv2" + ":" + defaultImageVersion
)

var defaultResources = &apiv1.ResourceRequirements{
	Requests: apiv1.ResourceList{
		apiv1.ResourceCPU: resource.MustParse("10m"),
	},
}

var defaultProxyResources = &apiv1.ResourceRequirements{
	Requests: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("100m"),
		apiv1.ResourceMemory: resource.MustParse("128Mi"),
	},
	Limits: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("2000m"),
		apiv1.ResourceMemory: resource.MustParse("1024Mi"),
	},
}

// SetDefaults sets the defaults values for all the components
func SetDefaults(config *Linkerd) {
	if config.Spec.Controller.Resources == nil {
		config.Spec.Controller.Resources = defaultResources
	}
}
