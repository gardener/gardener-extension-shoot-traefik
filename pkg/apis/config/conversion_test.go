// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/gardener/gardener-extension-shoot-traefik/pkg/apis/config"
	configinstall "github.com/gardener/gardener-extension-shoot-traefik/pkg/apis/config/install"
)

func TestConfigAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config API Suite")
}

var _ = Describe("TraefikConfig decoding", func() {
	var decoder runtime.Decoder

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		configinstall.Install(scheme)
		decoder = serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()
	})

	// The core non-disruption guarantee: an existing shoot's v1alpha1 providerConfig
	// (with the "spec" wrapper) and a new v1alpha2 providerConfig (flat) must decode to
	// exactly the same internal configuration.
	It("should decode v1alpha1 (spec-wrapped) and v1alpha2 (flat) to the same internal config", func() {
		v1alpha1Raw := []byte(`{
			"apiVersion": "traefik.extensions.gardener.cloud/v1alpha1",
			"kind": "TraefikConfig",
			"spec": {
				"replicas": 3,
				"ingressProvider": "KubernetesIngressNGINX",
				"logLevel": "Info",
				"dashboard": true,
				"httpEntrypoint": "Redirect"
			}
		}`)

		v1alpha2Raw := []byte(`{
			"apiVersion": "traefik.extensions.gardener.cloud/v1alpha2",
			"kind": "TraefikConfig",
			"replicas": 3,
			"ingressProvider": "KubernetesIngressNGINX",
			"logLevel": "Info",
			"dashboard": true,
			"httpEntrypoint": "Redirect"
		}`)

		var fromV1alpha1, fromV1alpha2 config.TraefikConfig
		Expect(runtime.DecodeInto(decoder, v1alpha1Raw, &fromV1alpha1)).To(Succeed())
		Expect(runtime.DecodeInto(decoder, v1alpha2Raw, &fromV1alpha2)).To(Succeed())

		expected := config.TraefikConfig{
			Replicas:        3,
			IngressProvider: config.IngressProviderKubernetesIngressNGINX,
			LogLevel:        "Info",
			Dashboard:       true,
			HTTPEntrypoint:  config.HTTPEntrypointRedirect,
		}

		// Ignore TypeMeta, which differs by apiVersion/kind of the source object.
		fromV1alpha1.TypeMeta = expected.TypeMeta
		fromV1alpha2.TypeMeta = expected.TypeMeta

		Expect(fromV1alpha1).To(Equal(expected))
		Expect(fromV1alpha2).To(Equal(expected))
	})
})
