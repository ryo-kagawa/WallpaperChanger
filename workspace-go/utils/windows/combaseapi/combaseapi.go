// https://learn.microsoft.com/ja-jp/windows/win32/api/combaseapi

package combaseapi

import (
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/windef"
	"golang.org/x/sys/windows"
)

var (
	combase = windows.NewLazySystemDLL("combase.dll")
)

var (
	// https://docs.microsoft.com/ja-JP/windows/win32/api/combaseapi/nf-combaseapi-cotaskmemfree
	procCoTaskMemFree = combase.NewProc("CoTaskMemFree")
)

func CoTaskMemFree(pv windef.LPVOID) {
	procCoTaskMemFree.Call(
		uintptr(pv),
	)
}
