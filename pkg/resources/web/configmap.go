package web

import (
	"encoding/json"

	apiv1 "k8s.io/api/core/v1"

	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) configMap() runtime.Object {
	return &apiv1.ConfigMap{
		ObjectMeta: templates.ObjectMeta(configMapName, nil, r.Config),
		Data: map[string]string{
			"cni_network_config": r.networkConfig(),
		},
	}
}

func (r *Reconciler) networkConfig() string {
	config := map[string]interface{}{
		"cniVersion": "0.3.1",
		"name":       "linekrd-cni",
		"type":       "linekrd-cni",
		// "log_level":  r.Config.Spec.SidecarInjector.InitCNIConfiguration.LogLevel,
		"kubernetes": map[string]interface{}{
			"kubeconfig": "__KUBECONFIG_FILEPATH__",
			// "cni_bin_dir":        r.Config.Spec.SidecarInjector.InitCNIConfiguration.BinDir,
			// "exclude_namespaces": r.Config.Spec.SidecarInjector.InitCNIConfiguration.ExcludeNamespaces,
		},
	}

	marshaledConfig, _ := json.Marshal(config)
	return string(marshaledConfig)
}
