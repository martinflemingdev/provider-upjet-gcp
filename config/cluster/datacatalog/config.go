// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package datacatalog

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_data_catalog_entry_group", func(r *config.Resource) {
		r.MarkAsRequired("region")
	})

	p.AddResourceConfigurator("google_data_catalog_taxonomy_iam_member", func(r *config.Resource) {
		r.MarkAsRequired("member", "role")
		r.References["taxonomy"] = config.Reference{
			TerraformName: "google_data_catalog_taxonomy",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
	})

	p.AddResourceConfigurator("google_data_catalog_policy_tag_iam_member", func(r *config.Resource) {
		r.MarkAsRequired("member", "role")
		r.References["policy_tag"] = config.Reference{
			TerraformName: "google_data_catalog_policy_tag",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
	})
}
