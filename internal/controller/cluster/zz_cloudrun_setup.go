// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	domainmapping "github.com/upbound/provider-gcp/v2/internal/controller/cluster/cloudrun/domainmapping"
	service "github.com/upbound/provider-gcp/v2/internal/controller/cluster/cloudrun/service"
	serviceiammember "github.com/upbound/provider-gcp/v2/internal/controller/cluster/cloudrun/serviceiammember"
	v2job "github.com/upbound/provider-gcp/v2/internal/controller/cluster/cloudrun/v2job"
	v2jobiammember "github.com/upbound/provider-gcp/v2/internal/controller/cluster/cloudrun/v2jobiammember"
	v2service "github.com/upbound/provider-gcp/v2/internal/controller/cluster/cloudrun/v2service"
	v2serviceiammember "github.com/upbound/provider-gcp/v2/internal/controller/cluster/cloudrun/v2serviceiammember"
	v2workerpool "github.com/upbound/provider-gcp/v2/internal/controller/cluster/cloudrun/v2workerpool"
	v2workerpooliammember "github.com/upbound/provider-gcp/v2/internal/controller/cluster/cloudrun/v2workerpooliammember"
)

// Setup_cloudrun creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_cloudrun(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		domainmapping.Setup,
		service.Setup,
		serviceiammember.Setup,
		v2job.Setup,
		v2jobiammember.Setup,
		v2service.Setup,
		v2serviceiammember.Setup,
		v2workerpool.Setup,
		v2workerpooliammember.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_cloudrun creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_cloudrun(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		domainmapping.SetupGated,
		service.SetupGated,
		serviceiammember.SetupGated,
		v2job.SetupGated,
		v2jobiammember.SetupGated,
		v2service.SetupGated,
		v2serviceiammember.SetupGated,
		v2workerpool.SetupGated,
		v2workerpooliammember.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
