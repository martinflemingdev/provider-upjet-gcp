// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	dataset "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/dataset"
	deploymentresourcepool "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/deploymentresourcepool"
	endpoint "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/endpoint"
	endpointwithmodelgardendeployment "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/endpointwithmodelgardendeployment"
	featurestore "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/featurestore"
	featurestoreentitytype "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/featurestoreentitytype"
	index "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/index"
	indexendpoint "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/indexendpoint"
	indexendpointdeployedindex "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/indexendpointdeployedindex"
	ragengineconfig "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/ragengineconfig"
	tensorboard "github.com/upbound/provider-gcp/internal/controller/namespaced/vertexai/tensorboard"
)

// Setup_vertexai creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_vertexai(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		dataset.Setup,
		deploymentresourcepool.Setup,
		endpoint.Setup,
		endpointwithmodelgardendeployment.Setup,
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

// SetupGated_vertexai creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_vertexai(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		dataset.SetupGated,
		deploymentresourcepool.SetupGated,
		endpoint.SetupGated,
		endpointwithmodelgardendeployment.SetupGated,
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
