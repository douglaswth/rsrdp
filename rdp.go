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
	"os"
	"time"
)

func rdpLaunch(instance *Instance, private bool, index int, arguments []string, prompt bool, username string, timeout, interval time.Duration) error {
	err := instance.Wait(private, index, prompt, timeout, interval)
	if err != nil {
		return err
	}

	ipAddress, err := instance.IpAddress(private, index)
	if err != nil {
		return err
	}

	rdpWriteParameter(os.Stderr, "full address", ipAddress)
	rdpWriteParameter(os.Stderr, "username", username)
	fmt.Println(instance.AdminPassword)

	return nil
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
