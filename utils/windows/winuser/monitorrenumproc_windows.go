// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-monitorenumproc

package winuser

import (
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/windef"
)

type MONITORENUMPROC func(unnamedParam1 windows.HMONITOR, unnamedParam2 windows.HDC, unnamedPara3 windef.LPRECT, unnamedParam4 windows.LPARAM) windows.BOOL
