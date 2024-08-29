// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/

package winuser

import (
	"unsafe"

	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/intsafe"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/windef"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/winnt"
	"golang.org/x/sys/windows"
)

var (
	user32 = windows.NewLazySystemDLL("user32.dll")
)

var (
	// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/nf-winuser-enumdisplaymonitors
	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/nf-winuser-getmonitorinfow
	procGetMonitorInfoW = user32.NewProc("GetMonitorInfoW")
	// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/nf-winuser-systemparametersinfow
	procSystemParametersInfoW = user32.NewProc("SystemParametersInfoW")
)

const (
	CCHDEVICENAME int = 32
)

// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/ns-winuser-monitorinfo
type MONITORINFO struct {
	CbSize    intsafe.DWORD
	RcMonitor windows.Rect
	RcWork    windows.Rect
	DwFlags   intsafe.DWORD
}

// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/ns-winuser-monitorinfoexw
type MONITORINFOEX struct {
	MONITORINFO
	SzDevice [CCHDEVICENAME]winnt.WCHAR
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-monitorenumproc
type MONITORENUMPROC func(unnamedParam1 windef.HMONITOR, unnamedParam2 windef.HDC, unnamedPara3 windef.LPRECT, unnamedParam4 windef.LPARAM) windef.BOOL

const SPI_SETDESKWALLPAPER = 0x0014

const (
	SPIF_UPDATEINIFILE = 0x01
	SPIF_SENDCHANGE    = 0x02
)

func EnumDisplayMonitors(hdc windef.HDC, lprcClip windef.LPCRECT, lpfnEnum MONITORENUMPROC, dwData windef.LPARAM) windef.BOOL {
	ret, _, _ := procEnumDisplayMonitors.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(lprcClip)),
		windows.NewCallback(
			func(hmonitor uintptr, hdc uintptr, rect uintptr, lparam uintptr) uintptr {
				return uintptr(
					lpfnEnum(
						windef.HMONITOR(hmonitor),
						windef.HDC(hdc),
						windef.LPRECT(unsafe.Pointer(&rect)),
						windef.LPARAM(lparam),
					),
				)
			},
		),
		uintptr(0),
	)
	return windef.BOOL(ret)
}

func GetMonitorInfoW(hMonitor windef.HMONITOR, lpmi *MONITORINFOEX) windef.BOOL {
	lpmi.CbSize = intsafe.DWORD(unsafe.Sizeof(*lpmi))
	ret, _, _ := procGetMonitorInfoW.Call(
		uintptr(hMonitor),
		uintptr(unsafe.Pointer(lpmi)),
	)
	return windef.BOOL(ret)
}

func SystemParametersInfoW(uiAction windef.UINT, uiParam windef.UINT, pvParam winnt.PVOID, fWinIni windef.UINT) windef.BOOL {
	ret, _, _ := procSystemParametersInfoW.Call(
		uintptr(uiAction),
		uintptr(uiParam),
		uintptr(pvParam),
		uintptr(fWinIni),
	)
	return windef.BOOL(ret)
}
