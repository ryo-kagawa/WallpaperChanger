package shobjidlcore

import "github.com/ryo-kagawa/WallpaperChanger/utils/windows"

var (
	shell32 = windows.NewLazySystemDLL("shell32.dll")
)
