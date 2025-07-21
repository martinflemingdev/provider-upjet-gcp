// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package vertexai

import (
	"github.com/crossplane/upjet/pkg/config"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_vertex_ai_index", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "region")
	})
	p.AddResourceConfigurator("google_vertex_ai_tensorboard", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "region")
	})
	p.AddResourceConfigurator("google_vertex_ai_endpoint_iam_member", func(r *config.Resource) {
		r.References["endpoint"] = config.Reference{
			TerraformName: "google_vertex_ai_endpoint",
		}
	})
	p.AddResourceConfigurator("google_workbench_instance_iam_member", func(r *config.Resource) {
		r.References["name"] = config.Reference{
			TerraformName: "google_workbench_instance",
		}
	})
}
