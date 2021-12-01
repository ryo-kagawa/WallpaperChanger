// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/

package winuser

import "github.com/ryo-kagawa/WallpaperChanger/utils/windows"

var (
	user32 = windows.NewLazySystemDLL("user32.dll")
)
