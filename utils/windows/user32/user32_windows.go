// Win32APIの参考
// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/nf-winuser-systemparametersinfow

package user32

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

func SystemParametersInfo(filePath string) error {
	fileNameUTF16, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		return err
	}

	systemParametersInfo.Call(
		uiAction.SPI_SETDESKWALLPAPER,
		uintptr(0x0000),
		uintptr(unsafe.Pointer(fileNameUTF16)),
		fWinIni.SPIF_UPDATEINIFILE|fWinIni.SPIF_SENDCHANGE,
	)

	return nil
}
