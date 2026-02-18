// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package vertexai

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_vertex_ai_deployment_resource_pool", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_endpoint", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
	})

	p.AddResourceConfigurator("google_vertex_ai_endpoint_iam_member", func(r *config.Resource) {
		r.References["endpoint"] = config.Reference{
			TerraformName: "google_vertex_ai_endpoint",
		}
	})

	p.AddResourceConfigurator("google_vertex_ai_endpoint_with_model_garden_deployment", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "location")
	})

	p.AddResourceConfigurator("google_vertex_ai_index", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_index_endpoint", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_index_endpoint_deployed_index", func(r *config.Resource) {
		r.References["index_endpoint"] = config.Reference{
			TerraformName: "google_vertex_ai_index_endpoint",
		}
		r.References["index"] = config.Reference{
			TerraformName: "google_vertex_ai_index",
		}
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_metadata_store", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_rag_engine_config", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_tensorboard", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_feature_group", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_feature_group_feature", func(r *config.Resource) {
		r.References["feature_group"] = config.Reference{
			TerraformName: "google_vertex_ai_feature_group",
		}
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_feature_group_iam_member", func(r *config.Resource) {
		r.References["feature_group"] = config.Reference{
			TerraformName: "google_vertex_ai_feature_group",
		}
	})

	p.AddResourceConfigurator("google_vertex_ai_feature_online_store", func(r *config.Resource) {
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_feature_online_store_featureview", func(r *config.Resource) {
		r.References["feature_online_store"] = config.Reference{
			TerraformName: "google_vertex_ai_feature_online_store",
		}
		config.MarkAsRequired(r.TerraformResource, "region")
	})

	p.AddResourceConfigurator("google_vertex_ai_feature_online_store_featureview_iam_member", func(r *config.Resource) {
		r.References["feature_online_store"] = config.Reference{
			TerraformName: "google_vertex_ai_feature_online_store",
		}
		r.References["feature_view"] = config.Reference{
			TerraformName: "google_vertex_ai_feature_online_store_featureview",
		}
	})

	p.AddResourceConfigurator("google_vertex_ai_feature_online_store_iam_member", func(r *config.Resource) {
		r.References["feature_online_store"] = config.Reference{
			TerraformName: "google_vertex_ai_feature_online_store",
		}
	})
}