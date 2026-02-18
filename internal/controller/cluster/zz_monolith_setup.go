// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	providerconfig "github.com/upbound/provider-gcp/internal/controller/cluster/providerconfig"
	dataset "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/dataset"
	deploymentresourcepool "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/deploymentresourcepool"
	endpoint "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/endpoint"
	endpointwithmodelgardendeployment "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/endpointwithmodelgardendeployment"
	featuregroup "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/featuregroup"
	featuregroupfeature "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/featuregroupfeature"
	featureonlinestore "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/featureonlinestore"
	featureonlinestorefeatureview "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/featureonlinestorefeatureview"
	featurestore "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/featurestore"
	featurestoreentitytype "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/featurestoreentitytype"
	index "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/index"
	indexendpoint "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/indexendpoint"
	indexendpointdeployedindex "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/indexendpointdeployedindex"
	ragengineconfig "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/ragengineconfig"
	tensorboard "github.com/upbound/provider-gcp/internal/controller/cluster/vertexai/tensorboard"
)

// Setup_monolith creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_monolith(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		providerconfig.Setup,
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
		index.Setup,
		indexendpoint.Setup,
		indexendpointdeployedindex.Setup,
		ragengineconfig.Setup,
		tensorboard.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_monolith creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_monolith(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		providerconfig.SetupGated,
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
		index.SetupGated,
		indexendpoint.SetupGated,
		indexendpointdeployedindex.SetupGated,
		ragengineconfig.SetupGated,
		tensorboard.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
