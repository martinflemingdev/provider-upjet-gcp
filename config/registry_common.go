// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	// Note(ezgidemirel): we are importing this to embed provider schema document
	_ "embed"

	"os"
	"strings"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	conversiontfjson "github.com/crossplane/upjet/v2/pkg/types/conversion/tfjson"
	"github.com/crossplane/upjet/v2/pkg/types/name"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

const (
	resourcePrefix = "gcp"
	modulePath     = "github.com/upbound/provider-gcp"
)

var (
	//go:embed schema.json
	providerSchema string

	//go:embed provider-metadata.yaml
	providerMetadata []byte

	// oldSingletonListAPIs is a newline-delimited list of Terraform resource
	// names with converted singleton list APIs with at least CRD API version
	// containing the old singleton list API. This is to prevent the API
	// conversion for the newly added resources whose CRD APIs will already
	// use embedded objects instead of the singleton lists and thus, will
	// not possess a CRD API version with the singleton list. Thus, for
	// the newly added resources (resources added after the singleton lists
	// have been converted), we do not need the CRD API conversion
	// functions that convert between singleton lists and embedded objects,
	// but we need only the Terraform conversion functions.
	// This list is immutable and represents the set of resources with the
	// already generated CRD API versions with now converted singleton lists.
	// Because new resources should never have singleton lists in their
	// generated APIs, there should be no need to add them to this list.
	// However, bugs might result in exceptions in the future.
	// Please see:
	// https://github.com/crossplane-contrib/provider-upjet-gcp/pull/508
	// for more context on singleton list to embedded object conversions.
	//go:embed old-singleton-list-apis.txt
	oldSingletonListAPIs string
)

var skipList = []string{
	// Note(turkenh): Following two resources conflicts their singular versions
	// "google_access_context_manager_access_level" and
	// "google_access_context_manager_service_perimeter". Skipping for now.
	"google_access_context_manager_access_levels$",
	"google_access_context_manager_service_perimeters$",
	// Note(piotr): Following resources are potentially dangerous to implement
	// details in: https://github.com/upbound/official-providers/issues/587
	"google_kms_crypto_key_iam_policy",
	"google_kms_crypto_key_iam_binding",
	"google_kms_key_ring_iam_policy",
	"google_kms_key_ring_iam_binding",
	"google_cloudfunctions_function_iam_policy",
	"google_cloudfunctions_function_iam_binding",
	"google_compute_region_disk_iam_policy",
	"google_compute_region_disk_iam_binding",
	// Note(donovamuller): Following resources are potentially dangerous to implement
	// details in: https://github.com/upbound/official-providers/issues/521
	"google_project_iam_policy",
	"google_project_iam_binding",
	"google_organization_iam_binding",
	"google_service_account_iam_policy",
	"google_service_account_iam_binding",
	"google_cloud_run_service_iam_policy",
	"google_cloud_run_service_iam_binding",
	"google_pubsub_topic_iam_policy",
	"google_pubsub_topic_iam_binding",
	"google_pubsub_subscription_iam_policy",
	"google_pubsub_subscription_iam_binding",
	"google_compute_disk_iam_policy",
	"google_compute_disk_iam_binding",
	"google_compute_instance_iam_policy",
	"google_compute_instance_iam_binding",
	"google_compute_image_iam_policy",
	"google_compute_image_iam_binding",
	"google_notebooks_instance_iam_policy",
	"google_notebooks_instance_iam_binding",
	"google_notebooks_runtime_iam_policy",
	"google_notebooks_runtime_iam_binding",
	"google_secret_manager_secret_iam_policy",
	"google_secret_manager_secret_iam_binding",
	"google_sourcerepo_repository_iam_policy",
	"google_sourcerepo_repository_iam_binding",
	"google_spanner_instance_iam_policy",
	"google_spanner_instance_iam_binding",
	"google_spanner_database_iam_policy",
	"google_spanner_database_iam_binding",
	"google_compute_subnetwork_iam_policy",
	"google_compute_subnetwork_iam_binding",
	"google_endpoints_service_iam_policy",
	"google_endpoints_service_iam_binding",
	"google_endpoints_service_consumers_iam_policy",
	"google_endpoints_service_consumers_iam_binding",
}

