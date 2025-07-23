// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	dataset "github.com/upbound/provider-gcp/internal/controller/vertexai/dataset"
	endpoint "github.com/upbound/provider-gcp/internal/controller/vertexai/endpoint"
	endpointiammember "github.com/upbound/provider-gcp/internal/controller/vertexai/endpointiammember"
	featurestore "github.com/upbound/provider-gcp/internal/controller/vertexai/featurestore"
	featurestoreentitytype "github.com/upbound/provider-gcp/internal/controller/vertexai/featurestoreentitytype"
	tensorboard "github.com/upbound/provider-gcp/internal/controller/vertexai/tensorboard"
)

// Setup_vertexai creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_vertexai(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		dataset.Setup,
		endpoint.Setup,
		endpointiammember.Setup,
		featurestore.Setup,
		featurestoreentitytype.Setup,
		tensorboard.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
