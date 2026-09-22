// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package networkservices

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/upbound/provider-gcp/v3/config/cluster/common"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_network_services_agent_gateway", func(r *config.Resource) {
		r.MarkAsRequired("location")
		r.References["agent_connectivity_template"] = config.Reference{
			TerraformName: "google_network_services_agent_connectivity_template",
			Extractor:     common.ExtractResourceIDFuncPath,
		}
		r.References["network_config.dns_peering_config.target_network"] = config.Reference{
			TerraformName: "google_compute_network",
			Extractor:     common.ExtractResourceIDFuncPath,
		}
	})
}
