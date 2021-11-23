// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/nf-winuser-getmonitorinfow

package winuser

import (
	"unsafe"

	"github.com/ryo-kagawa/WallpaperChanger/utils/windows"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/windef"
)

var (
	getMonitorInfo = user32.NewProc("GetMonitorInfoW")
)

const (
	CCHDEVICENAME int = 32
)

// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/ns-winuser-monitorinfo
type MONITORINFO struct {
	CbSize    windows.DWORD
	RcMonitor windef.RECT
	RcWork    windef.RECT
	DwFlags   windows.DWORD
}

// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/ns-winuser-monitorinfoexw
type MONITORINFOEX struct {
	MONITORINFO
	SzDevice [CCHDEVICENAME]windows.WCHAR
}

func GetMonitorInfo(hMonitor windows.HMONITOR) (MONITORINFOEX, windows.BOOL) {
	monitorInfoEx := MONITORINFOEX{}
	monitorInfoEx.CbSize = windows.DWORD(unsafe.Sizeof(monitorInfoEx))
	ret, _, _ := getMonitorInfo.Call(
		uintptr(hMonitor),
		uintptr(unsafe.Pointer(&monitorInfoEx)),
	)
	return monitorInfoEx, windows.BOOL(ret)
}
