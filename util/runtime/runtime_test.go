// Copyright 2026 Cohesity, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "k8s.io/cri-api/pkg/apis/runtime/v1"
)

var errTest = fmt.Errorf("test")

func TestNewContainerRuntime(t *testing.T) {
	for _, tc := range []struct {
		name        string
		prepare     func(*fakeImpl)
		shouldError bool
	}{
		{
			name:        "valid",
			shouldError: false,
		},
		{
			name: "invalid: new runtime service fails",
			prepare: func(mock *fakeImpl) {
				mock.NewRemoteRuntimeServiceReturns(nil, errTest)
			},
			shouldError: true,
		},
		{
			name: "invalid: new image service fails",
			prepare: func(mock *fakeImpl) {
				mock.NewRemoteImageServiceReturns(nil, errTest)
			},
			shouldError: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			containerRuntime := NewContainerRuntime("")
			mock := &fakeImpl{}
			if tc.prepare != nil {
				tc.prepare(mock)
			}
			containerRuntime.SetImpl(mock)

			err := containerRuntime.Connect()

			assert.Equal(t, tc.shouldError, err != nil)
		})
	}
}

func TestIsRunning(t *testing.T) {
	for _, tc := range []struct {
		name        string
		prepare     func(*fakeImpl)
		shouldError bool
	}{
		{
			name:        "valid",
			shouldError: false,
		},
		{
			name: "invalid: runtime status fails",
			prepare: func(mock *fakeImpl) {
				mock.StatusReturns(nil, errTest)
			},
			shouldError: true,
		},
		{
			name: "invalid: runtime condition status not 'true'",
			prepare: func(mock *fakeImpl) {
				mock.StatusReturns(&v1.StatusResponse{
					Status: &v1.RuntimeStatus{
						Conditions: []*v1.RuntimeCondition{
							{
								Type:   v1.RuntimeReady,
								Status: false,
							},
						},
					},
				}, nil)
			},
			shouldError: true,
		},
		{
			name: "valid: runtime condition type does not match",
			prepare: func(mock *fakeImpl) {
				mock.StatusReturns(&v1.StatusResponse{
					Status: &v1.RuntimeStatus{
						Conditions: []*v1.RuntimeCondition{
							{
								Type:   v1.NetworkReady,
								Status: false,
							},
						},
					},
				}, nil)
			},
			shouldError: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			containerRuntime := NewContainerRuntime("")
			mock := &fakeImpl{}
			if tc.prepare != nil {
				tc.prepare(mock)
			}
			containerRuntime.SetImpl(mock)

			err := containerRuntime.IsRunning()

			assert.Equal(t, tc.shouldError, err != nil)
		})
	}
}

func TestListKubeContainers(t *testing.T) {
	for _, tc := range []struct {
		name        string
		expected    []string
		prepare     func(*fakeImpl)
		shouldError bool
	}{
		{
			name: "valid",
			prepare: func(mock *fakeImpl) {
				mock.ListPodSandboxReturns([]*v1.PodSandbox{
					{Id: "first"},
					{Id: "second"},
				}, nil)
			},
			expected:    []string{"first", "second"},
			shouldError: false,
		},
		{
			name: "invalid: list pod sandbox fails",
			prepare: func(mock *fakeImpl) {
				mock.ListPodSandboxReturns(nil, errTest)
			},
			shouldError: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			containerRuntime := NewContainerRuntime("")
			mock := &fakeImpl{}
			if tc.prepare != nil {
				tc.prepare(mock)
			}
			containerRuntime.SetImpl(mock)

			containers, err := containerRuntime.ListKubeContainers()

			assert.Equal(t, tc.shouldError, err != nil)
			assert.Equal(t, tc.expected, containers)
		})
	}
}

