// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package cloudfunctions2

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_cloudfunctions2_function_iam_member", func(r *config.Resource) {
		r.References["cloud_function"] = config.Reference{
			TerraformName: "google_cloudfunctions2_function",
		}
	})
}
