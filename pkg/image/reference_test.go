// Copyright (c) 2018 Sylabs, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package image

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sylabs/cri/pkg/singularity"
)

func TestParseImageRef(t *testing.T) {
	tt := []struct {
		name        string
		ref         string
		expect      *Reference
		expectError error
	}{
		{
			name: "library with tag",
			ref:  "cloud.sylabs.io/sashayakovtseva/test/image-server:1",
			expect: &Reference{
				uri:     singularity.LibraryDomain,
				tags:    []string{"cloud.sylabs.io/sashayakovtseva/test/image-server:1"},
				digests: nil,
			},
			expectError: nil,
		},
		{
			name: "library without tag",
			ref:  "cloud.sylabs.io/sashayakovtseva/test/image-server",
			expect: &Reference{
				uri:     singularity.LibraryDomain,
				tags:    []string{"cloud.sylabs.io/sashayakovtseva/test/image-server:latest"},
				digests: nil,
			},
			expectError: nil,
		},
		{
			name: "library with digest",
			ref:  "cloud.sylabs.io/sashayakovtseva/test/image-server:sha256.9327532a05078d7efd5a0ef9ace1ee5cd278653d8df53590e2fb7a4a34cb0bb8",
			expect: &Reference{
				uri:     singularity.LibraryDomain,
				tags:    nil,
				digests: []string{"cloud.sylabs.io/sashayakovtseva/test/image-server:sha256.9327532a05078d7efd5a0ef9ace1ee5cd278653d8df53590e2fb7a4a34cb0bb8"},
			},
			expectError: nil,
		},
		{
			name: "docker without tag",
			ref:  "gcr.io/cri-tools/test-image-tags",
			expect: &Reference{
				uri:     singularity.DockerDomain,
				tags:    []string{"gcr.io/cri-tools/test-image-tags:latest"},
				digests: nil,
			},
			expectError: nil,
		},
		{
			name: "docker with tag",
			ref:  "docker.io/gcr.io/cri-tools/test-image-tags:1",
			expect: &Reference{
				uri:     singularity.DockerDomain,
				tags:    []string{"gcr.io/cri-tools/test-image-tags:1"},
				digests: nil,
			},
			expectError: nil,
		},
		{
			name: "docker with digest",
			ref:  "docker.io/gcr.io/cri-tools/test-image-digest@sha256:9179135b4b4cc5a8721e09379244807553c318d92fa3111a65133241551ca343",
			expect: &Reference{
				uri:     singularity.DockerDomain,
				tags:    nil,
				digests: []string{"gcr.io/cri-tools/test-image-digest@sha256:9179135b4b4cc5a8721e09379244807553c318d92fa3111a65133241551ca343"},
			},
			expectError: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ParseRef(tc.ref)
			require.Equal(t, tc.expectError, err)
			require.Equal(t, tc.expect, actual)
		})
	}
}

func TestNormalizedImageRef(t *testing.T) {
	tt := []struct {
		name   string
		ref    string
		expect string
	}{
		{
			name:   "docker image with tag",
			ref:    "gcr.io/cri-tools/test-image-tags:1",
			expect: "gcr.io/cri-tools/test-image-tags:1",
		},
		{
			name:   "docker image without tag",
			ref:    "gcr.io/cri-tools/test-image-tags",
			expect: "gcr.io/cri-tools/test-image-tags:latest",
		},
		{
			name:   "docker image with digest",
			ref:    "gcr.io/cri-tools/test-image-digest@sha256:9179135b4b4cc5a8721e09379244807553c318d92fa3111a65133241551ca343",
			expect: "gcr.io/cri-tools/test-image-digest@sha256:9179135b4b4cc5a8721e09379244807553c318d92fa3111a65133241551ca343",
		},
		{
			name:   "library image with tag",
			ref:    "cloud.sylabs.io/sashayakovtseva/test/image-server:latest",
			expect: "cloud.sylabs.io/sashayakovtseva/test/image-server:latest",
		},
		{
			name:   "library image without tag",
			ref:    "cloud.sylabs.io/sashayakovtseva/test/image-server",
			expect: "cloud.sylabs.io/sashayakovtseva/test/image-server:latest",
		},
		{
			name:   "library image with digest",
			ref:    "cloud.sylabs.io/sashayakovtseva/test/image-server:sha256.9327532a05078d7efd5a0ef9ace1ee5cd278653d8df53590e2fb7a4a34cb0bb8",
			expect: "cloud.sylabs.io/sashayakovtseva/test/image-server:sha256.9327532a05078d7efd5a0ef9ace1ee5cd278653d8df53590e2fb7a4a34cb0bb8",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := NormalizedImageRef(tc.ref)
			require.Equal(t, tc.expect, actual)
		})
	}
}

