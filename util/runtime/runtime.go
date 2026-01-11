// Copyright 2026 Cohesity, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	errorsutil "k8s.io/apimachinery/pkg/util/errors"
	criapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/klog/v2"
)

const (
	RemoveContainerRetry = 5
)

var (
	// DefaultCRISocket is the default CRI socket
	DefaultCRISocket = "unix:///var/run/containerd/containerd.sock"

	// defaultKnownCRISockets holds the set of known CRI endpoints.
	defaultKnownCRISockets = []string{
		"unix:///var/run/containerd/containerd.sock",
		"unix:///var/run/crio/crio.sock",
		"unix:///var/run/docker.sock",
	}

	// ErrRuntimeConditionNotTrue is returned when a container runtime condition is not true.
	ErrRuntimeConditionNotTrue = errors.New("container runtime condition is not true")

	// ErrMultipleCRIEndpoints is returned when multiple CRI endpoints are found on the host.
	ErrMultipleCRIEndpoints = errors.New("found multiple CRI endpoints on the host")
)

// ContainerRuntime is an interface for working with container runtimes
type ContainerRuntime interface {
	Connect() error
	SetImpl(impl)
	IsRunning() error
	ListKubeContainers() ([]string, error)
	RemoveContainers(containers []string) error
}

// CRIRuntime is a struct that interfaces with the CRI
type CRIRuntime struct {
	impl           impl
	runtimeService criapi.RuntimeService
	imageService   criapi.ImageManagerService
	criSocket      string
}

// defaultTimeout is the default timeout inherited by crictl
const defaultTimeout = 2 * time.Second

// NewContainerRuntime sets up and returns a ContainerRuntime struct
func NewContainerRuntime(criSocket string) ContainerRuntime {
	return &CRIRuntime{
		impl:      &defaultImpl{},
		criSocket: criSocket,
	}
}

// SetImpl can be used to set the internal implementation for testing purposes.
func (runtime *CRIRuntime) SetImpl(impl impl) {
	runtime.impl = impl
}

// Connect establishes a connection with the CRI runtime.
func (runtime *CRIRuntime) Connect() error {
	runtimeService, err := runtime.impl.NewRemoteRuntimeService(runtime.criSocket, defaultTimeout)
	if err != nil {
		return fmt.Errorf("failed to create new CRI runtime service: %w", err)
	}
	runtime.runtimeService = runtimeService

	imageService, err := runtime.impl.NewRemoteImageService(runtime.criSocket, defaultTimeout)
	if err != nil {
		return fmt.Errorf("failed to create new CRI image service: %w", err)
	}
	runtime.imageService = imageService

	return nil
}

// IsRunning checks if runtime is running.
func (runtime *CRIRuntime) IsRunning() error {
	ctx, cancel := defaultContext()
	defer cancel()

	res, err := runtime.impl.Status(ctx, runtime.runtimeService, false)
	if err != nil {
		return fmt.Errorf("container runtime is not running: %w", err)
	}

	for _, condition := range res.GetStatus().GetConditions() {
		if condition.GetType() == runtimeapi.RuntimeReady && // NetworkReady will not be tested on purpose
			!condition.GetStatus() {
			return fmt.Errorf(
				"%w: condition %q, reason: %s, message: %s",
				ErrRuntimeConditionNotTrue, condition.GetType(), condition.GetReason(), condition.GetMessage(),
			)
		}
	}

	return nil
}

// ListKubeContainers lists running k8s CRI pods
func (runtime *CRIRuntime) ListKubeContainers() ([]string, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	sandboxes, err := runtime.impl.ListPodSandbox(ctx, runtime.runtimeService, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list pod sandboxes: %w", err)
	}

	pods := []string{}
	for _, sandbox := range sandboxes {
		klog.Infof("Pod Name: %s, Pod ID: %s", sandbox.GetMetadata().GetName(), sandbox.GetId())
		pods = append(pods, sandbox.GetId())
	}
	return pods, nil
}

// RemoveContainers removes running k8s pods
func (runtime *CRIRuntime) RemoveContainers(containers []string) error {
	errs := []error{}
	for _, container := range containers {
		var lastErr error
		for i := 0; i < RemoveContainerRetry; i++ {
			klog.Infof("Attempting to remove container %v", container)

			ctx, cancel := defaultContext()
			if err := runtime.impl.StopPodSandbox(ctx, runtime.runtimeService, container); err != nil {
				lastErr = fmt.Errorf("failed to stop running pod %s: %w", container, err)
				cancel()
				continue
			}
			klog.Infof("Successfully stopped container %v", container)
			cancel()

			ctx, cancel = defaultContext()
			if err := runtime.impl.RemovePodSandbox(ctx, runtime.runtimeService, container); err != nil {
				lastErr = fmt.Errorf("failed to remove pod %s: %w", container, err)
				cancel()
				continue
			}
			klog.Infof("Successfully removed container %v", container)
			cancel()

			lastErr = nil
			break
		}

		if lastErr != nil {
			errs = append(errs, lastErr)
		}
	}
	return errorsutil.NewAggregate(errs)
}

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}

// detectCRISocketImpl is separated out only for test purposes, DON'T call it directly, use DetectCRISocket instead
func detectCRISocketImpl(isSocket func(string) bool, knownCRISockets []string) (string, error) {
	foundCRISockets := []string{}

	for _, socket := range knownCRISockets {
		if isSocket(socket) {
			foundCRISockets = append(foundCRISockets, socket)
		}
	}

	switch len(foundCRISockets) {
	case 0:
		// Fall back to the default socket if no CRI is detected, we can error out later on if we need it
		return DefaultCRISocket, nil
	case 1:
		// Precisely one CRI found, use that
		return foundCRISockets[0], nil
	default:
		// Multiple CRIs installed?
		return "", fmt.Errorf("%w. Please define which one do you wish to use by setting the 'criSocket' field in the kubeadm configuration file: %s",
			ErrMultipleCRIEndpoints, strings.Join(foundCRISockets, ", "))
	}
}

// DetectCRISocket uses a list of known CRI sockets to detect one. If more than one or none is discovered, an error is returned.
func DetectCRISocket() (string, error) {
	return detectCRISocketImpl(isExistingSocket, defaultKnownCRISockets)
}
