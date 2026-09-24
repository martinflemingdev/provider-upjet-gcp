// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	repository "github.com/upbound/provider-gcp/v3/internal/controller/cluster/dataform/repository"
)

// Setup_dataform creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_dataform(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		repository.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_dataform creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_dataform(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		repository.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager_dataform registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager_dataform(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		repository.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
