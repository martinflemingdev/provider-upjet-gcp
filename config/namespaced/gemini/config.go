// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package gemini

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_gemini_code_repository_index", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
	})

	p.AddResourceConfigurator("google_gemini_code_tools_setting", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "enabled_tool")
	})

	p.AddResourceConfigurator("google_gemini_code_tools_setting_binding", func(r *config.Resource) {
		r.References["code_tools_setting_id"] = config.Reference{
			TerraformName: "google_gemini_code_tools_setting",
		}
		config.MarkAsRequired(r.TerraformResource, "target")
	})

	p.AddResourceConfigurator("google_gemini_data_sharing_with_google_setting", func(r *config.Resource) {
	})

	p.AddResourceConfigurator("google_gemini_data_sharing_with_google_setting_binding", func(r *config.Resource) {
		r.References["data_sharing_with_google_setting_id"] = config.Reference{
			TerraformName: "google_gemini_data_sharing_with_google_setting",
		}
		config.MarkAsRequired(r.TerraformResource, "target")
	})

	p.AddResourceConfigurator("google_gemini_gemini_gcp_enablement_setting", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
	})

	p.AddResourceConfigurator("google_gemini_gemini_gcp_enablement_setting_binding", func(r *config.Resource) {
		r.References["gemini_gcp_enablement_setting_id"] = config.Reference{
			TerraformName: "google_gemini_gemini_gcp_enablement_setting",
		}
		config.MarkAsRequired(r.TerraformResource, "target")
	})

	p.AddResourceConfigurator("google_gemini_logging_setting", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
	})

	p.AddResourceConfigurator("google_gemini_logging_setting_binding", func(r *config.Resource) {
		r.References["logging_setting_id"] = config.Reference{
			TerraformName: "google_gemini_logging_setting",
		}
		config.MarkAsRequired(r.TerraformResource, "target")
	})

	p.AddResourceConfigurator("google_gemini_release_channel_setting", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
	})

	p.AddResourceConfigurator("google_gemini_release_channel_setting_binding", func(r *config.Resource) {
		r.References["release_channel_setting_id"] = config.Reference{
			TerraformName: "google_gemini_release_channel_setting",
		}
		config.MarkAsRequired(r.TerraformResource, "target")
	})

	p.AddResourceConfigurator("google_gemini_repository_group", func(r *config.Resource) {
		r.References["code_repository_index"] = config.Reference{
			TerraformName: "google_gemini_code_repository_index",
		}
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "repositories")
	})

	p.AddResourceConfigurator("google_gemini_repository_group_iam_member", func(r *config.Resource) {
		r.References["repository_group_id"] = config.Reference{
			TerraformName: "google_gemini_repository_group",
		}
		r.References["code_repository_index"] = config.Reference{
			TerraformName: "google_gemini_code_repository_index",
		}
	})
}
