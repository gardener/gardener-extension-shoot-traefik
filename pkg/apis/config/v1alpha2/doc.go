// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// +k8s:deepcopy-gen=package
// +k8s:defaulter-gen=TypeMeta
// +k8s:conversion-gen=github.com/gardener/gardener-extension-shoot-traefik/pkg/apis/config
// +groupName=traefik.extensions.gardener.cloud

// Package v1alpha2 provides the v1alpha2 version of the external API types.
// Unlike v1alpha1, the configuration fields live at the top level with no spec wrapper,
// following the convention of other Gardener extension providerConfig APIs.
package v1alpha2
