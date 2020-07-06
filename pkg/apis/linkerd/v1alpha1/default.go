package v1alpha1

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/resource"

	apiv1 "k8s.io/api/core/v1"
)

const (
	linkerdImageHub              = "docker.io/davideberdin"
	linkerdImageVersion          = "v0.1.0"
	defaultImageHub              = "gcr.io/linkerd-io/"
	defaultImageVersion          = "2.8.1"
	defaultProxyInitVersion      = "v1.3.3"
	defaultCollectorImageHub     = "omnition"
	defaultCollectorImageVersion = "0.1.11"
	defaultJaegerImageHub        = "jaegertracing"
	defaultJaegerImageVersion    = "1.17.1"
	defaultImagePullPolicy       = "IfNotPresent"
	defaultNetworkName           = "cluster.local"
	defaultLogLevel              = "log-level:info"
	// replicas
	defaultReplicaCount = 1
	defaultMinReplicas  = 1
	defaultMaxReplicas  = 5
	// images
	defaultControllerImage = defaultImageHub + "/" + "controller" + ":" + defaultImageVersion
	defaultProxyImage      = defaultImageHub + "/" + "proxy" + ":" + defaultImageVersion
	defaultProxyInitImage  = defaultImageHub + "/" + "proxy-init" + ":" + defaultImageVersion
	defaultWebImage        = defaultImageHub + "/" + "web" + ":" + defaultImageVersion
	defaultGrafanaImage    = defaultImageHub + "/" + "grafana" + ":" + defaultImageVersion
	defaultCollectorImage  = defaultCollectorImageHub + "/" + "opencensus-collector" + ":" + defaultCollectorImageVersion
	defaultJaegerImage     = defaultJaegerImageHub + "/" + "all-in-one" + ":" + defaultJaegerImageVersion
	// resources
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

// var defaultControllerServicePorts = []ServicePort{
// 	{ServicePort: corev1.ServicePort{Name: "http", Port: int32(8085), TargetPort: intstr.FromString("8085")}},
// }

// SetDefaults sets the defaults values for all the components
func SetDefaults(config *Linkerd) {
	// controller
	if config.Spec.Controller.Image == nil {
		config.Spec.Controller.Image = util.StrPointer(defaultControllerImage)
	}
	if config.Spec.Controller.Resources == nil {
		config.Spec.Controller.Resources = defaultResources
	}
	// if config.Spec.Controller

	// destination

	// heartbeat

	// identity

	// prometheus
	if config.Spec.Prometheus.Resources == nil {
		config.Spec.Prometheus.Resources = &apiv1.ResourceRequirements{
			Limits: apiv1.ResourceList{
				apiv1.ResourceCPU:    resource.MustParse("4"),
				apiv1.ResourceMemory: resource.MustParse("8Gi"),
			},
			Requests: apiv1.ResourceList{
				apiv1.ResourceCPU:    resource.MustParse("300m"),
				apiv1.ResourceMemory: resource.MustParse("300Mi"),
			},
		}
	}

	// proxyinjector
	if config.Spec.ProxyInjector.Resources == nil {
		config.Spec.ProxyInjector.Resources = &apiv1.ResourceRequirements{
			Limits: apiv1.ResourceList{
				apiv1.ResourceCPU:    resource.MustParse("1"),
				apiv1.ResourceMemory: resource.MustParse("250Mi"),
			},
			Requests: apiv1.ResourceList{
				apiv1.ResourceCPU:    resource.MustParse("100m"),
				apiv1.ResourceMemory: resource.MustParse("50Mi"),
			},
		}
	}

	// psp

	// serviceprofile

	// tap

	// trafficsplit

	// web
}
