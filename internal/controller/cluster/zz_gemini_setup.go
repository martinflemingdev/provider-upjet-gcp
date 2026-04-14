// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	coderepositoryindex "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/coderepositoryindex"
	codetoolssetting "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/codetoolssetting"
	codetoolssettingbinding "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/codetoolssettingbinding"
	datasharingwithgooglesetting "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/datasharingwithgooglesetting"
	datasharingwithgooglesettingbinding "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/datasharingwithgooglesettingbinding"
	geminigcpenablementsetting "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/geminigcpenablementsetting"
	geminigcpenablementsettingbinding "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/geminigcpenablementsettingbinding"
	loggingsetting "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/loggingsetting"
	loggingsettingbinding "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/loggingsettingbinding"
	releasechannelsetting "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/releasechannelsetting"
	releasechannelsettingbinding "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/releasechannelsettingbinding"
	repositorygroup "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/repositorygroup"
	repositorygroupiammember "github.com/upbound/provider-gcp/v2/internal/controller/cluster/gemini/repositorygroupiammember"
)

// Setup_gemini creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_gemini(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		coderepositoryindex.Setup,
		codetoolssetting.Setup,
		codetoolssettingbinding.Setup,
		datasharingwithgooglesetting.Setup,
		datasharingwithgooglesettingbinding.Setup,
		geminigcpenablementsetting.Setup,
		geminigcpenablementsettingbinding.Setup,
		loggingsetting.Setup,
		loggingsettingbinding.Setup,
		releasechannelsetting.Setup,
		releasechannelsettingbinding.Setup,
		repositorygroup.Setup,
		repositorygroupiammember.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_gemini creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_gemini(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		coderepositoryindex.SetupGated,
		codetoolssetting.SetupGated,
		codetoolssettingbinding.SetupGated,
		datasharingwithgooglesetting.SetupGated,
		datasharingwithgooglesettingbinding.SetupGated,
		geminigcpenablementsetting.SetupGated,
		geminigcpenablementsettingbinding.SetupGated,
		loggingsetting.SetupGated,
		loggingsettingbinding.SetupGated,
		releasechannelsetting.SetupGated,
		releasechannelsettingbinding.SetupGated,
		repositorygroup.SetupGated,
		repositorygroupiammember.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
