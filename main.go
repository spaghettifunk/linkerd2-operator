/*
Copyright 2020 The Linkerd2 Operator authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	linkerdv1alpha1 "github.com/spaghettifunk/linkerd2-operator/api/v1alpha1"
	"github.com/spaghettifunk/linkerd2-operator/controllers"
	"github.com/spaghettifunk/linkerd2-operator/pkg/k8sutil"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const (
	watchNamespaceEnvVar = "WATCH_NAMESPACE"
	podNamespaceEnvVar   = "POD_NAMESPACE"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(linkerdv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var waitBeforeExitDuration time.Duration
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.DurationVar(&waitBeforeExitDuration, "wait-before-exit-duration", time.Duration(3)*time.Second, "Wait for workers to finish before exiting and removing finalizers")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "d76a60c4.linkerd.io",
		MapperProvider:     k8sutil.NewCachedRESTMapper,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	log.Info("Watch namespaces.")

	namespace, err := getWatchNamespace()
	if err != nil {
		log.Error(err, "Failed to get watch namespace")
		os.Exit(1)
	}
	if namespace != "" {
		log.Info("watch namespace", "namespace", namespace)
	} else {
		log.Info("watch all namespaces")
	}

	log.Info("Registering Components.")

	reconciler := &controllers.ReconcileLinkerd{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Linkerd"),
		Scheme: mgr.GetScheme(),
	}

	if err = reconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Linkerd")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}

	// Cleanup
	// log.Info("removing finalizer from Linkerd resources")
	// err = linkerd.RemoveFinalizers(mgr.GetClient())
	// if err != nil {
	// 	log.Error(err, "could not remove finalizers from Linkerd resources")
	// }
}

func getWatchNamespace() (string, error) {
	podNamespace, found := os.LookupEnv(podNamespaceEnvVar)
	if !found {
		return "", errors.Errorf("%s env variable must be specified and cannot be empty", podNamespaceEnvVar)
	}
	watchNamespace, found := os.LookupEnv(watchNamespaceEnvVar)
	if found {
		if watchNamespace != "" && watchNamespace != podNamespace {
			return "", errors.New("watch namespace must be either empty or equal to pod namespace")
		}

	}
	return watchNamespace, nil
}
