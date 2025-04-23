package kernel32

import (
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/intsafe"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/windef"
	"golang.org/x/sys/windows"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
)

var (
	// https://learn.microsoft.com/ja-jp/windows/console/allocconsole
	procAllocConsole = kernel32.NewProc("AllocConsole")
	// https://learn.microsoft.com/ja-jp/windows/console/attachconsole
	procAttachConsole = kernel32.NewProc("AttachConsole")
)

var (
	ATTACH_PARENT_PROCESS intsafe.DWORD = -1 & (1<<32 - 1)
)

func AllocConsole() windef.BOOL {
	ret, _, _ := procAllocConsole.Call()
	return windef.BOOL(ret)
}

func AttachConsole(dwProcessId intsafe.DWORD) windef.BOOL {
	ret, _, _ := procAttachConsole.Call(uintptr(dwProcessId))
	return windef.BOOL(ret)
}
