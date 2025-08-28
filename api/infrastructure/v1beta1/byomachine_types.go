// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// MachineFinalizer allows ReconcileByoMachine to clean up Byo
	// resources associated with ByoMachine before removing it from the
	// API Server.
	MachineFinalizer = "byomachine.infrastructure.cluster.x-k8s.io"
)

// ByoMachineSpec defines the desired state of ByoMachine.
type ByoMachineSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// Label Selector to choose the byohost
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	ProviderID string `json:"providerID,omitempty"`

	// InstallerRef is an optional reference to a installer-specific resource that holds
	// the details of InstallationSecret to be used to install BYOH Bundle.
	// +optional
	InstallerRef *corev1.ObjectReference `json:"installerRef,omitempty"`
}

// NetworkStatus provides information about one of a VM's networks.
type NetworkStatus struct {
	// Connected is a flag that indicates whether this network is currently
	// connected to the VM.
	Connected bool `json:"connected,omitempty"`

	// IPAddrs is one or more IP addresses reported by vm-tools.
	// +optional
	IPAddrs []string `json:"ipAddrs,omitempty"`

	// MACAddr is the MAC address of the network device.
	MACAddr string `json:"macAddr"`

	// NetworkInterfaceName is the name of the network interface.
	// +optional
	NetworkInterfaceName string `json:"networkInterfaceName,omitempty"`

	// IsDefault is a flag that indicates whether this interface name is where
	// the default gateway sit on.
	IsDefault bool `json:"isDefault,omitempty"`
}

// ByoMachineStatus defines the observed state of ByoMachine.
type ByoMachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// HostInfo has the attached host platform details.
	// +optional
	HostInfo HostInfo `json:"hostinfo,omitempty"`

	// +optional
	Ready bool `json:"ready"`

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the ByoMachine resource.
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
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=byomachines,scope=Namespaced,shortName=byom
// +kubebuilder:subresource:status

// ByoMachine is the Schema for the byomachines API.
type ByoMachine struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of ByoMachine
	// +required
	Spec ByoMachineSpec `json:"spec"`

	// status defines the observed state of ByoMachine
	// +optional
	Status ByoMachineStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ByoMachineList contains a list of ByoMachine.
type ByoMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ByoMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ByoMachine{}, &ByoMachineList{})
}

// GetConditions returns the conditions of ByoMachine status
func (byoMachine *ByoMachine) GetConditions() clusterv1.Conditions {
	return byoMachine.Status.Conditions
}

// SetConditions sets the conditions of ByoMachine status
func (byoMachine *ByoMachine) SetConditions(conditions clusterv1.Conditions) {
	byoMachine.Status.Conditions = conditions
}
