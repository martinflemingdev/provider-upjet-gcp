// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/upbound/provider-gcp/config/cluster/common"
)

// terraformPluginSDKExternalNameConfigs contains all external name configurations
// belonging to Terraform resources to be reconciled under the no-fork
// architecture for this provider.
var terraformPluginSDKExternalNameConfigs = map[string]config.ExternalName{
	// vertexai
	//
	// No Import
	"google_vertex_ai_dataset": config.IdentifierFromProvider,
	// Imported by using the following projects/{{project}}/locations/{{region}}/featurestores/{{name}}
	"google_vertex_ai_featurestore": config.IdentifierFromProvider,
	// Imported by using the following {{featurestore}}/entityTypes/{{name}}
	"google_vertex_ai_featurestore_entitytype": config.IdentifierFromProvider,
	// Imported by using the following projects/{{project}}/locations/{{region}}/tensorboards/{{name}}
	"google_vertex_ai_tensorboard": config.TemplatedStringAsIdentifier("display_name", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/tensorboards/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{region}}/deploymentResourcePools/{{name}}
	"google_vertex_ai_deployment_resource_pool": config.TemplatedStringAsIdentifier("name", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/deploymentResourcePools/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{location}}/endpoints/{{name}}
	"google_vertex_ai_endpoint": config.TemplatedStringAsIdentifier("name", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.location }}/endpoints/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{location}}/endpoints/{{endpoint}} roles/viewer user:jane@example.com
	"google_vertex_ai_endpoint_iam_member": config.IdentifierFromProvider,
	// No Import
	"google_vertex_ai_endpoint_with_model_garden_deployment": config.IdentifierFromProvider,
	// Imported by using the following projects/{{project}}/locations/{{region}}/indexes/{{name}}
	"google_vertex_ai_index": config.TemplatedStringAsIdentifier("name", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/indexes/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{region}}/indexEndpoints/{{name}}
	"google_vertex_ai_index_endpoint": config.TemplatedStringAsIdentifier("name", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/indexEndpoints/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{region}}/indexEndpoints/{{index_endpoint}}/deployedIndex/{{deployed_index_id}}
	"google_vertex_ai_index_endpoint_deployed_index": config.TemplatedStringAsIdentifier("deployed_index_id", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/indexEndpoints/{{ .parameters.index_endpoint }}/deployedIndex/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{region}}/metadataStores/{{name}}
	"google_vertex_ai_metadata_store": config.TemplatedStringAsIdentifier("name", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/metadataStores/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{region}}/ragEngineConfig
	"google_vertex_ai_rag_engine_config": config.TemplatedStringAsIdentifier("", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/ragEngineConfig"),
	// Imported by using the following projects/{{project}}/locations/{{region}}/featureGroups/{{name}}
	"google_vertex_ai_feature_group": config.TemplatedStringAsIdentifier("name", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/featureGroups/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{region}}/featureGroups/{{feature_group}}/features/{{name}}
	"google_vertex_ai_feature_group_feature": config.TemplatedStringAsIdentifier("name", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/featureGroups/{{ .parameters.feature_group }}/features/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{region}}/featureGroups/{{feature_group}} roles/viewer user:jane@example.com
	"google_vertex_ai_feature_group_iam_member": config.IdentifierFromProvider,
	// Imported by using the following projects/{{project}}/locations/{{region}}/featureOnlineStores/{{name}}
	"google_vertex_ai_feature_online_store": config.TemplatedStringAsIdentifier("name", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/featureOnlineStores/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{region}}/featureOnlineStores/{{feature_online_store}}/featureViews/{{name}}
	"google_vertex_ai_feature_online_store_featureview": config.TemplatedStringAsIdentifier("name", "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.region }}/featureOnlineStores/{{ .parameters.feature_online_store }}/featureViews/{{ .external_name }}"),
	// Imported by using the following projects/{{project}}/locations/{{region}}/featureOnlineStores/{{feature_online_store}}/featureViews/{{featureview}} roles/viewer user:jane@example.com
	"google_vertex_ai_feature_online_store_featureview_iam_member": config.IdentifierFromProvider,
	// Imported by using the following projects/{{project}}/locations/{{region}}/featureOnlineStores/{{feature_online_store}} roles/viewer user:jane@example.com
	"google_vertex_ai_feature_online_store_iam_member": config.IdentifierFromProvider,
}

// cliReconciledExternalNameConfigs contains all external name configurations
// belonging to Terraform resources to be reconciled under the CLI-based
// architecture for this provider.
var cliReconciledExternalNameConfigs = map[string]config.ExternalName{}

// TemplatedStringAsIdentifierWithNoName uses TemplatedStringAsIdentifier but
// without the name initializer. This allows it to be used in cases where the ID
// is constructed with parameters and a provider-defined value, meaning no
// user-defined input. Since the external name is not user-defined, the name
// initializer has to be disabled.
func TemplatedStringAsIdentifierWithNoName(tmpl string) config.ExternalName {
	e := config.TemplatedStringAsIdentifier("", tmpl)
	e.DisableNameInitializer = true
	return e
}

// resourceConfigurator applies all external name configs
// listed in the table terraformPluginSDKExternalNameConfigs and
// cliReconciledExternalNameConfigs and sets the version
// of those resources to v1beta1. For those resource in
// terraformPluginSDKExternalNameConfigs, it also sets
// config.Resource.UseNoForkClient to `true`.
func resourceConfigurator() config.ResourceOption {
	return func(r *config.Resource) {
		// if configured both for the no-fork and CLI based architectures,
		// no-fork configuration prevails
		e, configured := terraformPluginSDKExternalNameConfigs[r.Name]
		if !configured {
			e, configured = cliReconciledExternalNameConfigs[r.Name]
		}
		if !configured {
			return
		}
		r.Version = common.VersionV1Beta1
		r.ExternalName = e
	}
}
