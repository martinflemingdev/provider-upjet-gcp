package billing

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures billing resources.
func Configure(p *config.Provider) {
	// Budget
	p.AddResourceConfigurator("google_billing_budget", func(r *config.Resource) {
		// Required arguments
		config.MarkAsRequired(r.TerraformResource, "amount")
		config.MarkAsRequired(r.TerraformResource, "billing_account")
	})

	// ProjectInfo
	p.AddResourceConfigurator("google_billing_project_info", func(r *config.Resource) {
		// Required arguments
		config.MarkAsRequired(r.TerraformResource, "billing_account")
	})

	// Billing Account IAM Member
	p.AddResourceConfigurator("google_billing_account_iam_member", func(r *config.Resource) {
		// Required arguments
		config.MarkAsRequired(r.TerraformResource, "billing_account_id")
		config.MarkAsRequired(r.TerraformResource, "role")
		config.MarkAsRequired(r.TerraformResource, "member")
	})
}
