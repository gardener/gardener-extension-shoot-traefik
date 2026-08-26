// SPDX-FileCopyrightText: Copyright Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package traefik

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/gardener/gardener-extension-shoot-traefik/pkg/apis/config"
)

const (
	// Namespace is the namespace where Traefik will be deployed in the shoot cluster.
	Namespace = "kube-system"

	// DeploymentName is the name of the Traefik deployment.
	DeploymentName = "traefik"

	// ServiceAccountName is the name of the service account for Traefik.
	ServiceAccountName = "traefik"

	// ManagedResourceName is the name of the ManagedResource for Traefik.
	ManagedResourceName = "extension-traefik"

	// ImageName is the name of the Traefik image in the image vector.
	ImageName = "traefik"

	// SeedManagedResourceName is the name of the seed-class ManagedResource
	// that contains the DNSRecord for the Traefik ingress wildcard domain.
	SeedManagedResourceName = "extension-traefik-ingress-dns"

	// LogLevelInfo is the default Traefik log level.
	LogLevelInfo = "Info"

	// LabelName is the "app.kubernetes.io/name" label key.
	LabelName = "app.kubernetes.io/name"

	// LabelInstance is the "app.kubernetes.io/instance" label key.
	LabelInstance = "app.kubernetes.io/instance"

	// LabelComponent is the "app.kubernetes.io/component" label key.
	LabelComponent = "app.kubernetes.io/component"

	// LabelManagedBy is the "app.kubernetes.io/managed-by" label key.
	LabelManagedBy = "app.kubernetes.io/managed-by"

	// LabelComponentValue is the value for the ingress-controller component label.
	LabelComponentValue = "ingress-controller"

	// LabelManagedByValue is the value for the gardener managed-by label.
	LabelManagedByValue = "gardener"

	// IngressControllerRoleName is the name of the ClusterRole and ClusterRoleBinding for Traefik.
	IngressControllerRoleName = "traefik-ingress-controller"

	// PingPath is the path used for Traefik health probes.
	PingPath = "/ping"

	// VerbGet is the "get" RBAC verb.
	VerbGet = "get"

	// VerbList is the "list" RBAC verb.
	VerbList = "list"

	// VerbWatch is the "watch" RBAC verb.
	VerbWatch = "watch"

	// AnnotationIsDefaultClass is the annotation that marks a class as default.
	AnnotationIsDefaultClass = "ingressclass.kubernetes.io/is-default-class"

	// AnnotationIsDefaultClassValue is the value for the is-default-class annotation.
	AnnotationIsDefaultClassValue = "true"

	// TraefikIngressController is the controller value for the Traefik ingress class.
	TraefikIngressController = "traefik.io/ingress-controller"

	// ConfirmationDeletionValue is the value for the deletion confirmation annotation.
	ConfirmationDeletionValue = "true"

	// IngressClassNGINX is the ingress class name for nginx-compatible mode.
	IngressClassNGINX = "nginx"

	// EntrypointWeb is the name of the plain-HTTP entrypoint and its Service/container port.
	EntrypointWeb = "web"

	// EntrypointWebSecure is the name of the HTTPS entrypoint and its Service/container port.
	EntrypointWebSecure = "websecure"

	// ArgEntrypointWebAddress is the Traefik argument configuring the web entrypoint address.
	ArgEntrypointWebAddress = "--entrypoints.web.address=:8000"
)

// ValidLogLevels contains the set of log levels supported by Traefik.
var ValidLogLevels = sets.New(
	"Debug",
	"Info",
	"Warn",
	"Error",
	"Fatal",
	"Panic",
)

// ValidHTTPEntrypoints contains the set of HTTP entrypoint modes supported by the extension.
var ValidHTTPEntrypoints = sets.New(
	config.HTTPEntrypointEnabled,
	config.HTTPEntrypointRedirect,
	config.HTTPEntrypointDisabled,
)
