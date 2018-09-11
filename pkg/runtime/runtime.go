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
	"fmt"
	"os/exec"

	"github.com/sylabs/cri/pkg/singularity"
	"k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

// SingularityRuntime implements k8s RuntimeService interface.
type SingularityRuntime struct {
	singularity string
}

// NewSingularityRuntime initializes and returns SingularityRuntime.
// Singularity must be installed on the host otherwise it will return an error.
func NewSingularityRuntime() (*SingularityRuntime, error) {
	s, err := exec.LookPath(singularity.RuntimeName)
	if err != nil {
		return nil, fmt.Errorf("could not find %s daemon on this machine: %v", singularity.RuntimeName, err)
	}
	return &SingularityRuntime{
		singularity: s,
	}, nil
}

// Version returns the runtime name, runtime version and runtime API version
func (s *SingularityRuntime) Version(ctx context.Context, req *v1alpha2.VersionRequest) (*v1alpha2.VersionResponse, error) {
	const kubeAPIVersion = "0.1.0"

	syVersion, err := exec.Command(s.singularity, "version").Output()
	if err != nil {
		return nil, err
	}

	return &v1alpha2.VersionResponse{
		Version:           kubeAPIVersion, // todo or use req.Version?
		RuntimeName:       singularity.RuntimeName,
		RuntimeVersion:    string(syVersion),
		RuntimeApiVersion: string(syVersion),
	}, nil
}

// RunPodSandbox creates and starts a pod-level sandbox. Runtimes must ensure
// the sandbox is in the ready state on success.
func (s *SingularityRuntime) RunPodSandbox(ctx context.Context, req *v1alpha2.RunPodSandboxRequest) (*v1alpha2.RunPodSandboxResponse, error) {
	return &v1alpha2.RunPodSandboxResponse{}, nil
}

// StopPodSandbox stops any running process that is part of the sandbox and
// reclaims network resources (e.g., IP addresses) allocated to the sandbox.
// If there are any running containers in the sandbox, they must be forcibly
// terminated.
// This call is idempotent, and must not return an error if all relevant
// resources have already been reclaimed. kubelet will call StopPodSandbox
// at least once before calling RemovePodSandbox. It will also attempt to
// reclaim resources eagerly, as soon as a sandbox is not needed. Hence,
// multiple StopPodSandbox calls are expected.
func (s *SingularityRuntime) StopPodSandbox(context.Context, *v1alpha2.StopPodSandboxRequest) (*v1alpha2.StopPodSandboxResponse, error) {
	return &v1alpha2.StopPodSandboxResponse{}, nil
}

// RemovePodSandbox removes the sandbox. If there are any running containers
// in the sandbox, they must be forcibly terminated and removed.
// This call is idempotent, and must not return an error if the sandbox has
// already been removed.
func (s *SingularityRuntime) RemovePodSandbox(context.Context, *v1alpha2.RemovePodSandboxRequest) (*v1alpha2.RemovePodSandboxResponse, error) {
	return &v1alpha2.RemovePodSandboxResponse{}, nil
}

// PodSandboxStatus returns the status of the PodSandbox. If the PodSandbox is not
// present, returns an error.
func (s *SingularityRuntime) PodSandboxStatus(context.Context, *v1alpha2.PodSandboxStatusRequest) (*v1alpha2.PodSandboxStatusResponse, error) {
	return &v1alpha2.PodSandboxStatusResponse{}, nil
}

// ListPodSandbox returns a list of PodSandboxes.
func (s *SingularityRuntime) ListPodSandbox(context.Context, *v1alpha2.ListPodSandboxRequest) (*v1alpha2.ListPodSandboxResponse, error) {
	return &v1alpha2.ListPodSandboxResponse{}, nil
}

// CreateContainer creates a new container in specified PodSandbox
func (s *SingularityRuntime) CreateContainer(context.Context, *v1alpha2.CreateContainerRequest) (*v1alpha2.CreateContainerResponse, error) {
	return &v1alpha2.CreateContainerResponse{}, nil
}

// StartContainer starts the container.
func (s *SingularityRuntime) StartContainer(context.Context, *v1alpha2.StartContainerRequest) (*v1alpha2.StartContainerResponse, error) {
	return &v1alpha2.StartContainerResponse{}, nil
}

// StopContainer stops a running container with a grace period (i.e., timeout).
// This call is idempotent, and must not return an error if the container has
// already been stopped.
// TODO: what must the runtime do after the grace period is reached?
func (s *SingularityRuntime) StopContainer(context.Context, *v1alpha2.StopContainerRequest) (*v1alpha2.StopContainerResponse, error) {
	return &v1alpha2.StopContainerResponse{}, nil
}

