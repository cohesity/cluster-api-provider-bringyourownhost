// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// K8sInstallerConfigTemplateSpec defines the desired state of K8sInstallerConfigTemplate.
type K8sInstallerConfigTemplateSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	Template K8sInstallerConfigTemplateResource `json:"template"`
}

// K8sInstallerConfigTemplateStatus defines the observed state of K8sInstallerConfigTemplate.
type K8sInstallerConfigTemplateStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// K8sInstallerConfigTemplate is the Schema for the k8sinstallerconfigtemplates API.
type K8sInstallerConfigTemplate struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of K8sInstallerConfigTemplate
	// +required
	Spec K8sInstallerConfigTemplateSpec `json:"spec"`

	// status defines the observed state of K8sInstallerConfigTemplate
	// +optional
	Status K8sInstallerConfigTemplateStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// K8sInstallerConfigTemplateList contains a list of K8sInstallerConfigTemplate.
type K8sInstallerConfigTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []K8sInstallerConfigTemplate `json:"items"`
}

type K8sInstallerConfigTemplateResource struct {
	// Spec is the specification of the desired behavior of the installer config.
	Spec K8sInstallerConfigSpec `json:"spec"`
}

func init() {
	SchemeBuilder.Register(&K8sInstallerConfigTemplate{}, &K8sInstallerConfigTemplateList{})
}
