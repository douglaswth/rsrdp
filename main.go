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
	_ "gopkg.in/rightscale/rsc.v1/rsapi"
	_ "gopkg.in/rightscale/rsc.v1/cm15"
)

var (
	app       = kingpin.New("rsrdp", "Launch Windows Remote Desktop for a RightScale Server, ServerArray, or Instance.")
	config    = app.Flag("config", "Set the config file path.").Short('c').Default(configPath).String()
	account   = app.Flag("account", "Set the RightScale account ID.").Short('a').Int()
	host      = app.Flag("host", "RightScale login endpoint (e.g. 'us-3.rightscale.com')").Short('h').String()
	arguments = app.Flag("argument", "Argument to the Remote Desktop command (specify multiple times for multiple arguments)").Short('A').Strings()
	urls      = app.Arg("url", "RightScale Server, ServerArray, or Instance URL").Required().Strings()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	fmt.Println(*account)
	fmt.Println(*config)
	fmt.Println(*host)
	fmt.Println(*arguments, len(*arguments))
	fmt.Println(*urls, len(*urls))
}
