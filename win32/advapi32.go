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

// +build windows

package win32

// #include "Windows.h"
// #include "WinCred.h"
import "C"

import (
	"syscall"
	"unsafe"
)

const (
	CRED_FLAGS_PROMPT_NOW      DWORD = 0x2
	CRED_FLAGS_USERNAME_TARGET DWORD = 0x4
)

const (
	CRED_TYPE_GENERIC                 DWORD = 0x1
	CRED_TYPE_DOMAIN_PASSWORD         DWORD = 0x2
	CRED_TYPE_DOMAIN_CERTIFICATE      DWORD = 0x3
	CRED_TYPE_DOMAIN_VISIBLE_PASSWORD DWORD = 0x4
	CRED_TYPE_GENERIC_CERTIFICATE     DWORD = 0x5
	CRED_TYPE_DOMAIN_EXTENDED         DWORD = 0x6
	CRED_TYPE_MAXIMUM                 DWORD = 0x7
	CRED_TYPE_MAXIMUM_EX              DWORD = CRED_TYPE_MAXIMUM + 1000
)

const (
	CRED_PERSIST_SESSION       DWORD = 0x1
	CRED_PERSIST_LOCAL_MACHINE DWORD = 0x2
	CRED_PERSIST_ENTERPRISE    DWORD = 0x3
)

var (
	modAdvapi32    = syscall.NewLazyDLL("Advapi32.dll")
	procCredWriteW = modAdvapi32.NewProc("CredWriteW")
)

type CREDENTIAL struct {
	Flags          DWORD
	Type           DWORD
	TargetName     string
	Comment        string
	LastWritten    FILETIME
	CredentialBlob string
	Persist        DWORD
	Attributes     []CREDENTIAL_ATTRIBUTE
	TargetAlias    string
	UserName       string
}

type CREDENTIAL_ATTRIBUTE struct {
	Keyword string
	Flags   DWORD
	Value   []byte
}

func CredWrite(credential *CREDENTIAL, flags DWORD) error {
	targetName, err := syscall.UTF16PtrFromString(credential.TargetName)
	if err != nil {
		return err
	}

	comment, err := syscall.UTF16PtrFromString(credential.Comment)
	if err != nil {
		return err
	}

	credentialBlob, err := syscall.UTF16PtrFromString(credential.CredentialBlob)
	if err != nil {
		return err
	}

	attributes := make([]C.CREDENTIAL_ATTRIBUTEW, len(credential.Attributes))
	for _, attribute := range credential.Attributes {
		keyword, err := syscall.UTF16PtrFromString(attribute.Keyword)
		if err != nil {
			return err
		}
		attributes = append(attributes, C.CREDENTIAL_ATTRIBUTEW{
			Keyword:   C.LPWSTR(unsafe.Pointer(keyword)),
			Flags:     C.DWORD(attribute.Flags),
			ValueSize: C.DWORD(len(attribute.Value)),
			Value:     C.LPBYTE(unsafe.Pointer(&attribute.Value[0])),
		})
	}

	targetAlias, err := syscall.UTF16PtrFromString(credential.TargetAlias)
	if err != nil {
		return err
	}

	userName, err := syscall.UTF16PtrFromString(credential.UserName)
	if err != nil {
		return err
	}

	credentialw := C.CREDENTIALW{
		Flags:      C.DWORD(credential.Flags),
		Type:       C.DWORD(credential.Type),
		TargetName: C.LPWSTR(unsafe.Pointer(targetName)),
		Comment:    C.LPWSTR(unsafe.Pointer(comment)),
		LastWritten: C.FILETIME{
			dwLowDateTime:  C.DWORD(credential.LastWritten.DwLowDateTime),
			dwHighDateTime: C.DWORD(credential.LastWritten.DwHighDateTime),
		},
		CredentialBlobSize: C.DWORD(len(credential.CredentialBlob) * 2),
		CredentialBlob:     C.LPBYTE(unsafe.Pointer(credentialBlob)),
		Persist:            C.DWORD(credential.Persist),
		AttributeCount:     C.DWORD(len(attributes)),
		TargetAlias:        C.LPWSTR(unsafe.Pointer(targetAlias)),
		UserName:           C.LPWSTR(unsafe.Pointer(userName)),
	}
	if len(attributes) != 0 {
		credentialw.Attributes = &attributes[0]
	}

	result, _, err := procCredWriteW.Call(uintptr(unsafe.Pointer(&credentialw)), uintptr(flags))
	if result == 0 {
		return err
	}
	return nil
}