func TestReferenceDigests(t *testing.T) {
	ref := &Reference{
		digests: []string{
			"gcr.io/cri-tools/test-image@sha256:9179135b4b4cc5a8721e09379244807553c318d92fa3111a65133241551ca343",
			"gcr.io/cri-tools/test-image@sha256:9327532a05078d7efd5a0ef9ace1ee5cd278653d8df53590e2fb7a4a34cb0bb8",
		},
	}

	digests := ref.Digests()
	digests = append(digests, "will-not-affect")
	require.NotEqual(t, ref.Digests(), digests)
	require.NotEqual(t, ref.digests, digests)

	ref.AddDigests([]string{
		"gcr.io/cri-tools/test-image@sha256:73a84b7ecd215008166111f3beb0a8da142535afafa68439e6292d173bc1251f",
		"gcr.io/cri-tools/test-image@sha256:d50278eebfe4ca5655cc28503983f7c947914a34fbbb805481657d39e98f33f0",
	})
	require.ElementsMatch(t, []string{
		"gcr.io/cri-tools/test-image@sha256:9179135b4b4cc5a8721e09379244807553c318d92fa3111a65133241551ca343",
		"gcr.io/cri-tools/test-image@sha256:9327532a05078d7efd5a0ef9ace1ee5cd278653d8df53590e2fb7a4a34cb0bb8",
		"gcr.io/cri-tools/test-image@sha256:73a84b7ecd215008166111f3beb0a8da142535afafa68439e6292d173bc1251f",
		"gcr.io/cri-tools/test-image@sha256:d50278eebfe4ca5655cc28503983f7c947914a34fbbb805481657d39e98f33f0",
	}, ref.Digests())

	ref.RemoveDigest("gcr.io/cri-tools/test-image@sha256:9327532a05078d7efd5a0ef9ace1ee5cd278653d8df53590e2fb7a4a34cb0bb8")
	require.ElementsMatch(t, []string{
		"gcr.io/cri-tools/test-image@sha256:9179135b4b4cc5a8721e09379244807553c318d92fa3111a65133241551ca343",
		"gcr.io/cri-tools/test-image@sha256:73a84b7ecd215008166111f3beb0a8da142535afafa68439e6292d173bc1251f",
		"gcr.io/cri-tools/test-image@sha256:d50278eebfe4ca5655cc28503983f7c947914a34fbbb805481657d39e98f33f0",
	}, ref.Digests())

}

func TestReferenceTags(t *testing.T) {
	ref := &Reference{
		tags: []string{
			"gcr.io/cri-tools/test-image-tags:1",
			"gcr.io/cri-tools/test-image-tags:2",
		},
	}

	tags := ref.Tags()
	tags = append(tags, "will-not-affect")
	require.NotEqual(t, ref.Tags(), tags)
	require.NotEqual(t, ref.tags, tags)

	ref.AddTags([]string{"new-tag", "new-tag-2"})
	require.ElementsMatch(t, []string{
		"gcr.io/cri-tools/test-image-tags:1",
		"gcr.io/cri-tools/test-image-tags:2",
		"new-tag",
		"new-tag-2",
	}, ref.Tags())

	ref.RemoveTag("gcr.io/cri-tools/test-image-tags:2")
	require.ElementsMatch(t, []string{
		"gcr.io/cri-tools/test-image-tags:1",
		"new-tag",
		"new-tag-2",
	}, ref.Tags())

}
