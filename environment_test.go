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
	"testing"

	. "github.com/onsi/gomega"
)

var testingEnvironment = Environment{
	Account:      54321,
	Host:         "localhost",
	RefreshToken: "def1234567890abcdef1234567890abcdef12345",
}

func TestEnvironmentClient15(t *testing.T) {
	RegisterTestingT(t)

	firstClient := testingEnvironment.Client15()
	Expect(firstClient).NotTo(BeNil())
	Expect(fmt.Sprintln(firstClient)).To(Equal(fmt.Sprintln(testingEnvironment.client15)))

	secondClient := testingEnvironment.Client15()
	Expect(fmt.Sprintln(secondClient)).To(Equal(fmt.Sprintln(firstClient)))
}

func TestEnvironmentClient16(t *testing.T) {
	RegisterTestingT(t)

	firstClient := testingEnvironment.Client16()
	Expect(firstClient).NotTo(BeNil())
	Expect(fmt.Sprintln(firstClient)).To(Equal(fmt.Sprintln(testingEnvironment.client16)))

	secondClient := testingEnvironment.Client16()
	Expect(fmt.Sprintln(secondClient)).To(Equal(fmt.Sprintln(firstClient)))
}
