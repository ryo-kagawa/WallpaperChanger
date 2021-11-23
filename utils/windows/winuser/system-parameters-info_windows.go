// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/nf-winuser-systemparametersinfow

package winuser

import (
	"unsafe"

	"github.com/ryo-kagawa/WallpaperChanger/utils/windows"
)

var (
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

var uiAction = struct {
	SPI_SETDESKWALLPAPER windows.UINT
}{
	SPI_SETDESKWALLPAPER: 0x0014,
}

var fWinIni = struct {
	SPIF_UPDATEINIFILE windows.UINT
	SPIF_SENDCHANGE    windows.UINT
}{
	SPIF_UPDATEINIFILE: 0x01,
	SPIF_SENDCHANGE:    0x02,
}

func SystemParametersInfo(uiAction windows.UINT, uiParam windows.UINT, pvParam windows.PVOID, fWinIni windows.UINT) (uintptr, uintptr, error) {
	return systemParametersInfo.Call(
		uintptr(uiAction),
		uintptr(uiParam),
		uintptr(pvParam),
		uintptr(fWinIni),
	)
}

func SPI_SETDESKWALLPAPER(filePath string) error {
	fileNameUTF16, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		return err
	}

	_, _, err = SystemParametersInfo(
		uiAction.SPI_SETDESKWALLPAPER,
		0x0000,
		windows.PVOID(uintptr(unsafe.Pointer(fileNameUTF16))),
		fWinIni.SPIF_UPDATEINIFILE|fWinIni.SPIF_SENDCHANGE,
	)
	if err == windows.DS_S_SUCCESS {
		err = nil
	}
	return err
}
