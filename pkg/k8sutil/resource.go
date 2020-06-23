package k8sutil

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Reconcile(log logr.Logger, client runtimeClient.Client, desired runtime.Object, desiredState DesiredState) error {
	return nil
}
