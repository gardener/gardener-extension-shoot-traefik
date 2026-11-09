// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressProviderType defines the type of Kubernetes Ingress provider to use.
type IngressProviderType string

const (
	// IngressProviderKubernetesIngress is the standard Kubernetes Ingress provider.
	IngressProviderKubernetesIngress IngressProviderType = "KubernetesIngress"
	// IngressProviderKubernetesIngressNGINX is the NGINX-compatible Kubernetes Ingress provider.
	// This provider supports NGINX Ingress Controller annotations, making it easier to migrate
	// from NGINX Ingress Controller to Traefik.
	IngressProviderKubernetesIngressNGINX IngressProviderType = "KubernetesIngressNGINX"
)

// HTTPEntrypointType defines how the Traefik LoadBalancer handles plain HTTP (port 80).
type HTTPEntrypointType string

const (
	// HTTPEntrypointEnabled exposes the Service port 80 and serves plain HTTP (default).
	HTTPEntrypointEnabled HTTPEntrypointType = "Enabled"
	// HTTPEntrypointRedirect exposes the Service port 80 and redirects all HTTP requests to HTTPS (301).
	HTTPEntrypointRedirect HTTPEntrypointType = "Redirect"
	// HTTPEntrypointDisabled does not expose the Service port 80. The container web entrypoint (:8000)
	// is kept regardless, because the /ping health probes depend on it.
	HTTPEntrypointDisabled HTTPEntrypointType = "Disabled"
)

// TraefikConfigSpec defines the desired state of [TraefikConfig]
type TraefikConfigSpec struct {
	// Replicas is the number of Traefik replicas to deploy.
	// Defaults to 2 if not specified.
	Replicas int32 `json:"replicas,omitempty"`

	// IngressProvider specifies which Kubernetes Ingress provider to use.
	// Valid values are:
	// - "KubernetesIngress" (default): Standard Kubernetes Ingress provider
	// - "KubernetesIngressNGINX": NGINX-compatible provider with support for NGINX annotations
	//
	// Use KubernetesIngressNGINX when migrating from NGINX Ingress Controller to maintain
	// compatibility with existing NGINX-specific annotations.
	IngressProvider IngressProviderType `json:"ingressProvider,omitempty"`

	// LogLevel sets the Traefik log level.
	// Valid values are: Debug, Info, Warn, Error, Fatal, Panic
	// Defaults to "Info" if not specified.
	LogLevel string `json:"logLevel,omitempty"`

	// Dashboard enables the Traefik dashboard.
	// The dashboard is exposed on port 9000 and accessible via port-forwarding.
	// Enabling the API and the dashboard in production is not recommended, because it will expose all
	// configuration elements, including sensitive data, for which access should be reserved to administrators.
	// Defaults to false if not specified.
	Dashboard bool `json:"dashboard,omitempty"`

	// HTTPEntrypoint controls how the Traefik LoadBalancer handles plain HTTP (port 80).
	// Valid values are:
	// - "Enabled" (default): expose Service port 80 and serve plain HTTP.
	// - "Redirect": expose Service port 80 and redirect all HTTP requests to HTTPS (301).
	// - "Disabled": do not expose Service port 80 (HTTPS only).
	// Defaults to "Enabled" if not specified.
	HTTPEntrypoint HTTPEntrypointType `json:"httpEntrypoint,omitempty"`

	// ProxyProtocol configures the PROXY protocol on the Traefik LoadBalancer so that the
	// real client IP is preserved. This is required on AWS, where Gardener provisions a
	// Classic Load Balancer with TCP (L4) listeners that discard the client IP unless PROXY
	// protocol carries it. When unset, PROXY protocol is enabled automatically on AWS shoots
	// (matching the retired nginx-ingress addon) and disabled on all other providers.
	ProxyProtocol *ProxyProtocolConfig `json:"proxyProtocol,omitempty"`
}

// ProxyProtocolConfig configures PROXY protocol handling for the Traefik LoadBalancer.
type ProxyProtocolConfig struct {
	// Enabled explicitly turns PROXY protocol on or off. When nil, PROXY protocol defaults
	// to enabled on AWS shoots and disabled on all other providers.
	Enabled *bool `json:"enabled,omitempty"`

	// TrustedIPs is the list of CIDRs (or bare IPs) Traefik trusts to send a PROXY protocol
	// header on the web and websecure entrypoints. When empty, it defaults to the shoot's
	// node and pod network CIDRs, since the LoadBalancer lives inside the shoot VPC.
	TrustedIPs []string `json:"trustedIPs,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TraefikConfig is the configuration schema for the Traefik extension.
// This extension deploys Traefik ingress controller to shoot clusters
// as a replacement for the nginx-ingress-controller which is out of maintenance.
//
// Deprecated: This version nests the configuration under a "spec" field. Use the
// v1alpha2 version, which places the configuration fields at the top level (no "spec"
// wrapper), following the convention of other Gardener extension providerConfig APIs.
// v1alpha1 is still accepted for backwards compatibility with existing shoots and is
// converted transparently.
type TraefikConfig struct {
	metav1.TypeMeta `json:",inline"`

	// Spec provides the Traefik extension configuration spec.
	Spec TraefikConfigSpec `json:"spec"`
}
