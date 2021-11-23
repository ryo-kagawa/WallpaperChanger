package windows

import (
	"golang.org/x/sys/windows"
)

func NewCallback(fn interface{}) uintptr {
	return windows.NewCallback(fn)
}

func UTF16PtrFromString(s string) (*uint16, error) {
	return windows.UTF16PtrFromString(s)
}
