// https://learn.microsoft.com/ja-jp/windows/win32/api/shlobj_core/

package shobjidlcore

import (
	"unsafe"

	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/shtypes"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/windef"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/winnt"
	"golang.org/x/sys/windows"
)

var (
	shell32 = windows.NewLazySystemDLL("shell32.dll")
)

var (
	// https://docs.microsoft.com/ja-jp/windows/win32/api/shlobj_core/nf-shlobj_core-shbrowseforfolderw
	procSHBrowseForFolderW = shell32.NewProc("SHBrowseForFolderW")
	// https://docs.microsoft.com/ja-jp/windows/win32/api/shlobj_core/nf-shlobj_core-shgetpathfromidlistw
	procSHGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")
)

// https://docs.microsoft.com/ja-jp/previous-versions/windows/desktop/legacy/bb762598(v=vs.85)
type BFFCALLBACK uintptr

// https://docs.microsoft.com/ja-jp/windows/win32/api/shlobj_core/ns-shlobj_core-browseinfow
type BROWSEINFO struct {
	HwndOwner      windef.HWND
	PidlRoot       shtypes.PCIDLIST_ABSOLUTE
	PszDisplayName winnt.LPWSTR
	LpszTitle      winnt.LPCWSTR
	UlFlags        windef.UINT
	Lpfn           BFFCALLBACK
	LParam         windef.LPARAM
	IImage         int32
}
type LPBROWSEINFO *BROWSEINFO

// https://docs.microsoft.com/ja-jp/windows/win32/api/shlobj_core/ns-shlobj_core-browseinfow
const (
	BIF_NEWDIALOGSTYLE = 0x00000040
)

func SHBrowseForFolderW(lpbi LPBROWSEINFO) shtypes.PIDLIST_ABSOLUTE {
	ret, _, _ := procSHBrowseForFolderW.Call(
		uintptr(unsafe.Pointer(lpbi)),
	)

	return shtypes.PIDLIST_ABSOLUTE(ret)
}

func SHGetPathFromIDListW(pidl shtypes.PCIDLIST_ABSOLUTE, pszPath winnt.LPWSTR) windef.BOOL {
	ret, _, _ := procSHGetPathFromIDListW.Call(
		uintptr(pidl),
		uintptr(unsafe.Pointer(pszPath)),
	)
	return windef.BOOL(ret)
}
