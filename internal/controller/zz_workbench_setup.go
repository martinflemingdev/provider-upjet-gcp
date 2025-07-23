// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	instance "github.com/upbound/provider-gcp/internal/controller/workbench/instance"
	instanceiammember "github.com/upbound/provider-gcp/internal/controller/workbench/instanceiammember"
)

// Setup_workbench creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_workbench(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		instance.Setup,
		instanceiammember.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
