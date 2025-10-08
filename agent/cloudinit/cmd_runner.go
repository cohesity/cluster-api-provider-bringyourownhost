// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cloudinit

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

//counterfeiter:generate . ICmdRunner
type ICmdRunner interface {
	RunCmd(context.Context, string) error
	RunCmdWithTimeout(context.Context, string, time.Duration) error
}

// CmdRunner default implementer of ICmdRunner
// TODO reevaluate empty interface/struct
type CmdRunner struct{}

// RunCmd executes the command string
func (r CmdRunner) RunCmd(ctx context.Context, cmd string) error {
	command := exec.CommandContext(ctx, "/bin/bash", "-c", cmd)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	if err := command.Run(); err != nil {
		return fmt.Errorf("failed to run command: %s: %w", cmd, err)
	}
	return nil
}

// RunCmdWithTimeout executes the command string with a specified timeout
// The function will wait for the command to complete within the timeout duration
func (r CmdRunner) RunCmdWithTimeout(ctx context.Context, cmd string, timeout time.Duration) error {
	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	command := exec.CommandContext(timeoutCtx, "/bin/bash", "-c", cmd)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	// Wait for the command to complete within the timeout
	if err := command.Run(); err != nil {
		return fmt.Errorf("failed to run command with timeout %v: %s: %w", timeout, cmd, err)
	}
	return nil
}
