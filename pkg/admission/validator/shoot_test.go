// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"context"
	"testing"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	configinstall "github.com/gardener/gardener-extension-shoot-traefik/pkg/apis/config/install"
)

const (
	testAPIVersion = "core.gardener.cloud/v1beta1"
	testKind       = "Shoot"
	testShootName  = "test-shoot"
	testNamespace  = "garden-test"
)

func TestValidator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Admission Validator Suite")
}

var _ = Describe("Shoot Validator", func() {
	var (
		validator *shootValidator
		scheme    *runtime.Scheme
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(gardencorev1beta1.AddToScheme(scheme)).To(Succeed())
		configinstall.Install(scheme)

		client := fake.NewClientBuilder().WithScheme(scheme).Build()
		decoder := serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()
		validator = &shootValidator{
			client:  client,
			decoder: decoder,
		}
	})

	Context("when shoot has traefik extension", func() {
		It("should allow shoot with purpose 'evaluation'", func() {
			purpose := gardencorev1beta1.ShootPurposeEvaluation
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: testAPIVersion,
					Kind:       testKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      testShootName,
					Namespace: testNamespace,
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: &purpose,
					Extensions: []gardencorev1beta1.Extension{
						{Type: ExtensionType},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should deny shoot with purpose 'production'", func() {
			purpose := gardencorev1beta1.ShootPurposeProduction
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: testAPIVersion,
					Kind:       testKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      testShootName,
					Namespace: testNamespace,
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: &purpose,
					Extensions: []gardencorev1beta1.Extension{
						{Type: ExtensionType},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("evaluation"))
		})

		It("should deny shoot with nil purpose", func() {
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: testAPIVersion,
					Kind:       testKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      testShootName,
					Namespace: testNamespace,
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: nil,
					Extensions: []gardencorev1beta1.Extension{
						{Type: ExtensionType},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should deny shoot with purpose 'development'", func() {
			purpose := gardencorev1beta1.ShootPurposeDevelopment
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: testAPIVersion,
					Kind:       testKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      testShootName,
					Namespace: testNamespace,
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: &purpose,
					Extensions: []gardencorev1beta1.Extension{
						{Type: ExtensionType},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("evaluation"))
		})

		It("should deny shoot with the nginx-ingress addon enabled", func() {
			purpose := gardencorev1beta1.ShootPurposeEvaluation
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: testAPIVersion,
					Kind:       testKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      testShootName,
					Namespace: testNamespace,
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: &purpose,
					Addons: &gardencorev1beta1.Addons{
						NginxIngress: &gardencorev1beta1.NginxIngress{
							Addon: gardencorev1beta1.Addon{Enabled: true},
						},
					},
					Extensions: []gardencorev1beta1.Extension{
						{Type: ExtensionType},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("nginx-ingress addon"))
		})

		It("should allow shoot with the nginx-ingress addon disabled", func() {
			purpose := gardencorev1beta1.ShootPurposeEvaluation
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: testAPIVersion,
					Kind:       testKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      testShootName,
					Namespace: testNamespace,
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: &purpose,
					Addons: &gardencorev1beta1.Addons{
						NginxIngress: &gardencorev1beta1.NginxIngress{
							Addon: gardencorev1beta1.Addon{Enabled: false},
						},
					},
					Extensions: []gardencorev1beta1.Extension{
						{Type: ExtensionType},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when shoot has traefik extension with providerConfig", func() {
		newShootWithProviderConfig := func(rawConfig string) *gardencorev1beta1.Shoot {
			purpose := gardencorev1beta1.ShootPurposeEvaluation
			ext := gardencorev1beta1.Extension{Type: ExtensionType}
			if rawConfig != "" {
				ext.ProviderConfig = &runtime.RawExtension{Raw: []byte(rawConfig)}
			}

			return &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: testAPIVersion,
					Kind:       testKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      testShootName,
					Namespace: testNamespace,
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose:    &purpose,
					Extensions: []gardencorev1beta1.Extension{ext},
				},
			}
		}

		It("should allow a valid httpEntrypoint", func() {
			shoot := newShootWithProviderConfig(`{"apiVersion":"traefik.extensions.gardener.cloud/v1alpha1","kind":"TraefikConfig","spec":{"httpEntrypoint":"Redirect"}}`)

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should deny an invalid httpEntrypoint", func() {
			shoot := newShootWithProviderConfig(`{"apiVersion":"traefik.extensions.gardener.cloud/v1alpha1","kind":"TraefikConfig","spec":{"httpEntrypoint":"Bogus"}}`)

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("httpEntrypoint"))
		})

		It("should deny an invalid logLevel", func() {
			shoot := newShootWithProviderConfig(`{"apiVersion":"traefik.extensions.gardener.cloud/v1alpha1","kind":"TraefikConfig","spec":{"logLevel":"Loud"}}`)

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("logLevel"))
		})

		It("should allow an empty httpEntrypoint (defaulted later)", func() {
			shoot := newShootWithProviderConfig(`{"apiVersion":"traefik.extensions.gardener.cloud/v1alpha1","kind":"TraefikConfig","spec":{"replicas":3}}`)

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should allow a valid httpEntrypoint in the flat v1alpha2 form", func() {
			shoot := newShootWithProviderConfig(`{"apiVersion":"traefik.extensions.gardener.cloud/v1alpha2","kind":"TraefikConfig","httpEntrypoint":"Redirect"}`)

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should deny an invalid httpEntrypoint in the flat v1alpha2 form", func() {
			shoot := newShootWithProviderConfig(`{"apiVersion":"traefik.extensions.gardener.cloud/v1alpha2","kind":"TraefikConfig","httpEntrypoint":"Bogus"}`)

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("httpEntrypoint"))
		})

		It("should deny an invalid logLevel in the flat v1alpha2 form", func() {
			shoot := newShootWithProviderConfig(`{"apiVersion":"traefik.extensions.gardener.cloud/v1alpha2","kind":"TraefikConfig","logLevel":"Loud"}`)

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("logLevel"))
		})
	})

	Context("when shoot does not have traefik extension", func() {
		It("should allow shoot without traefik extension regardless of purpose", func() {
			purpose := gardencorev1beta1.ShootPurposeProduction
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: testAPIVersion,
					Kind:       testKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      testShootName,
					Namespace: testNamespace,
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose: &purpose,
					Extensions: []gardencorev1beta1.Extension{
						{Type: "other-extension"},
					},
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should allow shoot with no extensions", func() {
			purpose := gardencorev1beta1.ShootPurposeProduction
			shoot := &gardencorev1beta1.Shoot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: testAPIVersion,
					Kind:       testKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      testShootName,
					Namespace: testNamespace,
				},
				Spec: gardencorev1beta1.ShootSpec{
					Purpose:    &purpose,
					Extensions: nil,
				},
			}

			err := validator.Validate(context.Background(), shoot, nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
