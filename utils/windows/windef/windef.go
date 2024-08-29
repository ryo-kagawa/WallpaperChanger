// https://learn.microsoft.com/ja-jp/windows/win32/api/windef

package windef

import (
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/basetsd"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/winnt"
	"golang.org/x/sys/windows"
)

const (
	FALSE BOOL = 0
	TRUE  BOOL = 1
)

type BOOL int32
type HDC winnt.HANDLE
type HMONITOR winnt.HANDLE
type HWND winnt.HANDLE
type LPARAM basetsd.LONG_PTR
type LPVOID uintptr
type LPCRECT *windows.Rect
type LPRECT *windows.Rect
type UINT uint32
