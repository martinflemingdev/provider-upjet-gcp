// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package logging

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/upbound/provider-gcp/v2/config/namespaced/common"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_logging_project_bucket_config", func(r *config.Resource) {
		r.References["cmek_settings.kms_key_name"] = config.Reference{
			TerraformName: "google_kms_crypto_key",
			Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
		}
	})
	p.AddResourceConfigurator("google_logging_project_sink", func(r *config.Resource) {
		r.References["custom_writer_identity"] = config.Reference{
			TerraformName: "google_service_account",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("member",true)`,
		}
		r.References["destination"] = config.Reference{
			TerraformName: "google_logging_project_bucket_config",
			Extractor:     common.ExtractLoggingBucketDestinationFuncPath,
		}
	})

	p.AddResourceConfigurator("google_logging_folder_bucket_config", func(r *config.Resource) {
		r.References["folder"] = config.Reference{
			TerraformName: "google_folder",
			Extractor:     common.ExtractFolderIDFuncPath,
		}
	})
	p.AddResourceConfigurator("google_logging_project_bucket_config", func(r *config.Resource) {
		originalRead := r.TerraformResource.Read
		r.TerraformResource.Read = func(data *schema.ResourceData, meta interface{}) error {
			err := originalRead(data, meta)
			if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
				data.SetId("")
				return nil
			}
			return err
		}
	})
}
