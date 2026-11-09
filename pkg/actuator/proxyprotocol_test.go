// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package actuator

import (
	"slices"
	"testing"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"

	"github.com/gardener/gardener-extension-shoot-traefik/pkg/apis/config"
)

func shoot(providerType string, nodes, pods *string) *gardencorev1beta1.Shoot {
	s := &gardencorev1beta1.Shoot{}
	s.Spec.Provider.Type = providerType
	if nodes != nil || pods != nil {
		s.Spec.Networking = &gardencorev1beta1.Networking{Nodes: nodes, Pods: pods}
	}

	return s
}

func TestResolveProxyProtocol(t *testing.T) {
	nodes := "10.250.0.0/16"
	pods := "100.96.0.0/11"

	tests := []struct {
		name           string
		cfg            *config.ProxyProtocolConfig
		shoot          *gardencorev1beta1.Shoot
		wantEnabled    bool
		wantTrustedIPs []string
	}{
		{
			name:           "AWS, no config, derives node+pod CIDRs",
			cfg:            nil,
			shoot:          shoot("aws", &nodes, &pods),
			wantEnabled:    true,
			wantTrustedIPs: []string{nodes, pods},
		},
		{
			name:           "non-AWS, no config, disabled",
			cfg:            nil,
			shoot:          shoot("gcp", &nodes, &pods),
			wantEnabled:    false,
			wantTrustedIPs: nil,
		},
		{
			name:           "AWS, explicitly disabled, override wins",
			cfg:            &config.ProxyProtocolConfig{Enabled: new(false)},
			shoot:          shoot("aws", &nodes, &pods),
			wantEnabled:    false,
			wantTrustedIPs: nil,
		},
		{
			name:           "GCP, explicitly enabled, override wins and derives CIDRs",
			cfg:            &config.ProxyProtocolConfig{Enabled: new(true)},
			shoot:          shoot("gcp", &nodes, &pods),
			wantEnabled:    true,
			wantTrustedIPs: []string{nodes, pods},
		},
		{
			name:           "explicit trustedIPs used verbatim, no derivation",
			cfg:            &config.ProxyProtocolConfig{Enabled: new(true), TrustedIPs: []string{"192.0.2.0/24"}},
			shoot:          shoot("aws", &nodes, &pods),
			wantEnabled:    true,
			wantTrustedIPs: []string{"192.0.2.0/24"},
		},
		{
			name:           "AWS, no networking, enabled with no trusted IPs",
			cfg:            nil,
			shoot:          shoot("aws", nil, nil),
			wantEnabled:    true,
			wantTrustedIPs: nil,
		},
		{
			name:           "AWS, only nodes CIDR set",
			cfg:            nil,
			shoot:          shoot("aws", &nodes, nil),
			wantEnabled:    true,
			wantTrustedIPs: []string{nodes},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEnabled, gotTrustedIPs := resolveProxyProtocol(tt.cfg, tt.shoot)
			if gotEnabled != tt.wantEnabled {
				t.Errorf("enabled: got %v, want %v", gotEnabled, tt.wantEnabled)
			}
			if !slices.Equal(gotTrustedIPs, tt.wantTrustedIPs) {
				t.Errorf("trustedIPs: got %v, want %v", gotTrustedIPs, tt.wantTrustedIPs)
			}
		})
	}
}
