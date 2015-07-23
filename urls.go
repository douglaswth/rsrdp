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
	"strconv"

	"gopkg.in/rightscale/rsc.v3/rsapi"
)

var (
	instanceHref    = regexp.MustCompile("^/api/clouds/(\\d+)/instances/[^/]+$")
	serverHref      = regexp.MustCompile("^/api/(?:deployments/\\d+/)?servers/\\d+$")
	serverArrayHref = regexp.MustCompile("^/api/(?:deployments/\\d+/)?server_arrays/\\d+$")
	instancePage    = regexp.MustCompile("^/acct/(\\d+)/clouds/(\\d+)/instances/(\\d+)$")
	serverPage      = regexp.MustCompile("^/acct/(\\d+)/servers/(\\d+)$")
	serverArrayPage = regexp.MustCompile("^/acct/(\\d+)/server_arrays/(\\d+)$")
	redirectPage    = regexp.MustCompile("^/acct/(\\d+)/redirect_to_ui_uri$")
)

func urlsToInstances(urls []string, prompt bool) ([]*Instance, error) {
	instances := make([]*Instance, 0, len(urls))

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
			arrayInstances, err := urlGetInstancesFromServerArrayHref(parsedUrl.Path, config.environment, prompt)
			if err != nil {
				return nil, err
			}
			instances = append(instances, arrayInstances...)
		case instancePage.MatchString(parsedUrl.Path):
			instance, err := urlGetInstanceFromInstancePage(parsedUrl, prompt)
			if err != nil {
				return nil, err
			}
			instances = append(instances, instance)
		case serverPage.MatchString(parsedUrl.Path):
			instance, err := urlGetInstanceFromServerPage(parsedUrl, prompt)
			if err != nil {
				return nil, err
			}
			instances = append(instances, instance)
		case serverArrayPage.MatchString(parsedUrl.Path):
			arrayInstances, err := urlGetInstancesFromServerArrayPage(parsedUrl, prompt)
			if err != nil {
				return nil, err
			}
			instances = append(instances, arrayInstances...)
		case redirectPage.MatchString(parsedUrl.Path):
			arrayInstances, err := urlGetInstancesFromRedirectPage(parsedUrl, prompt)
			if err != nil {
				return nil, err
			}
			instances = append(instances, arrayInstances...)
		default:
			return nil, fmt.Errorf("Error parsing URL: %s: unsupported URL format", url)
		}
	}

	return instances, nil
}

func urlGetInstanceFromInstanceHref(href string, environment *Environment, prompt bool) (*Instance, error) {
	client15 := environment.Client15()
	params := rsapi.ApiParams{}
	if !prompt {
		params["view"] = "sensitive"
	}
	instance, err := client15.InstanceLocator(href).Show(params)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving instance: %s: %s", href, err)
	}

	return &Instance{instance, environment}, nil
}

func urlGetInstanceFromServerHref(href string, environment *Environment, prompt bool) (*Instance, error) {
	client15 := environment.Client15()
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

func urlGetInstancesFromServerArrayHref(href string, environment *Environment, prompt bool) ([]*Instance, error) {
	client15 := environment.Client15()
	array, err := client15.ServerArrayLocator(href).Show(rsapi.ApiParams{})
	if err != nil {
		return nil, fmt.Errorf("Error retrieving array: %s: %s", href, err)
	}

	var currentInstancesHref string
	for _, link := range array.Links {
		if link["rel"] == "current_instances" {
			currentInstancesHref = link["href"]
			break
		}
	}

	params := rsapi.ApiParams{}
	if !prompt {
		params["view"] = "sensitive"
	}
	currentInstances, err := client15.InstanceLocator(currentInstancesHref).Index(params)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving array instances: %s: %s", currentInstancesHref, err)
	}

	instances := make([]*Instance, len(currentInstances))
	for index, instance := range currentInstances {
		instances[index] = &Instance{instance, environment}
	}

	return instances, nil
}

