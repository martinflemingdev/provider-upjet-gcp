// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package common

import (
	jresource "github.com/crossplane/upjet/v2/pkg/resource"
	"testing"
)

type loggingBucket struct {
	jresource.Terraformed
	id string
}

func (b *loggingBucket) GetID() string { return b.id }

func TestExtractLoggingBucketDestination(t *testing.T) {
	for _, tc := range []struct{ name, id, want string }{
		{name: "unobserved"},
		{name: "regional bucket", id: "projects/example/locations/us-central1/buckets/audit", want: "logging.googleapis.com/projects/example/locations/us-central1/buckets/audit"},
		{name: "default bucket", id: "projects/example/locations/global/buckets/_Default", want: "logging.googleapis.com/projects/example/locations/global/buckets/_Default"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := ExtractLoggingBucketDestination()(&loggingBucket{id: tc.id}); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
