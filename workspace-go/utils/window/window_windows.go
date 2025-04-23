package window

import (
	"errors"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/ryo-kagawa/WallpaperChanger/configs"
	"github.com/ryo-kagawa/WallpaperChanger/utils"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/combaseapi"
	shobjidlcore "github.com/ryo-kagawa/WallpaperChanger/utils/windows/shobjidl_core"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/shtypes"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/windef"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/winnt"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/winuser"
	"golang.org/x/image/bmp"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	outputFileName = "./wallpaper.bmp"
)

func getOutputFilePath() (string, error) {
	exeFileDirectory, err := utils.GetExeFileDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(exeFileDirectory, outputFileName), nil
}

func GetImageDirectoryPath() (string, error) {
	LpszTitle, err := windows.UTF16PtrFromString("壁紙フォルダーを選択してください")
	if err != nil {
		return "", err
	}
	browseInfo := shobjidlcore.LPBROWSEINFO(
		&shobjidlcore.BROWSEINFO{
			HwndOwner: windef.HWND(0),
			LpszTitle: winnt.LPCWSTR((*winnt.WCHAR)(LpszTitle)),
			UlFlags:   shobjidlcore.BIF_NEWDIALOGSTYLE,
		},
	)
	pidlistAbsolute := shobjidlcore.SHBrowseForFolderW(browseInfo)
	var pszPath [windows.MAX_PATH]winnt.WCHAR
	ok := shobjidlcore.SHGetPathFromIDListW(shtypes.PCIDLIST_ABSOLUTE(uintptr(pidlistAbsolute)), winnt.LPWSTR(&pszPath[0]))
	if ok != windef.TRUE {
		return "", errors.New("SHGetPathFromIDListWで失敗しました")
	}
	combaseapi.CoTaskMemFree(windef.LPVOID(pidlistAbsolute))
	value := windows.UTF16PtrToString((*uint16)(&pszPath[0]))
	return value, nil
}

func GetMonitorRectangleList() ([]configs.Rectangle, error) {
	rectangleList := []configs.Rectangle{}
	winuser.EnumDisplayMonitors(
		0,
		nil,
		winuser.MONITORENUMPROC(
			func(unnamedParam1 windef.HMONITOR, unnamedParam2 windef.HDC, unnamedPara3 windef.LPRECT, unnamedParam4 windef.LPARAM) windef.BOOL {
				lpmi := &winuser.MONITORINFOEX{}
				ok := winuser.GetMonitorInfoW(unnamedParam1, lpmi)
				if ok != windef.TRUE {
					return windef.FALSE
				}
				rectangleList = append(rectangleList, configs.Rectangle{
					X:      uint64(lpmi.RcMonitor.Left),
					Y:      uint64(lpmi.RcMonitor.Top),
					Width:  uint64(lpmi.RcMonitor.Right - lpmi.RcMonitor.Left),
					Height: uint64(lpmi.RcMonitor.Bottom - lpmi.RcMonitor.Top),
				})
				return windef.TRUE
			},
		),
		0,
	)

	return rectangleList, nil
}

func SetWallPaper(img image.Image) error {
	outputFilePath, err := getOutputFilePath()
	if err != nil {
		return err
	}
	writeFile := func(img image.Image) error {
		file, err := os.Create(outputFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		return bmp.Encode(file, img)
	}

	// ファイル出力
	err = writeFile(img)
	if err != nil {
		return err
	}

	filePath, err := filepath.Abs(outputFilePath)
	if err != nil {
		return err
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println(`key Control Panel\Desktop is Not Open`)
	}
	defer k.Close()
	value, valtype, err := k.GetIntegerValue("JPEGImportQuality")
	if err != nil {
		fmt.Println(`key Control Panel\Desktop value JPEGImportQuality is Not Get Integer value`)
	}
	if valtype != registry.DWORD || value != 0x00000064 {
		fmt.Println(`key Control Panel\Desktop value JPEGImportQuality is Not DWORD value 0x00000064`)
	}

	fileNameUTF16, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		return err
	}
	ok := winuser.SystemParametersInfoW(
		winuser.SPI_SETDESKWALLPAPER,
		0x0000,
		winnt.PVOID(uintptr(unsafe.Pointer(fileNameUTF16))),
		winuser.SPIF_UPDATEINIFILE|winuser.SPIF_SENDCHANGE,
	)
	if ok != windef.TRUE {
		return errors.New("壁紙の設定に失敗しました")
	}

	return nil
}
