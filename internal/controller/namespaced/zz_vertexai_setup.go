// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	cacheconfig "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/cacheconfig"
	dataset "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/dataset"
	deploymentresourcepool "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/deploymentresourcepool"
	endpoint "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/endpoint"
	endpointwithmodelgardendeployment "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/endpointwithmodelgardendeployment"
	featuregroup "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/featuregroup"
	featuregroupfeature "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/featuregroupfeature"
	featureonlinestore "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/featureonlinestore"
	featureonlinestorefeatureview "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/featureonlinestorefeatureview"
	featurestore "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/featurestore"
	featurestoreentitytype "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/featurestoreentitytype"
	featurestoreentitytypefeature "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/featurestoreentitytypefeature"
	index "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/index"
	indexendpoint "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/indexendpoint"
	indexendpointdeployedindex "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/indexendpointdeployedindex"
	ragengineconfig "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/ragengineconfig"
	reasoningengine "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/reasoningengine"
	reasoningengineiammember "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/reasoningengineiammember"
	tensorboard "github.com/upbound/provider-gcp/v2/internal/controller/namespaced/vertexai/tensorboard"
)

// Setup_vertexai creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_vertexai(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		cacheconfig.Setup,
		dataset.Setup,
		deploymentresourcepool.Setup,
		endpoint.Setup,
		endpointwithmodelgardendeployment.Setup,
		featuregroup.Setup,
		featuregroupfeature.Setup,
		featureonlinestore.Setup,
		featureonlinestorefeatureview.Setup,
		featurestore.Setup,
		featurestoreentitytype.Setup,
		featurestoreentitytypefeature.Setup,
		index.Setup,
		indexendpoint.Setup,
		indexendpointdeployedindex.Setup,
		ragengineconfig.Setup,
		reasoningengine.Setup,
		reasoningengineiammember.Setup,
		tensorboard.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_vertexai creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_vertexai(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		cacheconfig.SetupGated,
		dataset.SetupGated,
		deploymentresourcepool.SetupGated,
		endpoint.SetupGated,
		endpointwithmodelgardendeployment.SetupGated,
		featuregroup.SetupGated,
		featuregroupfeature.SetupGated,
		featureonlinestore.SetupGated,
		featureonlinestorefeatureview.SetupGated,
		featurestore.SetupGated,
		featurestoreentitytype.SetupGated,
		featurestoreentitytypefeature.SetupGated,
		index.SetupGated,
		indexendpoint.SetupGated,
		indexendpointdeployedindex.SetupGated,
		ragengineconfig.SetupGated,
		reasoningengine.SetupGated,
		reasoningengineiammember.SetupGated,
		tensorboard.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
