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
	neturl "net/url"
	"regexp"

	"gopkg.in/rightscale/rsc.v1/cm15"
	"gopkg.in/rightscale/rsc.v1/rsapi"
)

var (
	instanceHref    = regexp.MustCompile("^/api/clouds/[^/]+/instances/[^/]+$")
	serverHref      = regexp.MustCompile("^/api/(?:deployments/[^/]+/)?servers/[^/]+$")
	serverArrayHref = regexp.MustCompile("^/api/(?:deployments/[^/]+/)?server_arrays/[^/]+$")
	instancePage    = regexp.MustCompile("^/acct/([^/]+)/clouds/([^/]+)/instances/([^/]+)$")
	serverPage      = regexp.MustCompile("^/acct/([^/]+)/servers/([^/]+)$")
	serverArrayPage = regexp.MustCompile("^/acct/([^/]+)/server_arrays/([^/]+)$")
)

func urlsToInstances(urls []string, prompt bool) ([]*cm15.Instance, error) {
	instances := make([]*cm15.Instance, 0, len(urls))

	for _, url := range urls {
		parsedUrl, err := neturl.Parse(url)
		if err != nil {
			return nil, fmt.Errorf("Error parsing URL: %s", err)
		}

		switch {
		case instanceHref.MatchString(parsedUrl.Path):
			instance, err := urlGetInstanceFromInstanceHref(parsedUrl.Path, config.environment, prompt)
			if err != nil {
				return nil, err
			}
			instances = append(instances, instance)
		case serverHref.MatchString(parsedUrl.Path):
			instance, err := urlGetInstanceFromServerHref(parsedUrl.Path, config.environment, prompt)
			if err != nil {
				return nil, err
			}
			instances = append(instances, instance)
		case serverArrayHref.MatchString(parsedUrl.Path):
		case instancePage.MatchString(parsedUrl.Path):
		case serverPage.MatchString(parsedUrl.Path):
			/*
				fmt.Println(parsedUrl.Host, parsedUrl.Path, parsedUrl.RawQuery)

				query, err := neturl.ParseQuery(parsedUrl.RawQuery)
				instanceId := query.Get("instance_id")
				fmt.Println(query, err, instanceId)
			*/
		case serverArrayPage.MatchString(parsedUrl.Path):
		default:
			return nil, fmt.Errorf("Error parsing URL: %s: unsupported URL format", url)
		}
	}

	fmt.Println(len(instances), cap(instances))

	return instances, nil
}

func urlGetInstanceFromInstanceHref(href string, environment *Environment, prompt bool) (*cm15.Instance, error) {
	client15, err := environment.Client15()
	if err != nil {
		return nil, err
	}

	params := rsapi.ApiParams{}
	if !prompt {
		params["view"] = "sensitive"
	}
	instance, err := client15.InstanceLocator(href).Show(params)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving instance: %s: %s", href, err)
	}

	return instance, nil
}

func urlGetInstanceFromServerHref(href string, environment *Environment, prompt bool) (*cm15.Instance, error) {
	client15, err := environment.Client15()
	if err != nil {
		return nil, err
	}

	server, err := client15.ServerLocator(href).Show(rsapi.ApiParams{})
	if err != nil {
		return nil, fmt.Errorf("Error retrieving server: %s: %s", href, err)
	}

	var currentInstanceHref string
	for _, link := range server.Links {
		if link["rel"] == "current_instance" {
			currentInstanceHref = link["href"]
			break
		}
	}
	if currentInstanceHref == "" {
		return nil, fmt.Errorf("Error retrieving server: %s: server has no current instance", href)
	}

	return urlGetInstanceFromInstanceHref(currentInstanceHref, environment, prompt)
}
