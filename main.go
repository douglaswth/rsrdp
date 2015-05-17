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

	"gopkg.in/alecthomas/kingpin.v1"
	//"github.com/rightscale/rsc/rsapi"
)

var (
	app         = kingpin.New("rsrdp", "Launch Windows Remote Desktop for a RightScale Server, ServerArray, or Instance.")
	configFile  = app.Flag("config", "Set the config file path.").Short('c').Default(configPath).String()
	environment = app.Flag("environment", "Set the RightScale login environment.").Short('e').String()
	account     = app.Flag("account", "Set the RightScale account ID.").Short('a').Int()
	host        = app.Flag("host", "RightScale login endpoint (e.g. 'us-3.rightscale.com')").Short('h').String()
	private     = app.Flag("private", "Connect to the Server, ServerArray, or Instance with the private interface instead of the public interface.").Short('p').Bool()
	index       = app.Flag("index", "Connect using the indexed public/private interface of the Server, ServerArray, or Instance.").Short('i').Int()
	arguments   = app.Flag("argument", "Argument to the Remote Desktop command (specify multiple times for multiple arguments)").Short('A').Strings()
	urls        = app.Arg("url", "RightScale Server, ServerArray, or Instance URL").Required().Strings()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	err := readConfig(*configFile, *environment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %s\n", err)
		os.Exit(1)
	}

	/*
		client16, err := config.environment.Client16()
		if err != nil {
		}
		instances, err := client16.InstanceLocator("/api/clouds/6/instances").Index(rsapi.ApiParams{})
		if err != nil {
		}

		for _, instance := range instances {
			if instance.State != "inactive" {
				fmt.Println(instance.Id, instance.LegacyId, instance.ResourceUid)
			}
		}
	*/
}
