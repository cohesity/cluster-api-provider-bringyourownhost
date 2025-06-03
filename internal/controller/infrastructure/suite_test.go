// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package infrastructure_test

import (
	"context"
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"

	infrastructurev1beta1 "github.com/cohesity/cluster-api-provider-bringyourownhost/api/infrastructure/v1beta1"
	controllers "github.com/cohesity/cluster-api-provider-bringyourownhost/internal/controller/infrastructure"
	"github.com/cohesity/cluster-api-provider-bringyourownhost/test/builder"

	// +kubebuilder:scaffold:imports

	fakeclientset "k8s.io/client-go/kubernetes/fake"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	ctrl "sigs.k8s.io/controller-runtime"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	testEnv                               *envtest.Environment
	k8sClient                             client.Client
	clientSetFake                         = fakeclientset.NewSimpleClientset()
	reconciler                            *controllers.ByoMachineReconciler
	byoClusterReconciler                  *controllers.ByoClusterReconciler
	byoAdmissionReconciler                *controllers.ByoAdmissionReconciler
	k8sInstallerConfigReconciler          *controllers.K8sInstallerConfigReconciler
	bootstrapKubeconfigReconciler         *controllers.BootstrapKubeconfigReconciler
	recorder                              *record.FakeRecorder
	byoCluster                            *infrastructurev1beta1.ByoCluster
	capiCluster                           *clusterv1.Cluster
	defaultClusterName                    = "my-cluster"
	defaultNodeName                       = "my-host"
	defaultByoHostName                    = "my-host"
	defaultMachineName                    = "my-machine"
	defaultByoMachineName                 = "my-byomachine"
	defaultK8sInstallerConfigName         = "my-k8sinstallerconfig"
	defaultK8sInstallerConfigTemplateName = "my-installer-template"
	defaultNamespace                      = "default"
	fakeBootstrapSecret                   = "fakeBootstrapSecret"
	k8sManager                            ctrl.Manager
	cfg                                   *rest.Config
	ctx                                   context.Context
	cancel                                context.CancelFunc
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	var err error
	err = infrastructurev1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "config", "crd", "bases"),
			filepath.Join(build.Default.GOPATH, "pkg", "mod", "sigs.k8s.io", "cluster-api@v1.9.8", "config", "crd", "bases"),
			filepath.Join(build.Default.GOPATH, "pkg", "mod", "sigs.k8s.io", "cluster-api@v1.9.8", "bootstrap", "kubeadm", "config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing: true,
	}

	// Retrieve the first found binary directory to allow running tests from IDEs
	if getFirstFoundEnvTestBinaryDir() != "" {
		testEnv.BinaryAssetsDirectory = getFirstFoundEnvTestBinaryDir()
	}

	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = infrastructurev1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = clusterv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = bootstrapv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme:  scheme.Scheme,
		Metrics: server.Options{BindAddress: ":6080"},
	})
	Expect(err).NotTo(HaveOccurred())

	byoCluster = builder.ByoCluster(defaultNamespace, defaultClusterName).
		WithBundleBaseRegistry("projects.registry.vmware.com/cluster_api_provider_bringyourownhost").
		WithBundleTag("1.0").
		Build()
	Expect(k8sManager.GetClient().Create(context.Background(), byoCluster)).Should(Succeed())

	capiCluster = builder.Cluster(defaultNamespace, defaultClusterName).WithInfrastructureRef(byoCluster).Build()
	Expect(k8sManager.GetClient().Create(context.Background(), capiCluster)).Should(Succeed())

	node := builder.Node(defaultNamespace, defaultNodeName).Build()
	k8sClient = fake.NewClientBuilder().WithObjects(
		capiCluster,
		node,
	).Build()

	recorder = record.NewFakeRecorder(32)
	reconciler = &controllers.ByoMachineReconciler{
		Client: k8sManager.GetClient(),
		Tracker: remote.NewTestClusterCacheTracker(logr.New(logf.NullLogSink{}),
			k8sClient, k8sClient, scheme.Scheme,
			client.ObjectKey{
				Name: capiCluster.Name, Namespace: capiCluster.Namespace,
			}),
		Recorder: recorder,
	}
	err = reconciler.SetupWithManager(context.TODO(), k8sManager)
	Expect(err).NotTo(HaveOccurred())

	byoClusterReconciler = &controllers.ByoClusterReconciler{
		Client: k8sManager.GetClient(),
	}
	err = byoClusterReconciler.SetupWithManager(context.TODO(), k8sManager)
	Expect(err).NotTo(HaveOccurred())

	byoAdmissionReconciler = &controllers.ByoAdmissionReconciler{
		ClientSet: clientSetFake,
	}
	err = byoAdmissionReconciler.SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	k8sInstallerConfigReconciler = &controllers.K8sInstallerConfigReconciler{
		Client: k8sManager.GetClient(),
	}
	err = k8sInstallerConfigReconciler.SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	bootstrapKubeconfigReconciler = &controllers.BootstrapKubeconfigReconciler{
		Client: k8sManager.GetClient(),
	}
	err = bootstrapKubeconfigReconciler.SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	go func() {
		err = k8sManager.GetCache().Start(ctx)
		Expect(err).NotTo(HaveOccurred())
	}()

	Expect(k8sManager.GetCache().WaitForCacheSync(context.TODO())).To(BeTrue())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	cancel()
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func WaitForObjectsToBePopulatedInCache(objects ...client.Object) {
	for _, object := range objects {
		objectCopy := object.DeepCopyObject().(client.Object)
		key := client.ObjectKeyFromObject(object)
		Eventually(func() (done bool) {
			if err := reconciler.Client.Get(context.TODO(), key, objectCopy); err != nil {
				return false
			}
			return true
		}).Should(BeTrue())
	}
}

func WaitForObjectToBeUpdatedInCache(object client.Object, testObjectUpdatedFunc func(client.Object) bool) {
	objectCopy := object.DeepCopyObject().(client.Object)
	key := client.ObjectKeyFromObject(object)
	Eventually(func() (done bool) {
		if err := reconciler.Client.Get(context.TODO(), key, objectCopy); err != nil {
			return false
		}
		if testObjectUpdatedFunc(objectCopy) {
			return true
		}
		return false
	}).Should(BeTrue())
}

// getFirstFoundEnvTestBinaryDir locates the first binary in the specified path.
// ENVTEST-based tests depend on specific binaries, usually located in paths set by
// controller-runtime. When running tests directly (e.g., via an IDE) without using
// Makefile targets, the 'BinaryAssetsDirectory' must be explicitly configured.
//
// This function streamlines the process by finding the required binaries, similar to
// setting the 'KUBEBUILDER_ASSETS' environment variable. To ensure the binaries are
// properly set up, run 'make setup-envtest' beforehand.
func getFirstFoundEnvTestBinaryDir() string {
	basePath := filepath.Join("..", "..", "..", "bin", "k8s")
	entries, err := os.ReadDir(basePath)
	if err != nil {
		logf.Log.Error(err, "Failed to read directory", "path", basePath)
		return ""
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return filepath.Join(basePath, entry.Name())
		}
	}
	return ""
}
