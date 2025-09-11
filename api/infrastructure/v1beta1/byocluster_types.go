// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// ClusterFinalizer allows ReconcileByoCluster to clean up Byo
	// resources associated with ByoCluster before removing it from the
	// API server.
	ClusterFinalizer = "byocluster.infrastructure.cluster.x-k8s.io"
)

// ByoClusterSpec defines the desired state of ByoCluster.
type ByoClusterSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint APIEndpoint `json:"controlPlaneEndpoint"`

	// BundleLookupBaseRegistry is the base Registry URL that is used for pulling byoh bundle images,
	// if not set, the default will be set to https://projects.registry.vmware.com/cluster_api_provider_bringyourownhost
	// +optional
	BundleLookupBaseRegistry string `json:"bundleLookupBaseRegistry,omitempty"`
}

// ByoClusterStatus defines the observed state of ByoCluster.
type ByoClusterStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	Ready bool `json:"ready,omitempty"`

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the ByoCluster resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// FailureDomains is a list of failure domain objects synced from the infrastructure provider.
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`
}

// APIEndpoint represents a reachable Kubernetes API endpoint.
type APIEndpoint struct {
	// Host is the hostname on which the API server is serving.
	Host string `json:"host"`

	// Port is the port on which the API server is serving.
	Port int32 `json:"port"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=byoclusters,scope=Namespaced,shortName=byoc
// +kubebuilder:subresource:status

// ByoCluster is the Schema for the byoclusters API.
type ByoCluster struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of ByoCluster
	// +required
	Spec ByoClusterSpec `json:"spec"`

	// status defines the observed state of ByoCluster
	// +optional
	Status ByoClusterStatus `json:"status,omitempty,omitzero"`
}

// GetConditions gets the condition for the ByoCluster status
func (byoCluster *ByoCluster) GetConditions() clusterv1.Conditions {
	return byoCluster.Status.Conditions
}

// SetConditions sets the conditions for the ByoCluster status
func (byoCluster *ByoCluster) SetConditions(conditions clusterv1.Conditions) {
	byoCluster.Status.Conditions = conditions
}

// GetV1Beta1Conditions gets the ByoCluster status conditions for v1beta1 compatibility
func (byoCluster *ByoCluster) GetV1Beta1Conditions() clusterv1.Conditions {
	return byoCluster.Status.Conditions
}

// SetV1Beta1Conditions sets the ByoCluster status conditions for v1beta1 compatibility
func (byoCluster *ByoCluster) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	byoCluster.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// ByoClusterList contains a list of ByoCluster.
type ByoClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ByoCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ByoCluster{}, &ByoClusterList{})
}
