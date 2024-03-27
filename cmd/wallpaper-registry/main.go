package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func main() {
	if !windows.GetCurrentProcessToken().IsElevated() {
		utf16PtrFromString := func(value string) *uint16 {
			ret, _ := syscall.UTF16PtrFromString(value)
			return ret
		}
		verb := "runas"
		exe, _ := os.Executable()
		args := strings.Join(os.Args[1:], " ")
		cwd, _ := os.Getwd()

		const SW_SHOWNORMAL int32 = 1

		err := windows.ShellExecute(
			0,
			utf16PtrFromString(verb),
			utf16PtrFromString(exe),
			utf16PtrFromString(args),
			utf16PtrFromString(cwd),
			SW_SHOWNORMAL,
		)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		panic(err)
	}
	defer key.Close()
	err = key.SetDWordValue("JPEGImportQuality", 100)
	if err != nil {
		panic(err)
	}
}
