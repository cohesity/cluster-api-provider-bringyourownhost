// Copyright 2026 Cohesity, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build !windows
// +build !windows

package runtime

import (
	"context"
	"net"
	"net/url"
)

// isExistingSocket checks if path exists and is domain socket
func isExistingSocket(path string) bool {
	u, err := url.Parse(path)
	if err != nil {
		// should not happen, since we are trying to access known / hardcoded sockets
		return false
	}

	var d net.Dialer
	c, err := d.DialContext(context.Background(), u.Scheme, u.Path)
	if err != nil {
		return false
	}
	_ = c.Close()
	return true
}
