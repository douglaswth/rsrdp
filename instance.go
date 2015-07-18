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
	"time"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/rightscale/rsc.v3/cm15"
)

type Instance struct {
	*cm15.Instance
	*Environment
}

func (instance *Instance) Href() string {
	for _, link := range instance.Links {
		if link["rel"] == "self" {
			return link["href"]
		}
	}

	panic(fmt.Errorf("No self href for instance: links %s", instance.Links))
}

func (instance *Instance) IpAddress(private bool, index int) (string, error) {
	ipAddresses := instance.PublicIpAddresses
	if private {
		ipAddresses = instance.PrivateIpAddresses
	}

	if index < 0 || index >= len(ipAddresses) {
		return "", fmt.Errorf("Interface index out of bounds: %d: instance %s %s", index, instance.Href(), ipAddresses)
	}

	return ipAddresses[index], nil
}

func (instance *Instance) Wait(private bool, index int, prompt bool, timeout, interval time.Duration) error {
	errChan := make(chan error, 1)
	go func() {
		for {
			_, err := instance.IpAddress(private, 0)
			if err == nil && (prompt || instance.AdminPassword != "") {
				errChan <- nil
				return
			}

			log15.Info("waiting for IP address and/or Administrator password", "instance", instance.Href(), "interval", interval)
			time.Sleep(interval)

			newInstance, err := urlGetInstanceFromInstanceHref(instance.Href(), instance.Environment, prompt)
			if err != nil {
				errChan <- err
				return
			}
			instance.Instance = newInstance.Instance
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("Timeout waiting for IP address and/or Administrator password: %s: %s", timeout, instance.Href())
	}
}
