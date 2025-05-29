// Copyright 2025 Cohesity, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package infrastructure_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrastructurev1beta1 "github.com/cohesity/cluster-api-provider-bringyourownhost/api/infrastructure/v1beta1"

	. "github.com/cohesity/cluster-api-provider-bringyourownhost/internal/controller/infrastructure"
)

var _ = Describe("ByoMachineTemplate Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		byomachinetemplate := &infrastructurev1beta1.ByoMachineTemplate{}

		BeforeEach(func(ctx SpecContext) {
			By("creating the custom resource for the Kind ByoMachineTemplate")
			err := k8sClient.Get(ctx, typeNamespacedName, byomachinetemplate)
			if err != nil && errors.IsNotFound(err) {
				resource := &infrastructurev1beta1.ByoMachineTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					// TODO(user): Specify other spec details if needed.
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func(ctx SpecContext) {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &infrastructurev1beta1.ByoMachineTemplate{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance ByoMachineTemplate")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func(ctx SpecContext) {
			By("Reconciling the created resource")
			controllerReconciler := &ByoMachineTemplateReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
