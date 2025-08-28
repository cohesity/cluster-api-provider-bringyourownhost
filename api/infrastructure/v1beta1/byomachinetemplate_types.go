// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ByoMachineTemplateSpec defines the desired state of ByoMachineTemplate.
type ByoMachineTemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	Template ByoMachineTemplateResource `json:"template"`
}

// ByoMachineTemplateStatus defines the observed state of ByoMachineTemplate.
type ByoMachineTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ByoMachineTemplate is the Schema for the byomachinetemplates API.
type ByoMachineTemplate struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of ByoMachineTemplate
	// +required
	Spec ByoMachineTemplateSpec `json:"spec"`

	// status defines the observed state of ByoMachineTemplate
	// +optional
	Status ByoMachineTemplateStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ByoMachineTemplateList contains a list of ByoMachineTemplate.
type ByoMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ByoMachineTemplate `json:"items"`
}

// ByoMachineTemplateResource defines the desired state of ByoMachineTemplateResource
type ByoMachineTemplateResource struct {
	// Spec is the specification of the desired behavior of the machine.
	Spec ByoMachineSpec `json:"spec"`
}

func init() {
	SchemeBuilder.Register(&ByoMachineTemplate{}, &ByoMachineTemplateList{})
}
