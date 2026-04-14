// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package discoveryengine

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_discovery_engine_acl_config", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
	})

	p.AddResourceConfigurator("google_discovery_engine_assistant", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "collection_id")
		config.MarkAsRequired(r.TerraformResource, "display_name")
		r.References["engine_id"] = config.Reference{
			TerraformName: "google_discovery_engine_search_engine",
		}
	})

	p.AddResourceConfigurator("google_discovery_engine_chat_engine", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "collection_id")
		config.MarkAsRequired(r.TerraformResource, "display_name")
		config.MarkAsRequired(r.TerraformResource, "data_store_ids")
		config.MarkAsRequired(r.TerraformResource, "chat_engine_config")
	})

	p.AddResourceConfigurator("google_discovery_engine_cmek_config", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "kms_key")
	})

	p.AddResourceConfigurator("google_discovery_engine_control", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "display_name")
		config.MarkAsRequired(r.TerraformResource, "solution_type")
		r.References["engine_id"] = config.Reference{
			TerraformName: "google_discovery_engine_search_engine",
		}
	})

	p.AddResourceConfigurator("google_discovery_engine_data_connector", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "data_source")
		config.MarkAsRequired(r.TerraformResource, "refresh_interval")
		config.MarkAsRequired(r.TerraformResource, "collection_display_name")
	})

	p.AddResourceConfigurator("google_discovery_engine_data_store", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "display_name")
		config.MarkAsRequired(r.TerraformResource, "industry_vertical")
		config.MarkAsRequired(r.TerraformResource, "content_config")
	})

	p.AddResourceConfigurator("google_discovery_engine_license_config", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "license_count")
		config.MarkAsRequired(r.TerraformResource, "subscription_tier")
		config.MarkAsRequired(r.TerraformResource, "start_date")
		config.MarkAsRequired(r.TerraformResource, "subscription_term")
	})

	p.AddResourceConfigurator("google_discovery_engine_recommendation_engine", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "display_name")
		config.MarkAsRequired(r.TerraformResource, "data_store_ids")
	})

	p.AddResourceConfigurator("google_discovery_engine_schema", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		r.References["data_store_id"] = config.Reference{
			TerraformName: "google_discovery_engine_data_store",
		}
	})

	p.AddResourceConfigurator("google_discovery_engine_search_engine", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "collection_id")
		config.MarkAsRequired(r.TerraformResource, "display_name")
		config.MarkAsRequired(r.TerraformResource, "data_store_ids")
		config.MarkAsRequired(r.TerraformResource, "search_engine_config")
	})

	p.AddResourceConfigurator("google_discovery_engine_serving_config", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		r.References["engine_id"] = config.Reference{
			TerraformName: "google_discovery_engine_search_engine",
		}
	})

	p.AddResourceConfigurator("google_discovery_engine_sitemap", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		r.References["data_store_id"] = config.Reference{
			TerraformName: "google_discovery_engine_data_store",
		}
	})

	p.AddResourceConfigurator("google_discovery_engine_target_site", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		config.MarkAsRequired(r.TerraformResource, "provided_uri_pattern")
		r.References["data_store_id"] = config.Reference{
			TerraformName: "google_discovery_engine_data_store",
		}
	})

	p.AddResourceConfigurator("google_discovery_engine_user_store", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
	})

	p.AddResourceConfigurator("google_discovery_engine_widget_config", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
		r.References["engine_id"] = config.Reference{
			TerraformName: "google_discovery_engine_search_engine",
		}
	})
}
