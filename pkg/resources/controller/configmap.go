package controller

import (
	"github.com/hoisie/mustache"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/templates"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var globalCfg = `{
    "linkerdNamespace": "linkerd",
    "cniEnabled": false,
    "version": "{{version}}",
    "identityContext": {
        "trustDomain": "{{trustDomain}}",
        "trustAnchorsPem": "",
        "issuanceLifetime": "86400s",
        "clockSkewAllowance": "20s",
        "scheme": "linkerd.io/tls"
    },
    "autoInjectContext": null,
    "omitWebhookSideEffects": false,
    "clusterDomain": "{{clusterDomain}}"
}`

var proxyCfg = `
{
    "proxyImage": {
        "imageName": "gcr.io/linkerd-io/proxy",
        "pullPolicy": "IfNotPresent"
    },
    "proxyInitImage": {
        "imageName": "gcr.io/linkerd-io/proxy-init",
        "pullPolicy": "IfNotPresent"
    },
    "controlPort": {
        "port": 4190
    },
    "ignoreInboundPorts": [],
    "ignoreOutboundPorts": [],
    "inboundPort": {
        "port": 4143
    },
    "adminPort": {
        "port": 4191
    },
    "outboundPort": {
        "port": 4140
    },
    "resource": {
        "requestCpu": "100m",
        "requestMemory": "20Mi",
        "limitCpu": "1",
        "limitMemory": "250Mi"
    },
    "proxyUid": "2102",
    "logLevel": {
        "level": "warn,linkerd=info"
    },
    "disableExternalProfiles": true,
    "proxyVersion": "{{version}}",
    "proxyInitImageVersion": "{{proxyInitImageVersion}}",
    "debugImage": {
        "imageName": "gcr.io/linkerd-io/debug",
        "pullPolicy": "IfNotPresent"
    },
    "debugImageVersion": "{{version}}",
    "destinationGetNetworks": "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"
}`

var installCfg = `
{
    "cliVersion": "{{version}}",
    "flags": [
        {
            "name": "ha",
            "value": "{{isHA}}"
        }
    ]
}
`

func (r *Reconciler) configmap() runtime.Object {
	return &apiv1.ConfigMap{
		ObjectMeta: templates.ObjectMetaWithAnnotations(configmapName, r.labels(), templates.DefaultAnnotations(string(r.Config.Spec.Version)), r.Config),
		Data: map[string]string{
			"global": mustache.Render(globalCfg, map[string]string{
				"version":       "stable-2.8.1",
				"trustDomain":   "cluster.local",
				"clusterDomain": "cluster.local",
			}),
			"proxy": mustache.Render(proxyCfg, map[string]string{
				"version":               "stable-2.8.1",
				"proxyInitImageVersion": "v1.3.3",
			}),
			"install": mustache.Render(installCfg, map[string]string{
				"version": "stable-2.8.1",
				"isHA":    "false",
			}),
		},
	}
}
