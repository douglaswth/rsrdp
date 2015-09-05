// The MIT License (MIT)
//
// Copyright (c) 2015 Douglas Thrift
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// +build !darwin,!windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func rdpLaunchNative(instance *Instance, private bool, index int, arguments []string, prompt bool, username string) error {
	client, options, err := rdpFindClient()
	if err != nil {
		return err
	}

	count := len(options) + len(arguments)
	if rdpIsRemmina(client) {
		count += 6
	} else {
		count += 5
		if !prompt {
			count += 2
		}
	}
	args := make([]string, 0, count)

	if rdpIsRemmina(client) {
		file, err := rdpCreateFile(instance, private, index, username, !prompt)
		if err != nil {
			return err
		}
		args = append(args, "--temporary", filepath.Dir(file), "--", client, "-c", file)
	} else {
		ipAddress, err := instance.IpAddress(private, index)
		if err != nil {
			return err
		}

		args = append(args, "--", client, "-u", username, ipAddress)
		if !prompt {
			args = append(args, "-p", "-")
		}
	}

	args = append(args, options...)
	args = append(args, arguments...)

	executable, err := rdpFindRunExecutable()
	if err != nil {
		return err
	}

	command := exec.Command(executable, args...)
	if !prompt && !rdpIsRemmina(client) {
		command.Stdin = strings.NewReader(instance.AdminPassword)
	}
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Start()
	if err != nil {
		return err
	}

	fmt.Println(command, command.Stdin)

	err = command.Process.Release()
	if err != nil {
		return err
	}

	return nil
}

func rdpFindClientNative() (string, error) {
	executables := []string{"remmina", "rdesktop"}
	for _, executable := range executables {
		_, err := exec.LookPath(executable)
		if err == nil {
			return executable, nil
		}
	}
	return "", fmt.Errorf("Error finding Remote Desktop client executable: none of %q found in $PATH", executables)
}
