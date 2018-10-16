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

package runtime

import (
	"context"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

func TestSingularityRuntimeService_Version(t *testing.T) {
	s, err := NewSingularityRuntime(nil)
	require.NoError(t, err, "could not create new runtime service")

	expectedVersion, err := exec.Command(s.singularity, "version").Output()
	require.NoError(t, err, "could not run version command against singularity")

	actualVersion, err := s.Version(context.Background(), &v1alpha2.VersionRequest{})
	require.NoError(t, err, "could not query runtime version")
	require.Equal(t, &v1alpha2.VersionResponse{
		Version:           "0.1.0",
		RuntimeName:       "singularity",
		RuntimeVersion:    string(expectedVersion),
		RuntimeApiVersion: string(expectedVersion),
	}, actualVersion, "runtime version mismatch")

}
