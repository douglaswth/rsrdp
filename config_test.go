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

var (
	nonexistentConfigFile               = "nonexistent/.rsrdp.yml"
	badEnvironmentConfigFile            = "test/bad_environment.rsrdp.yml"
	missingDefaultEnvironmentConfigFile = "test/missing_default_environment.rsrdp.yml"
	exampleConfigFile                   = "example/.rsrdp.yml"
)

func TestReadConfigWithNonexistentFile(t *testing.T) {
	RegisterTestingT(t)

	err := readConfig(nonexistentConfigFile, "")
	Expect(err).To(HaveOccurred())
}

func TestReadConfigWithBadEnvironment(t *testing.T) {
	RegisterTestingT(t)

	err := readConfig(badEnvironmentConfigFile, "")
	Expect(err).To(HaveOccurred())
}

func TestReadConfigWithMissingDefaultEnvironment(t *testing.T) {
	RegisterTestingT(t)

	err := readConfig(missingDefaultEnvironmentConfigFile, "")
	Expect(err).To(MatchError(missingDefaultEnvironmentConfigFile + ": could not find default environment: development"))
}

func TestReadConfigWithExample(t *testing.T) {
	RegisterTestingT(t)

	err := readConfig(exampleConfigFile, "")
	Expect(err).NotTo(HaveOccurred())
	Expect(config.environments).To(Equal(map[string]*Environment{
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
	Expect(config.environment).To(Equal(&Environment{
		Account:      12345,
		Host:         "us-3.rightscale.com",
		RefreshToken: "abcdef1234567890abcdef1234567890abcdef12",
	}))
}

func TestReadConfigWithExampleAndEnvironment(t *testing.T) {
	RegisterTestingT(t)

	err := readConfig(exampleConfigFile, "staging")
	Expect(err).NotTo(HaveOccurred())
	Expect(config.environments).To(Equal(map[string]*Environment{
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
	Expect(config.environment).To(Equal(&Environment{
		Account:      67890,
		Host:         "us-4.rightscale.com",
		RefreshToken: "fedcba0987654321febcba0987654321fedcba09",
	}))
}

func TestReadConfigWithExampleAndMissingEnvironment(t *testing.T) {
	RegisterTestingT(t)

	err := readConfig(exampleConfigFile, "development")
	Expect(err).To(MatchError(exampleConfigFile + ": could not find environment: development"))
}
