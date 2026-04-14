// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	aclconfig "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/aclconfig"
	assistant "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/assistant"
	chatengine "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/chatengine"
	cmekconfig "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/cmekconfig"
	control "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/control"
	dataconnector "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/dataconnector"
	datastore "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/datastore"
	licenseconfig "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/licenseconfig"
	recommendationengine "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/recommendationengine"
	schema "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/schema"
	searchengine "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/searchengine"
	servingconfig "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/servingconfig"
	sitemap "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/sitemap"
	targetsite "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/targetsite"
	userstore "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/userstore"
	widgetconfig "github.com/upbound/provider-gcp/v2/internal/controller/cluster/discoveryengine/widgetconfig"
)

// Setup_discoveryengine creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_discoveryengine(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		aclconfig.Setup,
		assistant.Setup,
		chatengine.Setup,
		cmekconfig.Setup,
		control.Setup,
		dataconnector.Setup,
		datastore.Setup,
		licenseconfig.Setup,
		recommendationengine.Setup,
		schema.Setup,
		searchengine.Setup,
		servingconfig.Setup,
		sitemap.Setup,
		targetsite.Setup,
		userstore.Setup,
		widgetconfig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_discoveryengine creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_discoveryengine(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		aclconfig.SetupGated,
		assistant.SetupGated,
		chatengine.SetupGated,
		cmekconfig.SetupGated,
		control.SetupGated,
		dataconnector.SetupGated,
		datastore.SetupGated,
		licenseconfig.SetupGated,
		recommendationengine.SetupGated,
		schema.SetupGated,
		searchengine.SetupGated,
		servingconfig.SetupGated,
		sitemap.SetupGated,
		targetsite.SetupGated,
		userstore.SetupGated,
		widgetconfig.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
