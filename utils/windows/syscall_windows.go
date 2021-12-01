package windows

import (
	"golang.org/x/sys/windows"
)

func UTF16ToString(s []uint16) string {
	return windows.UTF16ToString(s)
}

func UTF16PtrFromString(s string) (*uint16, error) {
	return windows.UTF16PtrFromString(s)
}

func UTF16PtrToString(p *uint16) string {
	return windows.UTF16PtrToString(p)
}

func NewCallback(fn interface{}) uintptr {
	return windows.NewCallback(fn)
}
