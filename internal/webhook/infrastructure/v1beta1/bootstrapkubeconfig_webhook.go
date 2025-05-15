// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	"context"
	b64 "encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrastructurev1beta1 "github.com/cohesity/cluster-api-provider-bringyourownhost/api/infrastructure/v1beta1"
)

// log is for logging in this package.
var bootstrapkubeconfiglog = logf.Log.WithName("bootstrapkubeconfig-resource")

// APIServerURLScheme is the url scheme for the APIServer
const APIServerURLScheme = "https"

func SetupBootstrapKubeconfigWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&infrastructurev1beta1.BootstrapKubeconfig{}).
		WithValidator(&BootstrapKubeconfigCustomValidator{}).
		Complete()
}

// +kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-bootstrapkubeconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=bootstrapkubeconfigs,verbs=create;update,versions=v1beta1,name=vbootstrapkubeconfig.kb.io,admissionReviewVersions=v1

// BootstrapKubeconfigCustomValidator struct is responsible for validating the BootstrapKubeconfig resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and
type BootstrapKubeconfigCustomValidator struct{}

var _ webhook.CustomValidator = &BootstrapKubeconfigCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type
func (v *BootstrapKubeconfigCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	bootstrapKubeconfig, ok := obj.(*infrastructurev1beta1.BootstrapKubeconfig)
	if !ok {
		return nil, fmt.Errorf("expected a BootstrapKubeconfig object but got %T", obj)
	}

	bootstrapkubeconfiglog.Info("Validation for BootstrapKubeconfig upon creation", "name", bootstrapKubeconfig.GetName())

	var allErrs field.ErrorList

	if err := validateAPIServer(bootstrapKubeconfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := validateCAData(bootstrapKubeconfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "BootstrapKubeconfig"},
		bootstrapKubeconfig.GetName(), allErrs)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type
func (v *BootstrapKubeconfigCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	bootstrapKubeconfig, ok := newObj.(*infrastructurev1beta1.BootstrapKubeconfig)
	if !ok {
		return nil, fmt.Errorf("expected a BootstrapKubeconfig object for the newObj but got %T", newObj)
	}

	bootstrapkubeconfiglog.Info("Validation for BootstrapKubeconfig upon update", "name", bootstrapKubeconfig.GetName())

	var allErrs field.ErrorList

	if err := validateAPIServer(bootstrapKubeconfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := validateCAData(bootstrapKubeconfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "BootstrapKubeconfig"},
		bootstrapKubeconfig.GetName(), allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type
func (v *BootstrapKubeconfigCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*infrastructurev1beta1.BootstrapKubeconfig)
	if !ok {
		return nil, fmt.Errorf("expected a BootstrapKubeconfig object but got %T", obj)
	}

	bootstrapkubeconfiglog.Info("Validation for BootstrapKubeconfig upon delete", "name", r.GetName())

	return nil, nil
}

func validateAPIServer(r *infrastructurev1beta1.BootstrapKubeconfig) field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.APIServer == "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("apiserver"), r.Spec.APIServer, "APIServer field cannot be empty"))
	}

	parsedURL, err := url.Parse(r.Spec.APIServer)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("apiserver"), r.Spec.APIServer, "APIServer URL is not valid"))
	}
	if !isURLValid(parsedURL) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("apiserver"), r.Spec.APIServer, "APIServer is not of the format https://hostname:port"))
	}

	return allErrs
}

func validateCAData(r *infrastructurev1beta1.BootstrapKubeconfig) field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.CertificateAuthorityData == "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("caData"), r.Spec.CertificateAuthorityData, "CertificateAuthorityData field cannot be empty"))
	}

	decodedCAData, err := b64.StdEncoding.DecodeString(r.Spec.CertificateAuthorityData)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("caData"), r.Spec.CertificateAuthorityData, "cannot base64 decode CertificateAuthorityData"))
	}

	block, _ := pem.Decode(decodedCAData)
	if block == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("caData"), r.Spec.CertificateAuthorityData, "CertificateAuthorityData is not PEM encoded"))
	}

	return allErrs
}

func isURLValid(parsedURL *url.URL) bool {
	if parsedURL.Host == "" || parsedURL.Scheme != APIServerURLScheme || parsedURL.Port() == "" {
		return false
	}
	return true
}
