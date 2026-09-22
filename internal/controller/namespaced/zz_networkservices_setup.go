// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	agentconnectivitytemplate "github.com/upbound/provider-gcp/v3/internal/controller/namespaced/networkservices/agentconnectivitytemplate"
	agentgateway "github.com/upbound/provider-gcp/v3/internal/controller/namespaced/networkservices/agentgateway"
	gateway "github.com/upbound/provider-gcp/v3/internal/controller/namespaced/networkservices/gateway"
)

// Setup_networkservices creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_networkservices(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		agentconnectivitytemplate.Setup,
		agentgateway.Setup,
		gateway.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_networkservices creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_networkservices(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		agentconnectivitytemplate.SetupGated,
		agentgateway.SetupGated,
		gateway.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager_networkservices registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager_networkservices(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		agentconnectivitytemplate.SetupWebhookWithManager,
		agentgateway.SetupWebhookWithManager,
		gateway.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
