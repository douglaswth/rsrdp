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

package win32

import (
	"syscall"
	"unsafe"
)

var (
	FOLDERID_RoamingAppData = syscall.GUID{0x3EB685DB, 0x65F9, 0x4CF6, [8]byte{0xA0, 0x3A, 0xE3, 0xEF, 0x65, 0x72, 0x9F, 0x3D}}
)

var (
	modShell32 = syscall.NewLazyDLL("Shell32.dll")
	procSHGetKnownFolderPath = modShell32.NewProc("SHGetKnownFolderPath")
)

func SHGetKnownFolderPath(rfid *syscall.GUID, dwFlags uint32, hToken syscall.Handle) (string, error) {
	var ppszPath *uint16
	result, _, _ := procSHGetKnownFolderPath.Call(uintptr(unsafe.Pointer(rfid)), uintptr(dwFlags), uintptr(hToken), uintptr(unsafe.Pointer(&ppszPath)))
	if HRESULT(result) != S_OK {
		return "", syscall.Errno(result)
	}
	defer CoTaskMemFree(uintptr(unsafe.Pointer(ppszPath)))

	return syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(ppszPath))[:]), nil
}
