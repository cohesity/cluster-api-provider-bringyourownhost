// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrastructurev1beta1 "github.com/cohesity/cluster-api-provider-bringyourownhost/api/infrastructure/v1beta1"
)

// nolint:unused
// log is for logging in this package.
var byohostlog = logf.Log.WithName("byohost-resource")

// SetupByoHostWebhookWithManager registers the webhook for ByoHost in the manager.
func SetupByoHostWebhookWithManager(mgr ctrl.Manager) error {
	// return ctrl.NewWebhookManagedBy(mgr).For(&infrastructurev1beta1.ByoHost{}).
	// 	WithValidator(&ByoHostCustomValidator{}).
	// 	Complete()

	schema := mgr.GetScheme()
	decoder := admission.NewDecoder(schema)
	validator := &ByoHostCustomValidator{
		Decoder: decoder,
	}

	customValidatorHandler := admission.WithCustomValidator(schema, &infrastructurev1beta1.ByoHost{}, validator)

	handler := admission.MultiValidatingHandler(customValidatorHandler, validator)

	webhookHandler := &webhook.Admission{Handler: handler}

	mgr.GetWebhookServer().Register("/validate-infrastructure-cluster-x-k8s-io-v1beta1-byohost", webhookHandler)

	return nil
}

// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-byohost,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=byohosts,verbs=create;update;delete,versions=v1beta1,name=vbyohost-v1beta1.kb.io,admissionReviewVersions=v1

// ByoHostCustomValidator struct is responsible for validating the ByoHost resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ByoHostCustomValidator struct {
	Decoder admission.Decoder
}

var _ webhook.CustomValidator = &ByoHostCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type ByoHost.
func (v *ByoHostCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	byohost, ok := obj.(*infrastructurev1beta1.ByoHost)
	if !ok {
		return nil, fmt.Errorf("expected a ByoHost object but got %T", obj)
	}
	byohostlog.Info("Validation for ByoHost upon creation", "name", byohost.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type ByoHost.
func (v *ByoHostCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	byohost, ok := newObj.(*infrastructurev1beta1.ByoHost)
	if !ok {
		return nil, fmt.Errorf("expected a ByoHost object for the newObj but got %T", newObj)
	}
	byohostlog.Info("Validation for ByoHost upon update", "name", byohost.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type ByoHost.
func (v *ByoHostCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	byohost, ok := obj.(*infrastructurev1beta1.ByoHost)
	if !ok {
		return nil, fmt.Errorf("expected a ByoHost object but got %T", obj)
	}
	byohostlog.Info("Validation for ByoHost upon deletion", "name", byohost.GetName())

	if byohost.Status.MachineRef != nil {
		return nil, errors.New("cannot delete ByoHost when MachineRef is assigned")
	}

	return nil, nil
}

// To allow byoh manager service account to patch ByoHost CR
const ManagerServiceAccount = "system:serviceaccount:byoh-system:byoh-controller-manager"

const (
	// minUsernamePartsForAgent is the minimum number of parts expected when splitting agent username
	// Agent usernames follow the pattern "system:serviceaccount:namespace:name" which has at least 2 parts when split by ":"
	minUsernamePartsForAgent = 2
)

// nolint: gocritic
// Handle handles all the requests for ByoHost resource
func (v *ByoHostCustomValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	var response admission.Response

	ctx = admission.NewContextWithRequest(ctx, req)

	switch req.Operation {
	case admissionv1.Create, admissionv1.Update:
		response = v.handleCreateUpdateReq(ctx, &req)
	// case admissionv1.Delete:
	// response = v.handleDelete(ctx, &req)
	default:
		response = admission.Allowed("")
	}
	return response
}

func (v *ByoHostCustomValidator) handleCreateUpdateReq(_ context.Context, req *admission.Request) admission.Response {
	byoHost := &infrastructurev1beta1.ByoHost{}

	err := v.Decoder.Decode(*req, byoHost)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	userName := req.UserInfo.Username
	// allow manager service account to patch ByoHost
	if userName == ManagerServiceAccount && req.Operation == admissionv1.Update {
		return admission.Allowed("")
	}

	substrs := strings.Split(userName, ":")
	if len(substrs) < minUsernamePartsForAgent {
		return admission.Denied(fmt.Sprintf("%s is not a valid agent username", userName))
	}

	if !strings.Contains(byoHost.Name, substrs[2]) {
		return admission.Denied(fmt.Sprintf("%s cannot create/update resource %s", userName, byoHost.Name))
	}

	return admission.Allowed("")
}

// func (v *ByoHostCustomValidator) handleDelete(ctx context.Context, req *admission.Request) admission.Response {
// 	byohost := &infrastructurev1beta1.ByoHost{}
// 	err := v.Decoder.DecodeRaw(req.OldObject, byohost)
// 	if err != nil {
// 		return admission.Errored(http.StatusBadRequest, err)
// 	}

// 	warnings, err := v.ValidateDelete(ctx, byohost)
// 	// Check the error message first.
// 	if err != nil {
// 		var apiStatus apierrors.APIStatus
// 		if errors.As(err, &apiStatus) {
// 			return admission.Response{
// 				AdmissionResponse: admissionv1.AdmissionResponse{
// 					Allowed: false,
// 					Result:  &apiStatus.Status(),
// 				},
// 			}.WithWarnings(warnings...)
// 		}
// 		return admission.Denied(err.Error()).WithWarnings(warnings...)
// 	}

// 	return admission.Allowed("").WithWarnings(warnings...)
// }
