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

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/douglaswth/rsrdp/win32"
)

func rdpLaunchNative(instance *Instance, private bool, index int, arguments []string, prompt bool, username string) error {
	ipAddress, err := instance.IpAddress(private, index)
	if err != nil {
		return err
	}

	count := 3 + len(arguments)
	if !prompt {
		count += 2
	}
	args := make([]string, 0, count)

	if !prompt {
		credential := win32.CREDENTIAL{
			Type:           win32.CRED_TYPE_GENERIC,
			TargetName:     ipAddress,
			Comment:        "Temporary RSRDP credential",
			CredentialBlob: instance.AdminPassword,
			Persist:        win32.CRED_PERSIST_SESSION,
			UserName:       username,
		}
		err = win32.CredWrite(&credential, 0)
		if err != nil {
			return fmt.Errorf("Error storing credential: %s", err)
		}
		args = append(args, "--credential", ipAddress)
	}

	file, err := rdpCreateFile(instance, private, index, username, false)
	if err != nil {
		return err
	}
	args = append(args, "--temporary", filepath.Dir(file), "--", "mstsc", file)
	args = append(args, arguments...)

	executable, err := rdpFindRunExecutable()
	if err != nil {
		return err
	}

	command := exec.Command(executable, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Start()
	if err != nil {
		return err
	}
	err = command.Process.Release()
	if err != nil {
		return err
	}

	return nil
}
