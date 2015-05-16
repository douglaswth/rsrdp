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
	"testing"

	. "github.com/onsi/gomega"
)

func TestNonexistentConfig(t *testing.T) {
	RegisterTestingT(t)

	nonexistentConfigFile := "nonexistent/.rsrdp.yml"
	configFile = &nonexistentConfigFile

	err := readConfig()
	Expect(err).To(HaveOccurred())
}

func TestReadBadEnvironmentConfig(t *testing.T) {
	RegisterTestingT(t)

	badEnvironmentConfigFile := "test/bad_environment.rsrdp.yml"
	configFile = &badEnvironmentConfigFile

	err := readConfig()
	Expect(err).To(HaveOccurred())
}

func TestReadMissingDefaultEnvironmentConfig(t *testing.T) {
	RegisterTestingT(t)

	missingDefaultEnvironmentConfigFile := "test/missing_default_environment.rsrdp.yml"
	configFile = &missingDefaultEnvironmentConfigFile

	err := readConfig()
	Expect(err).To(MatchError(missingDefaultEnvironmentConfigFile + ": could not find default environment: development"))
}

func TestReadExampleConfig(t *testing.T) {
	RegisterTestingT(t)

	exampleConfigFile := "example/.rsrdp.yml"
	configFile = &exampleConfigFile

	err := readConfig()
	Expect(err).NotTo(HaveOccurred())
	Expect(environments).To(Equal(map[string]Environment{
		"production": {
			Account:      12345,
			Host:         "us-3.rightscale.com",
			RefreshToken: "abcdef1234567890abcdef1234567890abcdef12",
		},
		"staging": {
			Account:      67890,
			Host:         "us-4.rightscale.com",
			RefreshToken: "fedcba0987654321febcba0987654321fedcba09",
		},
	}))
	Expect(defaultEnvironment).To(Equal(Environment{
		Account:      12345,
		Host:         "us-3.rightscale.com",
		RefreshToken: "abcdef1234567890abcdef1234567890abcdef12",
	}))
}
