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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kardianos/osext"
)

func rdpLaunch(instance *Instance, private bool, index int, arguments []string, prompt bool, username string, timeout, interval time.Duration) error {
	err := instance.Wait(private, index, prompt, timeout, interval)
	if err != nil {
		return err
	}

	return rdpLaunchNative(instance, private, index, arguments, prompt, username)
}

func rdpCreateFile(instance *Instance, private bool, index int, username string, password bool) (string, error) {
	ipAddress, err := instance.IpAddress(private, index)
	if err != nil {
		return "", err
	}

	dir, err := ioutil.TempDir("", "rsrdp")
	if err != nil {
		return "", fmt.Errorf("Error creating RDP directory: %s", err)
	}

	file, err := os.OpenFile(filepath.Join(dir, ipAddress+".rdp"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", fmt.Errorf("Error creating RDP file: %s", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = rdpWriteParameter(file, "full address", ipAddress)
	if err != nil {
		return "", err
	}
	_, err = rdpWriteParameter(file, "username", username)
	if err != nil {
		return "", err
	}
	if password {
		_, err = rdpWriteParameter(file, "password", instance.AdminPassword)
		if err != nil {
			return "", err
		}
	}

	return file.Name(), nil
}

func rdpWriteParameter(writer io.Writer, key string, value interface{}) (int, error) {
	switch value.(type) {
	case int:
		return fmt.Fprintf(writer, "%s:i:%d\r\n", key, value)
	case string:
		return fmt.Fprintf(writer, "%s:s:%s\r\n", key, value)
	default:
		return fmt.Fprintf(writer, "%s:b:%s\r\n", key, value)
	}
}

func rdpFindRunExecutable() (string, error) {
	folder, err := osext.ExecutableFolder()
	if err != nil {
		return "", fmt.Errorf("Error finding rsrdp-run executable: %s", err)
	}
	executable, err := exec.LookPath(filepath.Join(folder, "rsrdp-run", "rsrdp-run"))
	if err == nil {
		return executable, nil
	}

	switch err := err.(type) {
	case *exec.Error:
		if !os.IsNotExist(err.Err) {
			return "", fmt.Errorf("Error finding rsrdp-run executable: %s", err)
		}
	default:
		return "", fmt.Errorf("Error finding rsrdp-run executable: %s", err)
	}

	executable, err = exec.LookPath("rsrdp-run")
	if err != nil {
		return "", fmt.Errorf("Error finding rsrdp-run executable: %s", err)
	}
	return executable, nil
}
