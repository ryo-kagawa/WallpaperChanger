package window

import (
	"image"
	"os"
	"path/filepath"

	"github.com/ryo-kagawa/WallpaperChanger/configs"
	"github.com/ryo-kagawa/WallpaperChanger/utils"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/combaseapi"
	shobjidlcore "github.com/ryo-kagawa/WallpaperChanger/utils/windows/shobjidl_core"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/shtypes"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/windef"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/winuser"
	"golang.org/x/image/bmp"
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
			HwndOwner: windows.HWND(0),
			LpszTitle: windows.LPCWSTR((*windows.WCHAR)(LpszTitle)),
			UlFlags:   shobjidlcore.BIF_NEWDIALOGSTYLE,
		},
	)
	pidlistAbsolute := shobjidlcore.SHBrowseForFolder(browseInfo)
	_, pszPath := shobjidlcore.SHGetPathFromIDListW(shtypes.PCIDLIST_ABSOLUTE(uintptr(pidlistAbsolute)))
	combaseapi.CoTaskMemFree(windows.LPVOID(pidlistAbsolute))
	value := windows.UTF16PtrToString((*uint16)((*windows.WCHAR)(pszPath)))
	return value, nil
}

func GetMonitorRectangleList() ([]configs.Rectangle, error) {
	rectangleList := []configs.Rectangle{}
	winuser.EnumDisplayMonitors(
		0,
		nil,
		winuser.MONITORENUMPROC(
			func(unnamedParam1 windows.HMONITOR, unnamedParam2 windows.HDC, unnamedPara3 windef.LPRECT, unnamedParam4 windows.LPARAM) windows.BOOL {
				monitorInfoEx, ok := winuser.GetMonitorInfo(unnamedParam1)
				if ok != windows.TRUE {
					return windows.FALSE
				}
				rectangleList = append(rectangleList, configs.Rectangle{
					X:      uint64(monitorInfoEx.RcMonitor.Left),
					Y:      uint64(monitorInfoEx.RcMonitor.Top),
					Width:  uint64(monitorInfoEx.RcMonitor.Right - monitorInfoEx.RcMonitor.Left),
					Height: uint64(monitorInfoEx.RcMonitor.Bottom - monitorInfoEx.RcMonitor.Top),
				})
				return windows.TRUE
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
		defer file.Close()
		if err != nil {
			return err
		}
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

	winuser.SPI_SETDESKWALLPAPER(filePath)

	return nil
}
