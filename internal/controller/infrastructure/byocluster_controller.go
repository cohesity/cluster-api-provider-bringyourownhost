// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package infrastructure

import (
	"context"
	"reflect"
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrastructurev1beta1 "github.com/cohesity/cluster-api-provider-bringyourownhost/api/infrastructure/v1beta1"
)

var (
	// DefaultAPIEndpointPort default port for the API endpoint
	DefaultAPIEndpointPort    = 6443
	clusterControlledType     = &infrastructurev1beta1.ByoCluster{}
	clusterControlledTypeName = reflect.TypeOf(clusterControlledType).Elem().Name()
	clusterControlledTypeGVK  = infrastructurev1beta1.GroupVersion.WithKind(clusterControlledTypeName)
)

// ByoClusterReconciler reconciles a ByoCluster object
type ByoClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=byoclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=byoclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=byoclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

// Reconcile handles the ByoCluster reconciliations as part of the kubernetes
// reconciliation loop which aims to move the current state of the cluster
// closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *ByoClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	logger := log.FromContext(ctx)

	// Get the ByoCluster resource for this request.
	byoCluster := &infrastructurev1beta1.ByoCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, byoCluster); err != nil {
		if apierrors.IsNotFound(err) {
			logger.V(4).Info("ByoCluster not found, won't reconcile", "key", req.NamespacedName)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the CAPI Cluster.
	cluster, err := clusterutilv1.GetOwnerCluster(ctx, r.Client, byoCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		logger.Info("Waiting for Cluster Controller to set OwnerRef on ByoCluster")
		return reconcile.Result{}, nil
	}
	if annotations.IsPaused(cluster, byoCluster) {
		logger.V(4).Info("ByoCluster %s/%s linked to a cluster that is paused",
			byoCluster.Namespace, byoCluster.Name)
		return reconcile.Result{}, nil
	}

	// Create the patch helper.
	patchHelper, err := patch.NewHelper(byoCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(
			err,
			"failed to init patch helper for %s %s/%s",
			byoCluster.GroupVersionKind(),
			byoCluster.Namespace,
			byoCluster.Name)
	}

	// Always issue a patch when exiting this function so changes to the
	// resource are patched back to the API server.
	defer func() {
		if err := patchByoCluster(ctx, patchHelper, byoCluster); err != nil {
			logger.Error(err, "failed to patch ByoCluster")
			if reterr == nil {
				reterr = err
			}
		}
	}()

	// Handle deleted clusters
	if !byoCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, byoCluster)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, byoCluster)
}

func patchByoCluster(ctx context.Context, patchHelper *patch.Helper, byoCluster *infrastructurev1beta1.ByoCluster) error {
	// Always update the readyCondition by summarizing the state of other conditions.
	// A step counter is added to represent progress during the provisioning process (instead we are hiding it during the deletion process).
	conditions.SetSummary(byoCluster,
		conditions.WithStepCounterIf(byoCluster.ObjectMeta.DeletionTimestamp.IsZero()),
	)

	// Patch the object, ignoring conflicts on the conditions owned by this controller.
	return patchHelper.Patch(
		ctx,
		byoCluster,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
		}},
	)
}

// GetByoMachinesInCluster gets a cluster's ByoMachine resources.
func GetByoMachinesInCluster(
	ctx context.Context,
	controllerClient client.Client,
	namespace, clusterName string,
) ([]*infrastructurev1beta1.ByoMachine, error) {
	labels := map[string]string{clusterv1.ClusterNameLabel: clusterName}
	machineList := &infrastructurev1beta1.ByoMachineList{}

	if err := controllerClient.List(
		ctx, machineList,
		client.InNamespace(namespace),
		client.MatchingLabels(labels)); err != nil {
		return nil, err
	}

	machines := make([]*infrastructurev1beta1.ByoMachine, len(machineList.Items))
	for i := range machineList.Items {
		machines[i] = &machineList.Items[i]
	}

	return machines, nil
}

func (r ByoClusterReconciler) reconcileDelete(ctx context.Context, byoCluster *infrastructurev1beta1.ByoCluster) (reconcile.Result, error) {
	logger := log.FromContext(ctx)

	byoMachines, err := GetByoMachinesInCluster(ctx, r.Client, byoCluster.Namespace, byoCluster.Name)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err,
			"unable to list ByoMachines part of ByoCluster %s/%s", byoCluster.Namespace, byoCluster.Name)
	}

	if len(byoMachines) > 0 {
		logger.Info("Waiting for ByoMachines to be deleted", "count", len(byoMachines))
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}
	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(byoCluster, infrastructurev1beta1.ClusterFinalizer)

	return ctrl.Result{}, nil
}

func (r ByoClusterReconciler) reconcileNormal(ctx context.Context, byoCluster *infrastructurev1beta1.ByoCluster) (reconcile.Result, error) {
	// If the ByoCluster doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(byoCluster, infrastructurev1beta1.ClusterFinalizer)

	if byoCluster.Spec.ControlPlaneEndpoint.Port == 0 {
		byoCluster.Spec.ControlPlaneEndpoint.Port = int32(DefaultAPIEndpointPort)
	}

	byoCluster.Status.Ready = true

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ByoClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Watch the controlled, infrastructure resource.
		For(clusterControlledType).
		// Watch the CAPI resource that owns this infrastructure resource.
		Watches(&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(clusterutilv1.ClusterToInfrastructureMapFunc(ctx, infrastructurev1beta1.GroupVersion.WithKind(clusterControlledTypeGVK.Kind), mgr.GetClient(), &infrastructurev1beta1.ByoCluster{})),
		).
		Named("infrastructure-byocluster").
		Complete(r)
}
