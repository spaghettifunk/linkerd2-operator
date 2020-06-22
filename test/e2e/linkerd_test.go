package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	"github.com/spaghettifunk/linkerd2-operator/pkg/apis"
	operator "github.com/spaghettifunk/linkerd2-operator/pkg/apis/linkerd/v1alpha1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestLinkerd(t *testing.T) {
	linkerdList := &operator.LinkerdList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, linkerdList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("linekrd-group", func(t *testing.T) {
		t.Run("Cluster", LinkerdHighAvailability)
		t.Run("Cluster2", LinkerdHighAvailability)
	})
}

func linkerdScaleTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create linkerd custom resource
	exampleLinkerd := &operator.Linkerd{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-linkerd",
			Namespace: namespace,
		},
		Spec: operator.LinkerdSpec{
			Size: 3,
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), exampleLinkerd, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	// wait for example-linkerd to reach 3 replicas
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-linkerd", 3, retryInterval, timeout)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-linkerd", Namespace: namespace}, exampleLinkerd)
	if err != nil {
		return err
	}
	exampleLinkerd.Spec.Size = 4
	err = f.Client.Update(goctx.TODO(), exampleLinkerd)
	if err != nil {
		return err
	}

	// wait for example-linkerd to reach 4 replicas
	return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-linkerd", 4, retryInterval, timeout)
}

func LinkerdHighAvailability(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global
	// wait for linkerd2-operator to be ready
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "linkerd2-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = linkerdScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
