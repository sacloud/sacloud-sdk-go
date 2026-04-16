// Copyright 2022-2025 The sacloud/packages-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"io"
	"os/exec"
	"testing"
)

func StartCommand(t *testing.T, command string, args ...string) error {
	_, err := StartCommandWithStdErr(t, command, args...)
	return err
}

func StartCommandWithStdOut(t *testing.T, command string, args ...string) (io.Reader, error) {
	cmd := runCommand(t, command, args...)
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return stdOut, nil
}

func StartCommandWithStdErr(t *testing.T, command string, args ...string) (io.Reader, error) {
	cmd := runCommand(t, command, args...)
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return stdErr, nil
}

func RunCommand(t *testing.T, command string, args ...string) error {
	return runCommand(t, command, args...).Run()
}

func RunCommandWithOutput(t *testing.T, command string, args ...string) ([]byte, error) {
	return runCommand(t, command, args...).Output()
}

func RunCommandWithCombinedOutput(t *testing.T, command string, args ...string) ([]byte, error) {
	return runCommand(t, command, args...).CombinedOutput()
}

func runCommand(t *testing.T, command string, args ...string) *exec.Cmd {
	return exec.Command(command, args...)
}