func TestRemoveContainers(t *testing.T) {
	for _, tc := range []struct {
		name        string
		containers  []string
		prepare     func(*fakeImpl)
		shouldError bool
	}{
		{
			name: "valid",
		},
		{
			name:        "valid: two containers",
			containers:  []string{"1", "2"},
			shouldError: false,
		},
		{
			name:       "invalid: remove pod sandbox fails",
			containers: []string{"1"},
			prepare: func(mock *fakeImpl) {
				mock.RemovePodSandboxReturns(errTest)
			},
			shouldError: true,
		},
		{
			name:       "invalid: stop pod sandbox fails",
			containers: []string{"1"},
			prepare: func(mock *fakeImpl) {
				mock.StopPodSandboxReturns(errTest)
			},
			shouldError: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			containerRuntime := NewContainerRuntime("")
			mock := &fakeImpl{}
			if tc.prepare != nil {
				tc.prepare(mock)
			}
			containerRuntime.SetImpl(mock)

			err := containerRuntime.RemoveContainers(tc.containers)

			assert.Equal(t, tc.shouldError, err != nil)
		})
	}
}

func TestIsExistingSocket(t *testing.T) {
	// this test is not expected to work on Windows
	if runtime.GOOS == "windows" {
		return
	}

	const tempPrefix = "test.kubeadm.runtime.isExistingSocket."
	tests := []struct {
		name string
		proc func(*testing.T)
	}{
		{
			name: "Valid domain socket is detected as such",
			proc: func(t *testing.T) {
				tmpFile, err := os.CreateTemp("", tempPrefix)
				if err != nil {
					t.Fatalf("unexpected error by TempFile: %v", err)
				}
				theSocket := tmpFile.Name()
				_ = os.Remove(theSocket)
				_ = tmpFile.Close()

				var lc net.ListenConfig
				con, err := lc.Listen(context.Background(), "unix", theSocket)
				if err != nil {
					t.Fatalf("unexpected error while dialing a socket: %v", err)
				}
				defer func() { _ = con.Close() }()

				if !isExistingSocket("unix://" + theSocket) {
					t.Fatalf("isExistingSocket(%q) gave unexpected result. Should have been true, instead of false", theSocket)
				}
			},
		},
		{
			name: "Regular file is not a domain socket",
			proc: func(t *testing.T) {
				tmpFile, err := os.CreateTemp("", tempPrefix)
				if err != nil {
					t.Fatalf("unexpected error by TempFile: %v", err)
				}
				theSocket := tmpFile.Name()
				defer func() { _ = os.Remove(theSocket) }()
				_ = tmpFile.Close()

				if isExistingSocket(theSocket) {
					t.Fatalf("isExistingSocket(%q) gave unexpected result. Should have been false, instead of true", theSocket)
				}
			},
		},
		{
			name: "Non existent socket is not a domain socket",
			proc: func(t *testing.T) {
				const theSocket = "/non/existent/socket"
				if isExistingSocket(theSocket) {
					t.Fatalf("isExistingSocket(%q) gave unexpected result. Should have been false, instead of true", theSocket)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.proc)
	}
}

func TestDetectCRISocketImpl(t *testing.T) {
	tests := []struct {
		name            string
		existingSockets []string
		expectedError   bool
		expectedSocket  string
	}{
		{
			name:            "No existing sockets, use default",
			existingSockets: []string{},
			expectedError:   false,
			expectedSocket:  DefaultCRISocket,
		},
		{
			name:            "One valid CRI socket leads to success",
			existingSockets: []string{"unix:///foo/bar.sock"},
			expectedError:   false,
			expectedSocket:  "unix:///foo/bar.sock",
		},
		{
			name: "Multiple CRI sockets lead to an error",
			existingSockets: []string{
				"unix:///foo/bar.sock",
				"unix:///foo/baz.sock",
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			socket, err := detectCRISocketImpl(func(path string) bool {
				return slices.Contains(test.existingSockets, path)
			}, test.existingSockets)

			if (err != nil) != test.expectedError {
				t.Fatalf("detectCRISocketImpl returned unexpected result\n\tExpected error: %t\n\tGot error: %t", test.expectedError, err != nil)
			}
			if !test.expectedError && socket != test.expectedSocket {
				t.Fatalf("detectCRISocketImpl returned unexpected CRI socket\n\tExpected socket: %s\n\tReturned socket: %s",
					test.expectedSocket, socket)
			}
		})
	}
}
