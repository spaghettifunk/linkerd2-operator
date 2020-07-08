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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	"github.com/pkg/errors"
	linkerdv1alpha1 "github.com/spaghettifunk/linkerd2-operator/api/v1alpha1"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources"
	linkerdcontroller "github.com/spaghettifunk/linkerd2-operator/pkg/resources/controller"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/destination"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/heartbeat"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/identity"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/prometheus"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/proxyinjector"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/psp"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/tap"
	"github.com/spaghettifunk/linkerd2-operator/pkg/resources/web"
	"github.com/spaghettifunk/linkerd2-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const finalizerID = "linkerd2-operator.finializer.linkerd.io"
const linkerdSecretTypePrefix = "linkerd.io"

var log = logf.Log.WithName("controller")
var watchCreatedResourcesEvents bool

// Add creates a new Linkerd Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileLinkerd{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("linkerd-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Linkerd
	err = c.Watch(&source.Kind{Type: &linkerdv1alpha1.Linkerd{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Linkerd
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &linkerdv1alpha1.Linkerd{},
	})
	if err != nil {
		return err
	}

	return nil
}

// ReconcileLinkerd reconciles a Linkerd object
type ReconcileLinkerd struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Linkerd object and makes changes based on the state read
// and what is in the Linkerd.Spec
// +kubebuilder:rbac:groups=linkerd.linkerd.io,resources=linkerds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=linkerd.linkerd.io,resources=linkerds/status,verbs=get;update;patch
func (r *ReconcileLinkerd) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	// Fetch the Linkerd instance
	config := &linkerdv1alpha1.Linkerd{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, config)
	if err != nil {
		if k8errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("Linkerd resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get Linkerd")
		return reconcile.Result{}, err
	}

	logger.Info("Reconciling Linkerd")

	if !config.Spec.Version.IsSupported() {
		err = errors.New("intended Linkerd version is unsupported by this version of the operator")
		logger.Error(err, "", "version", config.Spec.Version)
		return reconcile.Result{
			Requeue: false,
		}, nil
	}

	// Set default values where not set
	linkerdv1alpha1.SetDefaults(config)

	// start reconciling loop
	result, err := r.reconcile(logger, config)
	if err != nil {
		updateErr := updateStatus(r.Client, config, linkerdv1alpha1.ReconcileFailed, err.Error(), logger)
		if updateErr != nil {
			logger.Error(updateErr, "failed to update state")
			return result, errors.WithStack(err)
		}
		return result, emperror.Wrap(err, "could not reconcile Linkerd")
	}
	return result, nil
}

func (r *ReconcileLinkerd) reconcile(logger logr.Logger, config *linkerdv1alpha1.Linkerd) (reconcile.Result, error) {
	if config.Status.Status == "" {
		err := updateStatus(r.Client, config, linkerdv1alpha1.Created, "", logger)
		if err != nil {
			return reconcile.Result{}, errors.WithStack(err)
		}
	}

	// for each component do a reconciliation
	reconcilers := []resources.ComponentReconciler{
		linkerdcontroller.New(r.Client, config),
		destination.New(r.Client, config),
		heartbeat.New(r.Client, config),
		identity.New(r.Client, config),
		prometheus.New(r.Client, config),
		proxyinjector.New(r.Client, config),
		// serviceprofile.New(r.Client, config),
		// trafficsplit.New(r.Client, config),
		web.New(r.Client, config),
		tap.New(r.Client, config),
		psp.New(r.Client, config),
	}
	for _, rec := range reconcilers {
		err := rec.Reconcile(logger)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	err := updateStatus(r.Client, config, linkerdv1alpha1.Available, "", logger)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	logger.Info("reconcile finished")

	return reconcile.Result{}, nil
}

func updateStatus(c client.Client, config *linkerdv1alpha1.Linkerd, status linkerdv1alpha1.ConfigState, errorMessage string, logger logr.Logger) error {

	typeMeta := config.TypeMeta
	config.Status.Status = status
	config.Status.ErrorMessage = errorMessage

	err := c.Status().Update(context.Background(), config)
	if k8errors.IsNotFound(err) {
		err = c.Update(context.Background(), config)
	}

	if err != nil {
		if !k8errors.IsConflict(err) {
			return emperror.Wrapf(err, "could not update Linkerd state to '%s'", status)
		}

		var actualConfig linkerdv1alpha1.Linkerd
		err := c.Get(context.TODO(), types.NamespacedName{
			Namespace: config.Namespace,
			Name:      config.Name,
		}, &actualConfig)

		if err != nil {
			return emperror.Wrap(err, "could not get config for updating status")
		}

		actualConfig.Status.Status = status
		actualConfig.Status.ErrorMessage = errorMessage

		err = c.Status().Update(context.Background(), &actualConfig)
		if k8errors.IsNotFound(err) {
			err = c.Update(context.Background(), &actualConfig)
		}
		if err != nil {
			return emperror.Wrapf(err, "could not update Linkerd state to '%s'", status)
		}
	}

	// update loses the typeMeta of the config that's used later when setting ownerrefs
	config.TypeMeta = typeMeta
	logger.Info("Linkerd state updated", "status", status)

	return nil
}

// RemoveFinalizers removes the finalizers from the context
func RemoveFinalizers(c client.Client) error {
	var linkerds linkerdv1alpha1.LinkerdList

	// fix this!
	// err := c.List(context.TODO(), &client.ListOptions{}, &linkerds)
	// if err != nil {
	// 	return emperror.Wrap(err, "could not list Linkerd resources")
	// }

	for _, linkerd := range linkerds.Items {
		linkerd.ObjectMeta.Finalizers = util.RemoveString(linkerd.ObjectMeta.Finalizers, finalizerID)
		if err := c.Update(context.Background(), &linkerd); err != nil {
			return emperror.WrapWith(err, "could not remove finalizer from Linkerd resource", "name", linkerd.GetName())
		}
		if err := updateStatus(c, &linkerd, linkerdv1alpha1.Unmanaged, "", log); err != nil {
			return emperror.Wrap(err, "could not update status of Linkerd resource")
		}
	}

	return nil
}

// SetupWithManager sets the reconciler with the manager
func (r *ReconcileLinkerd) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&linkerdv1alpha1.Linkerd{}).
		Complete(r)
}
