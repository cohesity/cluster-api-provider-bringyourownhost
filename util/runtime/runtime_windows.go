//go:build windows
// +build windows

/*
Copyright 2026 Cohesity, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package runtime

import (
	"net/url"

	winio "github.com/Microsoft/go-winio"
)

// isExistingSocket checks if path exists and is domain socket
func isExistingSocket(path string) bool {
	u, err := url.Parse(path)
	if err != nil {
		// should not happen, since we are trying to access known / hardcoded sockets
		return false
	}

	// the dial path must be without "npipe://"
	_, err = winio.DialPipe(u.Path, nil)
	if err != nil {
		return false
	}

	return true
}
