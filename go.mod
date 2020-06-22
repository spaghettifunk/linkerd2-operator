module github.com/spaghettifunk/linkerd2-operator

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/goph/emperror v0.17.2
	github.com/operator-framework/operator-sdk v0.18.1
	github.com/pkg/errors v0.9.1
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.2 // Required by prometheus-operator
)