// workaround for the TF Google v4.77.0-based no-fork release: We would like to
// keep the types in the generated CRDs intact
// (prevent number->int type replacements).
func getProviderSchema(s string) (*schema.Provider, error) {
	ps := tfjson.ProviderSchemas{}
	if err := ps.UnmarshalJSON([]byte(s)); err != nil {
		panic(err)
	}
	if len(ps.Schemas) != 1 {
		return nil, errors.Errorf("there should exactly be 1 provider schema but there are %d", len(ps.Schemas))
	}
	var rs map[string]*tfjson.Schema
	for _, v := range ps.Schemas {
		rs = v.ResourceSchemas
		break
	}
	return &schema.Provider{
		ResourcesMap: conversiontfjson.GetV2ResourceMap(rs),
	}, nil
}

// resourceList returns the list of resources that have external
// name configured in the specified table.
func resourceList(t map[string]ujconfig.ExternalName) []string {
	l := make([]string, len(t))
	i := 0
	for n := range t {
		// Expected format is regex and we'd like to have exact matches.
		l[i] = n + "$"
		i++
	}
	return l
}

// filterByGroup returns a new map containing only resources that belong to
// the specified API group. The group match is performed by checking if the
// resource name starts with "google_<group>_" (e.g. "google_bigquery_"). This
// is appropriate for GCP family providers where resources are named with the
// pattern "google_<service>_<resource>".
func filterByGroup(t map[string]ujconfig.ExternalName, group string) map[string]ujconfig.ExternalName {
	if group == "" {
		return t
	}
	out := make(map[string]ujconfig.ExternalName)
	prefix := "google_" + group + "_"
	for n, v := range t {
		if strings.HasPrefix(n, prefix) {
			out[n] = v
		}
	}
	return out
}

// detectGroupFromBinary extracts the API group from the binary name or environment
// For family providers, checks multiple sources in order of preference:
// 1. PROVIDER_GROUP environment variable (explicit override)
// 2. HOSTNAME environment variable (set by Kubernetes, e.g., "provider-gcp-bigquery-...")
// 3. Binary name (fallback for local development)
// For the monolith provider, it returns empty string
// The context parameter is used for logging to identify whether this is called from cluster or namespaced provider
func detectGroupFromBinary(context string) string {
	// First, check for explicit PROVIDER_GROUP environment variable
	// This is the most reliable method for explicit configuration
	if group := os.Getenv("PROVIDER_GROUP"); group != "" {
		println("DEBUG: detectGroupFromBinary ("+context+"): using PROVIDER_GROUP env var =", group)
		return group
	}

	// Second, check HOSTNAME which Kubernetes sets for pods
	// Format: provider-gcp-<group>-<hash>-<pod-id>
	// Example: provider-gcp-monitoring-5af302a3f2a2-5fbcb9fc84-8s6s5
	if hostname := os.Getenv("HOSTNAME"); hostname != "" {
		println("DEBUG: detectGroupFromBinary ("+context+"): parsing HOSTNAME =", hostname)

		// Split by '-' and look for pattern: provider-gcp-<group>-...
		parts := strings.Split(hostname, "-")
		if len(parts) >= 3 && parts[0] == "provider" && parts[1] == "gcp" {
			group := parts[2]
			println("DEBUG: detectGroupFromBinary ("+context+"): extracted group from HOSTNAME =", group)
			return group
		}
		println("DEBUG: detectGroupFromBinary (" + context + "): HOSTNAME doesn't match expected pattern")
	}

	// Fallback to binary name detection (useful for local development)
	if len(os.Args) == 0 {
		println("DEBUG: detectGroupFromBinary (" + context + "): os.Args is empty, returning empty group")
		return ""
	}

	// Get the binary name without path
	binaryName := os.Getenv("BINARY_NAME")
	if binaryName == "" && len(os.Args) > 0 {
		binaryName = os.Args[0]
		// Strip path if present
		if lastSlash := strings.LastIndex(binaryName, "/"); lastSlash != -1 {
			binaryName = binaryName[lastSlash+1:]
		}
	}
	println("DEBUG: detectGroupFromBinary ("+context+"): binary name =", binaryName)

	// Common binary names that indicate monolith provider
	if binaryName == "provider" || binaryName == "monolith" || strings.HasPrefix(binaryName, "provider-gcp") {
		println("DEBUG: detectGroupFromBinary (" + context + "): detected monolith provider, returning empty group")
		return ""
	}

	// For family providers built locally, the binary name is the group name
	// Remove any common suffixes or prefixes
	group := binaryName
	group = strings.TrimPrefix(group, "provider-")
	group = strings.TrimPrefix(group, "gcp-")

	println("DEBUG: detectGroupFromBinary ("+context+"): detected family provider group from binary =", group)
	return group
}

func init() {
	// GCP specific acronyms

	// Todo(turkenh): move to Terrajet?
	name.AddAcronym("idp", "IdP")
	name.AddAcronym("oauth", "OAuth")
}
