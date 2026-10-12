// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package bigquery

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	// These resources gained Terraform deletion_policy in v7.33.0. Keep it
	// observable without adding a second deletion control to the MR spec.
	for _, resourceName := range []string{
		"google_bigquery_connection",
		"google_bigquery_data_transfer_config",
		"google_bigquery_dataset",
		"google_bigquery_dataset_access",
		"google_bigquery_reservation",
		"google_bigquery_reservation_assignment",
		"google_bigquery_routine",
		"google_bigquery_table",
		"google_bigquery_datapolicy_data_policy",
		"google_bigquery_datapolicyv2_data_policy",
	} {
		p.AddResourceConfigurator(resourceName, func(r *config.Resource) {
			config.MoveToStatus(r.TerraformResource, "deletion_policy")
		})
	}

	p.AddResourceConfigurator("google_bigquery_datapolicy_data_policy", func(r *config.Resource) {
		r.Kind = "DataPolicy"
		r.References["policy_tag"] = config.Reference{
			TerraformName: "google_data_catalog_policy_tag",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
		r.References["data_masking_policy.routine"] = config.Reference{
			TerraformName: "google_bigquery_routine",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
	})
	p.AddResourceConfigurator("google_bigquery_datapolicy_data_policy_iam_member", func(r *config.Resource) {
		r.Kind = "DataPolicyIAMMember"
		// The IAM resource accepts the full policy ID and derives project and
		// location from it, so references also work across projects/locations.
		r.References["data_policy_id"] = config.Reference{
			TerraformName: "google_bigquery_datapolicy_data_policy",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
	})
	p.AddResourceConfigurator("google_bigquery_datapolicyv2_data_policy", func(r *config.Resource) {
		r.Kind = "DataPolicyV2"
		// BigQuery expects a namespaced key (project-or-org/key) and a short
		// value, rather than the numeric tagKeys/... and tagValues/... IDs.
		r.References["data_governance_tag.key"] = config.Reference{
			TerraformName: "google_tags_tag_key",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("namespaced_name",true)`,
		}
		r.References["data_governance_tag.value"] = config.Reference{
			TerraformName: "google_tags_tag_value",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("short_name",true)`,
		}
		r.References["data_masking_policy.routine"] = config.Reference{
			TerraformName: "google_bigquery_routine",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
	})
	p.AddResourceConfigurator("google_bigquery_datapolicyv2_data_policy_iam_member", func(r *config.Resource) {
		r.Kind = "DataPolicyV2IAMMember"
		r.References["data_policy_id"] = config.Reference{
			TerraformName: "google_bigquery_datapolicyv2_data_policy",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
	})
	p.AddResourceConfigurator("google_bigquery_dataset", func(r *config.Resource) {
		delete(r.References, "access.dataset.dataset.project_id")
	})
	p.AddResourceConfigurator("google_bigquery_dataset_access", func(r *config.Resource) {
		delete(r.References, "view.project_id")
		delete(r.References, "dataset.dataset.project_id")
	})
	p.AddResourceConfigurator("google_bigquery_dataset_iam_member", func(r *config.Resource) {
		r.References["dataset_id"] = config.Reference{
			TerraformName: "google_bigquery_dataset",
		}
	})
	p.AddResourceConfigurator("google_bigquery_dataset_iam_binding", func(r *config.Resource) {
		r.References["dataset_id"] = config.Reference{
			TerraformName: "google_bigquery_dataset",
		}
	})
	p.AddResourceConfigurator("google_bigquery_table_iam_binding", func(r *config.Resource) {
		r.References["dataset_id"] = config.Reference{
			TerraformName: "google_bigquery_dataset",
		}
		r.References["table_id"] = config.Reference{
			TerraformName: "google_bigquery_table",
		}
	})
	p.AddResourceConfigurator("google_bigquery_table_iam_policy", func(r *config.Resource) {
		r.References["dataset_id"] = config.Reference{
			TerraformName: "google_bigquery_dataset",
		}
		delete(r.References, "project")
	})
	p.AddResourceConfigurator("google_bigquery_table_iam_member", func(r *config.Resource) {
		r.References["dataset_id"] = config.Reference{
			TerraformName: "google_bigquery_dataset",
		}
		r.References["table_id"] = config.Reference{
			TerraformName: "google_bigquery_table",
		}
	})
	p.AddResourceConfigurator("google_bigquery_job", func(r *config.Resource) {
		r.References["query.destination_table.dataset_id"] = config.Reference{
			TerraformName: "google_bigquery_dataset",
		}
		r.References["extract.source_table.table_id"] = config.Reference{
			TerraformName: "google_bigquery_table",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
		r.References["load.destination_table.dataset_id"] = config.Reference{
			TerraformName: "google_bigquery_dataset",
		}
		r.References["extract.source_table.dataset_id"] = config.Reference{
			TerraformName: "google_bigquery_dataset",
		}
		r.References["query.destination_table.table_id"] = config.Reference{
			TerraformName: "google_bigquery_table",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
		r.References["load.destination_table.table_id"] = config.Reference{
			TerraformName: "google_bigquery_table",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
		r.References["copy.destination_table.dataset_id"] = config.Reference{
			TerraformName: "google_bigquery_dataset",
		}
		r.References["copy.destination_table.table_id"] = config.Reference{
			TerraformName: "google_bigquery_table",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}

		delete(r.References, "copy.destination_table.project_id")
		delete(r.References, "query.destination_table.project_id")
		delete(r.References, "load.destination_table.project_id")
		delete(r.References, "extract.source_table.project_id")
	})
	p.AddResourceConfigurator("google_bigquery_analytics_hub_data_exchange_iam_member", func(r *config.Resource) {
		r.References["data_exchange_id"] = config.Reference{
			TerraformName: "google_bigquery_analytics_hub_data_exchange",
		}
	})
	p.AddResourceConfigurator("google_bigquery_analytics_hub_data_exchange", func(r *config.Resource) {
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"sharing_environment_config"},
		}
		r.TerraformCustomDiff = func(diff *terraform.InstanceDiff, _ *terraform.InstanceState, _ *terraform.ResourceConfig) (*terraform.InstanceDiff, error) {
			if diff == nil || diff.Empty() || diff.Destroy || diff.Attributes == nil {
				return diff, nil
			}
			delete(diff.Attributes, "sharing_environment_config.0.default_exchange_config.#")
			return diff, nil
		}
	})
	p.AddResourceConfigurator("google_bigquery_analytics_hub_listing_iam_member", func(r *config.Resource) {
		r.References["listing_id"] = config.Reference{
			TerraformName: "google_bigquery_analytics_hub_listing",
		}
		r.References["data_exchange_id"] = config.Reference{
			TerraformName: "google_bigquery_analytics_hub_data_exchange",
			// Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
	})
}
