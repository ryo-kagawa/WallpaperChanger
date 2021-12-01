package shobjidlcore

import (
	"unsafe"

	"github.com/ryo-kagawa/WallpaperChanger/utils/windows"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/fileio"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/shtypes"
)

var (
	// https://docs.microsoft.com/ja-jp/windows/win32/api/shlobj_core/nf-shlobj_core-shbrowseforfolderw
	shBrowseForFolder = shell32.NewProc("SHBrowseForFolderW")
	// https://docs.microsoft.com/ja-jp/windows/win32/api/shlobj_core/nf-shlobj_core-shgetpathfromidlistw
	shGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")
)

// https://docs.microsoft.com/ja-jp/previous-versions/windows/desktop/legacy/bb762598(v=vs.85)
type BFFCALLBACK uintptr

func NewBFFCALLBACK(callback func(hwnd windows.HWND, uMsg windows.UINT, lParam windows.LPARAM, lpData windows.LPARAM) int) BFFCALLBACK {
	return BFFCALLBACK(
		windows.NewCallback(
			func(hwnd windows.HWND, uMsg windows.UINT, lParam windows.LPARAM, lpData windows.LPARAM) uintptr {
				return uintptr(callback(hwnd, uMsg, lParam, lpData))
			},
		),
	)
}

// https://docs.microsoft.com/ja-jp/windows/win32/api/shlobj_core/ns-shlobj_core-browseinfow
type BROWSEINFO struct {
	HwndOwner      windows.HWND
	PidlRoot       shtypes.PCIDLIST_ABSOLUTE
	PszDisplayName windows.LPWSTR
	LpszTitle      windows.LPCWSTR
	UlFlags        windows.UINT
	Lpfn           BFFCALLBACK
	LParam         windows.LPARAM
	IImage         int32
}
type LPBROWSEINFO *BROWSEINFO

// https://docs.microsoft.com/ja-jp/windows/win32/api/shlobj_core/ns-shlobj_core-browseinfow
const (
	BIF_NEWDIALOGSTYLE = 0x00000040
)

func SHBrowseForFolder(lpbi LPBROWSEINFO) shtypes.PIDLIST_ABSOLUTE {
	ret, _, _ := shBrowseForFolder.Call(
		uintptr(unsafe.Pointer(lpbi)),
	)

	return shtypes.PIDLIST_ABSOLUTE(ret)
}

func SHGetPathFromIDListW(pidl shtypes.PCIDLIST_ABSOLUTE) (windows.BOOL, windows.LPWSTR) {
	var pszPath [fileio.MAX_PATH]uint16
	ret, _, _ := shGetPathFromIDListW.Call(
		uintptr(pidl),
		uintptr(unsafe.Pointer(windows.LPWSTR((*windows.WCHAR)(&pszPath[0])))),
	)
	return windows.BOOL(ret), windows.LPWSTR((*windows.WCHAR)(&pszPath[0]))
}