// RemoveContainer removes the container. If the container is running, the
// container must be forcibly removed.
// This call is idempotent, and must not return an error if the container has
// already been removed.
func (s *SingularityRuntime) RemoveContainer(context.Context, *v1alpha2.RemoveContainerRequest) (*v1alpha2.RemoveContainerResponse, error) {
	return &v1alpha2.RemoveContainerResponse{}, nil
}

// ListContainers lists all containers by filters.
func (s *SingularityRuntime) ListContainers(context.Context, *v1alpha2.ListContainersRequest) (*v1alpha2.ListContainersResponse, error) {
	return &v1alpha2.ListContainersResponse{}, nil
}

// ContainerStatus returns status of the container. If the container is not
// present, returns an error.
func (s *SingularityRuntime) ContainerStatus(context.Context, *v1alpha2.ContainerStatusRequest) (*v1alpha2.ContainerStatusResponse, error) {
	return &v1alpha2.ContainerStatusResponse{}, nil
}

// UpdateContainerResources updates ContainerConfig of the container.
func (s *SingularityRuntime) UpdateContainerResources(context.Context, *v1alpha2.UpdateContainerResourcesRequest) (*v1alpha2.UpdateContainerResourcesResponse, error) {
	return &v1alpha2.UpdateContainerResourcesResponse{}, nil
}

// ReopenContainerLog asks runtime to reopen the stdout/stderr log file
// for the container. This is often called after the log file has been
// rotated. If the container is not running, container runtime can choose
// to either create a new log file and return nil, or return an error.
// Once it returns error, new container log file MUST NOT be created.
func (s *SingularityRuntime) ReopenContainerLog(context.Context, *v1alpha2.ReopenContainerLogRequest) (*v1alpha2.ReopenContainerLogResponse, error) {
	return &v1alpha2.ReopenContainerLogResponse{}, nil
}

// ExecSync runs a command in a container synchronously.
func (s *SingularityRuntime) ExecSync(context.Context, *v1alpha2.ExecSyncRequest) (*v1alpha2.ExecSyncResponse, error) {
	return &v1alpha2.ExecSyncResponse{}, nil
}

// Exec prepares a streaming endpoint to execute a command in the container.
func (s *SingularityRuntime) Exec(context.Context, *v1alpha2.ExecRequest) (*v1alpha2.ExecResponse, error) {
	return &v1alpha2.ExecResponse{}, nil
}

// Attach prepares a streaming endpoint to attach to a running container.
func (s *SingularityRuntime) Attach(context.Context, *v1alpha2.AttachRequest) (*v1alpha2.AttachResponse, error) {
	return &v1alpha2.AttachResponse{}, nil
}

// PortForward prepares a streaming endpoint to forward ports from a PodSandbox.
func (s *SingularityRuntime) PortForward(context.Context, *v1alpha2.PortForwardRequest) (*v1alpha2.PortForwardResponse, error) {
	return &v1alpha2.PortForwardResponse{}, nil
}

// ContainerStats returns stats of the container. If the container does not
// exist, the call returns an error.
func (s *SingularityRuntime) ContainerStats(context.Context, *v1alpha2.ContainerStatsRequest) (*v1alpha2.ContainerStatsResponse, error) {
	return &v1alpha2.ContainerStatsResponse{}, nil
}

// ListContainerStats returns stats of all running containers.
func (s *SingularityRuntime) ListContainerStats(context.Context, *v1alpha2.ListContainerStatsRequest) (*v1alpha2.ListContainerStatsResponse, error) {
	return &v1alpha2.ListContainerStatsResponse{}, nil
}

// UpdateRuntimeConfig updates the runtime configuration based on the given request.
func (s *SingularityRuntime) UpdateRuntimeConfig(context.Context, *v1alpha2.UpdateRuntimeConfigRequest) (*v1alpha2.UpdateRuntimeConfigResponse, error) {
	return &v1alpha2.UpdateRuntimeConfigResponse{}, nil
}

// Status returns the status of the runtime.
func (s *SingularityRuntime) Status(ctx context.Context, req *v1alpha2.StatusRequest) (*v1alpha2.StatusResponse, error) {
	runtimeReady := &v1alpha2.RuntimeCondition{
		Type:   v1alpha2.RuntimeReady,
		Status: true,
	}
	networkReady := &v1alpha2.RuntimeCondition{
		Type:   v1alpha2.NetworkReady,
		Status: true,
	}
	conditions := []*v1alpha2.RuntimeCondition{runtimeReady, networkReady}

	status := &v1alpha2.RuntimeStatus{Conditions: conditions}
	return &v1alpha2.StatusResponse{Status: status}, nil
}
