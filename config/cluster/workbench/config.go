// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package workbench

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_workbench_instance", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
	})

	p.AddResourceConfigurator("google_workbench_instance_iam_member", func(r *config.Resource) {
		r.References["name"] = config.Reference{
			TerraformName: "google_workbench_instance",
		}
	})
}
