package windows

import "golang.org/x/sys/windows"

func NewLazySystemDLL(name string) *windows.LazyDLL {
	return windows.NewLazySystemDLL(name)
}
