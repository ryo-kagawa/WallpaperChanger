// https://docs.microsoft.com/en-us/windows/win32/api/windef/ns-windef-rect

package windef

import "github.com/ryo-kagawa/WallpaperChanger/utils/windows"

type RECT struct {
	Left   windows.LONG
	Top    windows.LONG
	Right  windows.LONG
	Bottom windows.LONG
}

type LPCRECT *RECT
type LPRECT *RECT
