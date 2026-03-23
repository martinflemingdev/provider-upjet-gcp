// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	foldersettings "github.com/upbound/provider-gcp/internal/controller/namespaced/accessapproval/foldersettings"
	organizationsettings "github.com/upbound/provider-gcp/internal/controller/namespaced/accessapproval/organizationsettings"
	projectsettings "github.com/upbound/provider-gcp/internal/controller/namespaced/accessapproval/projectsettings"
)

// Setup_accessapproval creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_accessapproval(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		foldersettings.Setup,
		organizationsettings.Setup,
		projectsettings.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_accessapproval creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_accessapproval(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		foldersettings.SetupGated,
		organizationsettings.SetupGated,
		projectsettings.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
