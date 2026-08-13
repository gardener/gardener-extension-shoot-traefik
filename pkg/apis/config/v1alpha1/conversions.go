// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/conversion"

	"github.com/gardener/gardener-extension-shoot-traefik/pkg/apis/config"
)

// Convert_v1alpha1_TraefikConfig_To_config_TraefikConfig converts the deprecated,
// spec-wrapped v1alpha1 TraefikConfig into the flat internal TraefikConfig.
//
// This conversion is written by hand because the two versions have different shapes:
// v1alpha1 nests the settings under a "spec" field, while the internal type keeps them
// at the top level. conversion-gen cannot derive this automatically, so it detects this
// manual function and skips generating one.
//
//nolint:revive // Convert_<from>_To_<to> is the naming convention required by conversion-gen.
func Convert_v1alpha1_TraefikConfig_To_config_TraefikConfig(in *TraefikConfig, out *config.TraefikConfig, _ conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.Replicas = in.Spec.Replicas
	out.IngressProvider = config.IngressProviderType(in.Spec.IngressProvider)
	out.LogLevel = in.Spec.LogLevel
	out.Dashboard = in.Spec.Dashboard
	out.HTTPEntrypoint = config.HTTPEntrypointType(in.Spec.HTTPEntrypoint)

	return nil
}

// Convert_config_TraefikConfig_To_v1alpha1_TraefikConfig converts the flat internal
// TraefikConfig into the deprecated, spec-wrapped v1alpha1 TraefikConfig.
//
//nolint:revive // Convert_<from>_To_<to> is the naming convention required by conversion-gen.
func Convert_config_TraefikConfig_To_v1alpha1_TraefikConfig(in *config.TraefikConfig, out *TraefikConfig, _ conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.Spec.Replicas = in.Replicas
	out.Spec.IngressProvider = IngressProviderType(in.IngressProvider)
	out.Spec.LogLevel = in.LogLevel
	out.Spec.Dashboard = in.Dashboard
	out.Spec.HTTPEntrypoint = HTTPEntrypointType(in.HTTPEntrypoint)

	return nil
}
