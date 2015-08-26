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
	"path/filepath"

	"github.com/mattn/go-colorable"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/rightscale/rsc.v4/log"
)

var (
	app         = kingpin.New("rsrdp", "Launch Windows Remote Desktop for a RightScale Server, ServerArray, or Instance.")
	configFile  = app.Flag("config", "Set the config file path.").Short('c').Default(defaultConfigFile()).String()
	environment = app.Flag("environment", "Set the RightScale login environment.").Short('e').String()
	account     = app.Flag("account", "Set the RightScale account ID.").Short('a').Int()
	host        = app.Flag("host", "RightScale login endpoint (e.g. 'us-3.rightscale.com')").Short('h').String()
	private     = app.Flag("private", "Connect to the Server, ServerArray, or Instance with the private interface instead of the public interface.").Short('p').Bool()
	index       = app.Flag("index", "Connect using the indexed public/private interface of the Server, ServerArray, or Instance.").Short('i').Int()
	arguments   = app.Flag("argument", "Argument to the Remote Desktop command (specify multiple times for multiple arguments)").Short('A').Strings()
	prompt      = app.Flag("prompt", "Prompt for a username and password when launching Windows Remote Desktop rather than using the initial Adminstrator password from RightScale.").Short('P').Bool()
	username    = app.Flag("username", "The username to connect with").Default("Administrator").Short('u').String()
	timeout     = app.Flag("timeout", "The amount to wait for the Server, ServerArray, or Instance to have an IP address and/or Administrator password").Short('t').Default("5m").Duration()
	interval    = app.Flag("interval", "The amount of time between retries when waiting for the Server, ServerArray, or Instance to have an IP address and/or Administrator password").Short('I').Default("10s").Duration()
	urls        = app.Arg("url", "RightScale Server, ServerArray, or Instance URL").Required().Strings()
)

func main() {
	handler := log15.StreamHandler(colorable.NewColorableStdout(), log15.TerminalFormat())
	log15.Root().SetHandler(handler)
	log.Logger.SetHandler(handler)

	app.Writer(os.Stdout)
	kingpin.MustParse(app.Parse(os.Args[1:]))
	err := readConfig(*configFile, *environment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Error reading config file: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}

	instances, err := urlsToInstances(*urls, *prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}

	errChans := make([]chan error, len(instances))
	for errChanIndex, instance := range instances {
		errChans[errChanIndex] = make(chan error)
		go func(errChanIndex int, instance *Instance) {
			errChans[errChanIndex] <- rdpLaunch(instance, *private, *index, *arguments, *prompt, *username, *timeout, *interval)
		}(errChanIndex, instance)
	}

	errs := false
	for _, errChan := range errChans {
		err = <-errChan
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
			errs = true
		}
	}
	if errs {
		os.Exit(1)
	}
}