func urlGetInstanceFromInstancePage(url *neturl.URL, prompt bool) (*Instance, error) {
	submatches := instancePage.FindStringSubmatch(url.Path)
	account, _ := strconv.ParseInt(submatches[1], 0, 0)
	cloud, _ := strconv.ParseInt(submatches[2], 0, 0)
	legacyId, _ := strconv.ParseInt(submatches[3], 0, 0)

	environment, err := config.getEnvironment(int(account), url.Host)
	if err != nil {
		return nil, err
	}

	return urlGetInstanceFromLegacyId(int(cloud), int(legacyId), environment, prompt)
}

func urlGetInstanceFromServerPage(url *neturl.URL, prompt bool) (*Instance, error) {
	submatches := serverPage.FindStringSubmatch(url.Path)
	account, _ := strconv.ParseInt(submatches[1], 0, 0)
	href := "/api/servers/" + submatches[2]

	environment, err := config.getEnvironment(int(account), url.Host)
	if err != nil {
		return nil, err
	}

	instanceId := url.Query().Get("instance_id")
	if instanceId != "" {
		client15 := environment.Client15()
		server, err := client15.ServerLocator(href).Show(rsapi.ApiParams{})
		if err != nil {
			return nil, fmt.Errorf("Error retrieving server: %s: %s", href, err)
		}

		var nextInstanceHref string
		for _, link := range server.Links {
			if link["rel"] == "next_instance" {
				nextInstanceHref = link["href"]
				break
			}
		}
		if nextInstanceHref == "" {
			return nil, fmt.Errorf("Error retrieving server: %s: server has no next instance", href)
		}

		submatches := instanceHref.FindStringSubmatch(nextInstanceHref)
		cloud, _ := strconv.ParseInt(submatches[1], 0, 0)
		legacyId, err := strconv.ParseInt(instanceId, 0, 0)
		if err != nil {
			return nil, err
		}

		return urlGetInstanceFromLegacyId(int(cloud), int(legacyId), environment, prompt)
	}

	return urlGetInstanceFromServerHref(href, environment, prompt)
}

func urlGetInstancesFromServerArrayPage(url *neturl.URL, prompt bool) ([]*Instance, error) {
	submatches := serverArrayPage.FindStringSubmatch(url.Path)
	account, _ := strconv.ParseInt(submatches[1], 0, 0)
	href := "/api/server_arrays/" + submatches[2]

	environment, err := config.getEnvironment(int(account), url.Host)
	if err != nil {
		return nil, err
	}

	return urlGetInstancesFromServerArrayHref(href, environment, prompt)
}

func urlGetInstancesFromRedirectPage(url *neturl.URL, prompt bool) ([]*Instance, error) {
	submatches := redirectPage.FindStringSubmatch(url.Path)
	account, _ := strconv.ParseInt(submatches[1], 0, 0)

	environment, err := config.getEnvironment(int(account), url.Host)
	if err != nil {
		return nil, err
	}

	query := url.Query()
	resourceType := query.Get("resource_type")
	resourceUri := query.Get("resource_uri")
	instances := make([]*Instance, 1)

	switch resourceType {
	case "instance":
		instances[0], err = urlGetInstanceFromInstanceHref(resourceUri, environment, prompt)
		if err != nil {
			return nil, err
		}
	case "server":
		instances[0], err = urlGetInstanceFromServerHref(resourceUri, environment, prompt)
		if err != nil {
			return nil, err
		}
	case "server_array":
		return urlGetInstancesFromServerArrayHref(resourceUri, environment, prompt)
	default:
		return nil, fmt.Errorf("Error parsing URL: %s: unsupported resource type: %s", url, resourceType)
	}

	return instances, nil
}

func urlGetInstanceFromLegacyId(cloud, legacyId int, environment *Environment, prompt bool) (*Instance, error) {
	client16 := environment.Client16()
	instances, err := client16.InstanceLocator(fmt.Sprintf("/api/clouds/%d/instances", cloud)).Index(rsapi.ApiParams{})
	if err != nil {
		return nil, err
	}

	// TODO: remove print and uncomment loop when RSC and CM1.6 work correctly for collections
	fmt.Println(instances)
	/*for _, instance := range instances {
		if instance.LegacyId == legacyId {
			return urlGetInstanceFromInstanceHref(instance.Href, environment, prompt)
		}
	}*/

	return nil, fmt.Errorf("Could not find instance with legacy ID: %d", legacyId)
}
