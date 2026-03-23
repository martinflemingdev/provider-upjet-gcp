// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package accessapproval

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_folder_access_approval_settings", func(r *config.Resource) {
		r.Kind = "FolderSettings"
		config.MarkAsRequired(r.TerraformResource, "folder_id")
		config.MarkAsRequired(r.TerraformResource, "enrolled_services")
	})
	p.AddResourceConfigurator("google_organization_access_approval_settings", func(r *config.Resource) {
		r.Kind = "OrganizationSettings"
		config.MarkAsRequired(r.TerraformResource, "organization_id")
		config.MarkAsRequired(r.TerraformResource, "enrolled_services")
	})
	p.AddResourceConfigurator("google_project_access_approval_settings", func(r *config.Resource) {
		r.Kind = "ProjectSettings"
		config.MarkAsRequired(r.TerraformResource, "project_id")
		config.MarkAsRequired(r.TerraformResource, "enrolled_services")
	})
}
