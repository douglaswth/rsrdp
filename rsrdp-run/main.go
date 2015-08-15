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

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	credential = kingpin.Flag("credential", "Temporary RSRDP credential").String()
	executable = kingpin.Arg("executable", "Windows Remote Desktop client executable").Required().String()
	arguments  = kingpin.Arg("arguments", "Arguments to Windows Remote Desktop client").Required().Strings()
)

func main() {
	kingpin.Parse()

	command := exec.Command(*executable, *arguments...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Error running %s: %s", filepath.Base(os.Args[0]), *executable, err)
		os.Exit(1)
	}

	if *credential != "" {
		err = deleteCredential(*credential)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error deleting credential: %s\n", filepath.Base(os.Args[0]), err)
			os.Exit(1)
		}
	}
}
