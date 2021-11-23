// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/

package winuser

import (
	"golang.org/x/sys/windows"
)

var (
	user32 = windows.NewLazySystemDLL("user32.dll")
)
