// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/nf-winuser-enumdisplaymonitors

package winuser

import (
	"unsafe"

	"github.com/ryo-kagawa/WallpaperChanger/utils/windows"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/windef"
)

var (
	enumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
)

func EnumDisplayMonitors(hdc windows.HDC, lprcClip windef.LPCRECT, lpfnEnum MONITORENUMPROC, dwData windows.LPARAM) windows.BOOL {
	ret, _, _ := enumDisplayMonitors.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(lprcClip)),
		windows.NewCallback(
			func(hmonitor uintptr, hdc uintptr, rect uintptr, lparam uintptr) uintptr {
				return uintptr(
					lpfnEnum(
						windows.HMONITOR(hmonitor),
						windows.HDC(hdc),
						windef.LPRECT(unsafe.Pointer(&rect)),
						windows.LPARAM(lparam),
					),
				)
			},
		),
		uintptr(0),
	)
	return windows.BOOL(ret)
}
