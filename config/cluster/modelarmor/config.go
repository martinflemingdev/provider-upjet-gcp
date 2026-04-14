// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package modelarmor

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_model_armor_floorsetting", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "parent") // expects projects/{project}, folders/{folder} or organizations/{organization}
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "filter_config")
	})

	p.AddResourceConfigurator("google_model_armor_template", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "filter_config")
	})
}
