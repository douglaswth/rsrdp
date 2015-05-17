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
	"log"
	"os"

	"gopkg.in/rightscale/rsc.v1/cm15"
	//"gopkg.in/rightscale/rsc.v1/cm16"
	"github.com/rightscale/rsc/cm16"
	"gopkg.in/rightscale/rsc.v1/rsapi"
)

type Environment struct {
	Account      int
	Host         string
	RefreshToken string `mapstructure:"refresh_token"`
	client15     *cm15.Api
	client16     *cm16.Api
}

func (environment *Environment) Client15() (*cm15.Api, error) {
	if environment.client15 == nil {
		auth := rsapi.NewOAuthAuthenticator(environment.RefreshToken)
		var err error
		environment.client15, err = cm15.New(environment.Account, environment.Host, auth, log.New(os.Stdout, "[CM 1.5] ", log.LstdFlags), nil)
		if err != nil {
			return nil, fmt.Errorf("Error initializing RightScale API 1.5 client: %s", err)
		}
	}
	return environment.client15, nil
}

func (environment *Environment) Client16() (*cm16.Api, error) {
	if environment.client16 == nil {
		auth := rsapi.NewOAuthAuthenticator(environment.RefreshToken)
		var err error
		environment.client16, err = cm16.New(environment.Account, environment.Host, auth, log.New(os.Stdout, "[CM 1.6] ", log.LstdFlags), nil)
		if err != nil {
			return nil, fmt.Errorf("Error initializing RightScale API 1.6 client: %s", err)
		}
	}
	return environment.client16, nil
}
