package linkerd

import (
	"context"
	"fmt"

	linkerdv1alpha1 "github.com/spaghettifunk/linkerd2-operator/pkg/apis/linkerd/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_linkerd")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic. Delete these comments after modifying this file.
**/

// Add creates a new Linkerd Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileLinkerd{client: mgr.GetClient(), scheme: mgr.GetScheme()}
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

// blank assignment to verify that ReconcileLinkerd implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileLinkerd{}

// ReconcileLinkerd reconciles a Linkerd object
type ReconcileLinkerd struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Linkerd object and makes changes based on the state read
// and what is in the Linkerd.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileLinkerd) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Linkerd")

	// Fetch the Linkerd instance
	linkerd := &linkerdv1alpha1.Linkerd{}
	err := r.client.Get(context.TODO(), request.NamespacedName, linkerd)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Linkerd resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Linkerd")
		return reconcile.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: linkerd.Name, Namespace: linkerd.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForLinkerd(linkerd)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// Deployment already exists - don't requeue
	reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	return reconcile.Result{}, nil
}

// deploymentForLinkerd returns a linkerd Deployment object
func (r *ReconcileLinkerd) deploymentForLinkerd(m *linkerdv1alpha1.Linkerd) *appsv1.Deployment {
	ls := labelsForLinkerd(m.Name)
	replicas := m.Spec.Size
	logLevel := m.Spec.LogLevel

	prometheuURL := fmt.Sprintf("http://linkerd-prometheus.%s.svc.cluster.local:9090", m.Namespace)
	dstAddr := fmt.Sprintf("linkerd-dst.%s.svc.cluster.local:8086", m.Namespace)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "gcr.io/linkerd-io/controller:latest",
						Name:  "controller",
						Command: []string{
							"public-api",
							"-prometheus-url=" + prometheuURL,
							"-destination-addr=" + dstAddr,
							"-controller-namespace=" + m.Namespace,
							"-log-level=" + logLevel,
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 11211,
							Name:          "memcached",
						}},
					}},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// labelsForLinkerd returns the labels for selecting the resources
// belonging to the given linkerd CR name.
func labelsForLinkerd(name string) map[string]string {
	return map[string]string{"app": "linkerd", "linkerd_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
