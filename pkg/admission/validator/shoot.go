// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

// Package validator provides admission webhook validators for the Traefik extension.
package validator

import (
	"context"
	"fmt"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/gardener/gardener-extension-shoot-traefik/pkg/apis/config"
	"github.com/gardener/gardener-extension-shoot-traefik/pkg/traefik"
)

const (
	// Name is the name of the shoot validator webhook.
	Name = "shoot-validator"
	// ExtensionType is the type of extension being validated.
	ExtensionType = "shoot-traefik"
)

// shootValidator validates Shoot resources for the Traefik extension.
type shootValidator struct {
	client  client.Client
	decoder runtime.Decoder
}

// NewShootValidatorWebhook creates a new webhook for validating Shoot resources.
// It ensures that the Traefik extension can only be enabled for shoots with
// purpose "evaluation".
func NewShootValidatorWebhook(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	decoder := serializer.NewCodecFactory(mgr.GetScheme(), serializer.EnableStrict).UniversalDecoder()

	return extensionswebhook.New(mgr, extensionswebhook.Args{
		Name:   Name,
		Path:   "/webhooks/validate-shoot-traefik",
		Target: extensionswebhook.TargetSeed,
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewShootValidator(mgr.GetClient(), decoder): {
				{Obj: &gardencorev1beta1.Shoot{}},
			},
		},
	})
}

// NewShootValidator creates a new shoot validator.
func NewShootValidator(c client.Client, decoder runtime.Decoder) extensionswebhook.Validator {
	return &shootValidator{
		client:  c,
		decoder: decoder,
	}
}

// Validate validates the given object (Shoot) on create and update operations.
func (v *shootValidator) Validate(ctx context.Context, newClient, old client.Object) error {
	shoot, ok := newClient.(*gardencorev1beta1.Shoot)
	if !ok {
		return fmt.Errorf("expected *gardencorev1beta1.Shoot but got %T", newClient)
	}

	return v.validateShoot(shoot)
}

// validateShoot validates that if the Traefik extension is enabled, the shoot must have
// purpose "evaluation", the nginx-ingress addon must not be enabled, and the TraefikConfig
// must be valid.
func (v *shootValidator) validateShoot(shoot *gardencorev1beta1.Shoot) error {
	// Find the Traefik extension and check whether it is enabled.
	var traefikExtension *gardencorev1beta1.Extension
	for i, ext := range shoot.Spec.Extensions {
		if ext.Type == ExtensionType {
			if ext.Disabled != nil && *ext.Disabled {
				return nil
			}
			traefikExtension = &shoot.Spec.Extensions[i]

			break
		}
	}

	// If no Traefik extension, validation passes
	if traefikExtension == nil {
		return nil
	}

	// Validate that the shoot purpose is "evaluation"
	if shoot.Spec.Purpose == nil || *shoot.Spec.Purpose != gardencorev1beta1.ShootPurposeEvaluation {
		purposeStr := "nil"
		if shoot.Spec.Purpose != nil {
			purposeStr = string(*shoot.Spec.Purpose)
		}

		return fmt.Errorf(
			"traefik extension can only be enabled for shoots with purpose 'evaluation'. "+
				"Current purpose: %s. Traefik acts as a replacement for the nginx ingress controller "+
				"and is only supported for evaluation clusters",
			purposeStr,
		)
	}

	// The traefik extension is a replacement for the nginx-ingress addon, so both must not
	// be enabled at the same time.
	if shoot.Spec.Addons != nil && shoot.Spec.Addons.NginxIngress != nil && shoot.Spec.Addons.NginxIngress.Enabled { //nolint:staticcheck // SA1019: Spec.Addons is deprecated but still needs to be validated.
		return fmt.Errorf(
			"traefik extension cannot be enabled while the nginx-ingress addon " +
				"(spec.addons.nginxIngress.enabled) is enabled. Traefik acts as a replacement for the " +
				"nginx ingress controller; please disable the nginx-ingress addon first",
		)
	}

	return v.validateProviderConfig(traefikExtension)
}

// validateProviderConfig decodes and validates the TraefikConfig carried in the
// extension's providerConfig, if any is set.
func (v *shootValidator) validateProviderConfig(ext *gardencorev1beta1.Extension) error {
	if ext.ProviderConfig == nil {
		return nil
	}

	cfg := &config.TraefikConfig{}
	if err := runtime.DecodeInto(v.decoder, ext.ProviderConfig.Raw, cfg); err != nil {
		return fmt.Errorf("failed to decode traefik providerConfig: %w", err)
	}

	if cfg.LogLevel != "" && !traefik.ValidLogLevels.Has(cfg.LogLevel) {
		return fmt.Errorf("invalid traefik logLevel %q: must be one of %v", cfg.LogLevel, sets.List(traefik.ValidLogLevels))
	}

	if cfg.HTTPEntrypoint != "" && !traefik.ValidHTTPEntrypoints.Has(cfg.HTTPEntrypoint) {
		return fmt.Errorf("invalid traefik httpEntrypoint %q: must be one of %v", cfg.HTTPEntrypoint, sets.List(traefik.ValidHTTPEntrypoints))
	}

	return nil
}
