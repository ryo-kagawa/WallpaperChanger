package combaseapi

import "github.com/ryo-kagawa/WallpaperChanger/utils/windows"

var (
	combase = windows.NewLazySystemDLL("combase.dll")
)

var (
	// https://docs.microsoft.com/ja-JP/windows/win32/api/combaseapi/nf-combaseapi-cotaskmemfree
	coTaskMemFree = combase.NewProc("CoTaskMemFree")
)

func CoTaskMemFree(pv windows.LPVOID) {
	coTaskMemFree.Call(
		uintptr(pv),
	)
}
