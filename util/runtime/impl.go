// Copyright 2026 Cohesity, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"context"
	"fmt"
	"time"

	criapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	criclient "k8s.io/cri-client/pkg"
)

type defaultImpl struct{}

type impl interface {
	NewRemoteRuntimeService(endpoint string, connectionTimeout time.Duration) (criapi.RuntimeService, error)
	NewRemoteImageService(endpoint string, connectionTimeout time.Duration) (criapi.ImageManagerService, error)
	Status(ctx context.Context, runtimeService criapi.RuntimeService, verbose bool) (*runtimeapi.StatusResponse, error)
	ListPodSandbox(ctx context.Context, runtimeService criapi.RuntimeService, filter *runtimeapi.PodSandboxFilter) ([]*runtimeapi.PodSandbox, error)
	StopPodSandbox(ctx context.Context, runtimeService criapi.RuntimeService, sandboxID string) error
	RemovePodSandbox(ctx context.Context, runtimeService criapi.RuntimeService, podSandboxID string) error
	PullImage(ctx context.Context, imageService criapi.ImageManagerService, image *runtimeapi.ImageSpec, auth *runtimeapi.AuthConfig, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error)
	ImageStatus(ctx context.Context, imageService criapi.ImageManagerService, image *runtimeapi.ImageSpec, verbose bool) (*runtimeapi.ImageStatusResponse, error)
}

func (*defaultImpl) NewRemoteRuntimeService(endpoint string, connectionTimeout time.Duration) (criapi.RuntimeService, error) {
	runtimeService, err := criclient.NewRemoteRuntimeService(endpoint, defaultTimeout, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote runtime service: %w", err)
	}
	return runtimeService, nil
}

func (*defaultImpl) NewRemoteImageService(endpoint string, connectionTimeout time.Duration) (criapi.ImageManagerService, error) {
	imageService, err := criclient.NewRemoteImageService(endpoint, connectionTimeout, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote image service: %w", err)
	}
	return imageService, nil
}

func (*defaultImpl) Status(ctx context.Context, runtimeService criapi.RuntimeService, verbose bool) (*runtimeapi.StatusResponse, error) {
	status, err := runtimeService.Status(ctx, verbose)
	if err != nil {
		return nil, fmt.Errorf("failed to get runtime status: %w", err)
	}
	return status, nil
}

func (*defaultImpl) ListPodSandbox(ctx context.Context, runtimeService criapi.RuntimeService, filter *runtimeapi.PodSandboxFilter) ([]*runtimeapi.PodSandbox, error) {
	sandboxes, err := runtimeService.ListPodSandbox(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list pod sandbox: %w", err)
	}
	return sandboxes, nil
}

func (*defaultImpl) StopPodSandbox(ctx context.Context, runtimeService criapi.RuntimeService, sandboxID string) error {
	if err := runtimeService.StopPodSandbox(ctx, sandboxID); err != nil {
		return fmt.Errorf("failed to stop pod sandbox: %w", err)
	}
	return nil
}

func (*defaultImpl) RemovePodSandbox(ctx context.Context, runtimeService criapi.RuntimeService, podSandboxID string) error {
	if err := runtimeService.RemovePodSandbox(ctx, podSandboxID); err != nil {
		return fmt.Errorf("failed to remove pod sandbox: %w", err)
	}
	return nil
}

func (*defaultImpl) PullImage(ctx context.Context, imageService criapi.ImageManagerService, image *runtimeapi.ImageSpec, auth *runtimeapi.AuthConfig, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	imageRef, err := imageService.PullImage(ctx, image, auth, podSandboxConfig)
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	return imageRef, nil
}

func (*defaultImpl) ImageStatus(ctx context.Context, imageService criapi.ImageManagerService, image *runtimeapi.ImageSpec, verbose bool) (*runtimeapi.ImageStatusResponse, error) {
	status, err := imageService.ImageStatus(ctx, image, verbose)
	if err != nil {
		return nil, fmt.Errorf("failed to get image status: %w", err)
	}
	return status, nil
}
