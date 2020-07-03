module github.com/spaghettifunk/linkerd2-operator

go 1.13

require (
	github.com/go-logr/logr v0.2.0
	github.com/goph/emperror v0.17.2
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/operator-framework/operator-sdk v0.18.1
	github.com/pkg/errors v0.9.1
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools v0.0.0-20200702044944-0cc1aa72b347 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	k8s.io/api v0.18.4
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/code-generator v0.18.4
	k8s.io/gengo v0.0.0-20200630090205-15d76db0a9e6 // indirect
	k8s.io/klog/v2 v2.3.0 // indirect
	k8s.io/kube-aggregator v0.18.4
	k8s.io/kube-openapi v0.0.0-20200615155156-dffdd1682719 // indirect
	k8s.io/kubernetes v1.13.0
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.2 // Required by prometheus-operator
)
