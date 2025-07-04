// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1beta1_test

import (
	b64 "encoding/base64"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"

	infrastructurev1beta1 "github.com/cohesity/cluster-api-provider-bringyourownhost/api/infrastructure/v1beta1"
	"github.com/cohesity/cluster-api-provider-bringyourownhost/test/builder"

	. "github.com/cohesity/cluster-api-provider-bringyourownhost/internal/webhook/infrastructure/v1beta1"
)

var _ = Describe("BootstrapKubeconfig Webhook", func() {
	var (
		obj       *infrastructurev1beta1.BootstrapKubeconfig
		oldObj    *infrastructurev1beta1.BootstrapKubeconfig
		validator BootstrapKubeconfigCustomValidator
	)

	var (
		err                         error
		defaultNamespace            = "default"
		testBootstrapKubeconfigName = "test-bootstrap-kubeconfig"
		testServerEmpty             = ""
		testServerInvalidURL        = "htt p://test.com"
		testServerWithoutScheme     = "abc.com"
		testServerWithoutHostname   = "https://test-server"
		testServerWithoutPort       = "https://test.com"
		testServerValid             = "https://abc.com:1234"
		testCADataEmpty             = ""
		testCADataInvalid           = "test-ca-data"
		testPEMDataInvalid          = b64.StdEncoding.EncodeToString([]byte(testCADataInvalid))
	)

	BeforeEach(func() {
		obj = &infrastructurev1beta1.BootstrapKubeconfig{}
		oldObj = &infrastructurev1beta1.BootstrapKubeconfig{}
		validator = BootstrapKubeconfigCustomValidator{}
		Expect(validator).NotTo(BeNil(), "Expected validator to be initialized")
		Expect(oldObj).NotTo(BeNil(), "Expected oldObj to be initialized")
		Expect(obj).NotTo(BeNil(), "Expected obj to be initialized")
		// TODO (user): Add any setup logic common to all tests
	})

	AfterEach(func() {
		// TODO (user): Add any teardown logic common to all tests
	})

	Context("When creating BootstrapKubeconfig under Validating Webhook", func() {
		It("Should deny creation if APIServer is not a valid URL", func() {
			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerInvalidURL).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("spec.apiserver: Invalid value: %q: APIServer URL is not valid", testServerInvalidURL))))
		})

		It("Should deny creation if APIServer field is empty", func() {
			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerEmpty).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring("spec.apiserver: Invalid value: \"\": APIServer field cannot be empty")))
		})

		It("Should deny creation if APIServer address does not have https scheme specified", func() {
			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerWithoutScheme).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("spec.apiserver: Invalid value: %q: APIServer is not of the format https://hostname:port", testServerWithoutScheme))))
		})

		It("Should deny creation if  APIServer address hostname is not specified", func() {
			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerWithoutHostname).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("spec.apiserver: Invalid value: %q: APIServer is not of the format https://hostname:port", testServerWithoutHostname))))
		})

		It("Should deny creation if APIServer address does not have the port info", func() {
			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerWithoutPort).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("spec.apiserver: Invalid value: %q: APIServer is not of the format https://hostname:port", testServerWithoutPort))))
		})

		It("Should deny creation if CertificateAuthorityData field is empty", func() {
			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerValid).
				WithCAData(testCADataEmpty).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring("spec.caData: Invalid value: \"\": CertificateAuthorityData field cannot be empty")))
		})

		It("Should deny creation if CertificateAuthorityData cannot be base64 decoded", func() {
			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerValid).
				WithCAData(testCADataInvalid).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("spec.caData: Invalid value: %q: cannot base64 decode CertificateAuthorityData", testCADataInvalid))))
		})

		It("Should deny creation if CertificateAuthorityData is not PEM encoded", func() {
			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerValid).
				WithCAData(testPEMDataInvalid).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("spec.caData: Invalid value: %q: CertificateAuthorityData is not PEM encoded", testPEMDataInvalid))))
		})

		It("Should admit creation if all fields are valid", func() {
			// use from config of envtest
			testCADataValid := b64.StdEncoding.EncodeToString(cfg.CAData)

			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerValid).
				WithCAData(testCADataValid).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When updating BootstrapKubeconfig under Validating Webhook", func() {
		// It("Should validate updates correctly", func() {
		//     By("simulating a valid update scenario")
		var (
			ph                         *patch.Helper
			createdBootstrapKubeconfig *infrastructurev1beta1.BootstrapKubeconfig
		)
		BeforeEach(func() {
			// use from config of envtest
			testCADataValid := b64.StdEncoding.EncodeToString(cfg.CAData)

			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerValid).
				WithCAData(testCADataValid).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).NotTo(HaveOccurred())

			createdBootstrapKubeconfig = &infrastructurev1beta1.BootstrapKubeconfig{}
			namespacedName := types.NamespacedName{Name: obj.Name, Namespace: defaultNamespace}
			Eventually(func() error {
				err = k8sClient.Get(ctx, namespacedName, createdBootstrapKubeconfig)
				if err != nil {
					return err
				}
				return nil
			}).Should(BeNil())

			// create a patch helper
			ph, err = patch.NewHelper(obj, k8sClient)
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			err = k8sClient.Delete(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should deny update if APIServer field is empty", func() {
			createdBootstrapKubeconfig.Spec.APIServer = testServerEmpty
			err = ph.Patch(ctx, createdBootstrapKubeconfig)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring("spec.apiserver: Invalid value: \"\": APIServer field cannot be empty")))
		})

		It("Should deny update if APIServer is not of the correct format", func() {
			createdBootstrapKubeconfig.Spec.APIServer = testServerWithoutHostname
			err = ph.Patch(ctx, createdBootstrapKubeconfig)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("spec.apiserver: Invalid value: %q: APIServer is not of the format https://hostname:port", testServerWithoutHostname))))
		})

		It("Should deny update if CertificateAuthorityData field is empty", func() {
			createdBootstrapKubeconfig.Spec.CertificateAuthorityData = testCADataEmpty
			err = ph.Patch(ctx, createdBootstrapKubeconfig)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring("spec.caData: Invalid value: \"\": CertificateAuthorityData field cannot be empty")))
		})

		It("Should deny update if CertificateAuthorityData cannot be base64 decoded", func() {
			createdBootstrapKubeconfig.Spec.CertificateAuthorityData = testCADataInvalid
			err = ph.Patch(ctx, createdBootstrapKubeconfig)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("spec.caData: Invalid value: %q: cannot base64 decode CertificateAuthorityData", testCADataInvalid))))
		})

		It("Should deny update if CertificateAuthorityData is not PEM encoded", func() {
			createdBootstrapKubeconfig.Spec.CertificateAuthorityData = testPEMDataInvalid
			err = ph.Patch(ctx, createdBootstrapKubeconfig)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("admission webhook \"vbootstrapkubeconfig-v1beta1.kb.io\" denied the request")))
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("spec.caData: Invalid value: %q: CertificateAuthorityData is not PEM encoded", testPEMDataInvalid))))
		})

		It("Should admit update if all fields are valid", func() {
			// patch a valid APIServer value
			createdBootstrapKubeconfig.Spec.APIServer = "https://1.2.3.4:5678"
			err = ph.Patch(ctx, createdBootstrapKubeconfig)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When deleting BootstrapKubeconfig under Validating Webhook", func() {
		var (
			createdBootstrapKubeconfig *infrastructurev1beta1.BootstrapKubeconfig
			namespacedName             types.NamespacedName
		)
		BeforeEach(func() {
			// use from config of envtest
			testCADataValid := b64.StdEncoding.EncodeToString(cfg.CAData)

			obj = builder.BootstrapKubeconfig(defaultNamespace, testBootstrapKubeconfigName).
				WithServer(testServerValid).
				WithCAData(testCADataValid).
				Build()
			err = k8sClient.Create(ctx, obj)
			Expect(err).NotTo(HaveOccurred())

			createdBootstrapKubeconfig = &infrastructurev1beta1.BootstrapKubeconfig{}
			namespacedName = types.NamespacedName{Name: obj.Name, Namespace: defaultNamespace}
			Eventually(func() error {
				err = k8sClient.Get(ctx, namespacedName, createdBootstrapKubeconfig)
				if err != nil {
					return err
				}
				return nil
			}).Should(BeNil())
		})

		It("Should admit delete always", func() {
			err = k8sClient.Delete(ctx, createdBootstrapKubeconfig)
			Expect(err).NotTo(HaveOccurred())

			deletedBootstrapKubeconfig := &infrastructurev1beta1.BootstrapKubeconfig{}
			Eventually(func() bool {
				err = k8sClient.Get(ctx, namespacedName, deletedBootstrapKubeconfig)
				if err != nil {
					if apierrors.IsNotFound(err) {
						return true
					}
				}
				return false
			}).Should(BeTrue())
		})
	})
})
