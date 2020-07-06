package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spaghettifunk/linkerd2-operator/pkg/apis"
	"github.com/spaghettifunk/linkerd2-operator/pkg/controller"
	"github.com/spaghettifunk/linkerd2-operator/pkg/k8sutil"

	_ "k8s.io/code-generator/cmd/client-gen/generators"

	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const watchNamespaceEnvVar = "WATCH_NAMESPACE"
const podNamespaceEnvVar = "POD_NAMESPACE"

func main() {
	var developmentMode bool
	flag.BoolVar(&developmentMode, "devel-mode", false, "Set development mode (mainly for logging)")
	var shutdownWaitDuration time.Duration
	flag.DurationVar(&shutdownWaitDuration, "shutdown-wait-duration", time.Duration(30)*time.Second, "Wait duration before shutting down")
	var waitBeforeExitDuration time.Duration
	flag.DurationVar(&waitBeforeExitDuration, "wait-before-exit-duration", time.Duration(3)*time.Second, "Wait for workers to finish before exiting and removing finalizers")
	flag.Parse()
	logf.SetLogger(zap.Logger())
	log := logf.Log.WithName("entrypoint")

	// Get a config to talk to the apiserver
	log.Info("setting up client for manager")
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "unable to set up client config")
		os.Exit(1)
	}

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

	// Create a new manager to provide shared dependencies and start components
	log.Info("setting up manager")
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:      namespace,
		MapperProvider: k8sutil.NewCachedRESTMapper,
	})
	if err != nil {
		log.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "unable add APIs to scheme")
		os.Exit(1)
	}

	stop := setupSignalHandler(mgr, log, shutdownWaitDuration)

	// Setup all Controllers
	log.Info("setting up controller")
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup all webhooks
	// log.Info("setting up webhooks")
	// if err := webhook.AddToManager(mgr); err != nil {
	// 	log.Error(err, "")
	// 	os.Exit(1)
	// }

	log.Info("Starting the Cmd.")
	if err := mgr.Start(stop); err != nil {
		log.Error(err, "unable to run the manager")
		os.Exit(1)
	}

	// Wait a bit for the workers to stop
	time.Sleep(waitBeforeExitDuration)

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

func setupSignalHandler(mgr manager.Manager, log logr.Logger, shutdownWaitDuration time.Duration) (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("termination signal arrived, shutdown gracefully")
		// wait a bit for deletion requests to arrive
		log.Info("wait a bit for CR deletion events to arrive", "waitSeconds", shutdownWaitDuration)
		time.Sleep(shutdownWaitDuration)
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
