package main

import (
	"os"
	"strings"

	"github.com/ryo-kagawa/go-utils/commandline"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type Command struct{}

var _ = (commandline.RootCommand)(Command{})

func (Command) Execute([]string) (string, error) {
	if !windows.GetCurrentProcessToken().IsElevated() {
		verb := "runas"
		exe, _ := os.Executable()
		args := strings.Join(os.Args[1:], " ")
		cwd, _ := os.Getwd()

		const SW_SHOWNORMAL int32 = 1

		err := windows.ShellExecute(
			0,
			windows.StringToUTF16Ptr(verb),
			windows.StringToUTF16Ptr(exe),
			windows.StringToUTF16Ptr(args),
			windows.StringToUTF16Ptr(cwd),
			SW_SHOWNORMAL,
		)
		if err != nil {
			return "", err
		}
		return "", nil
	}
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		return "", err
	}
	defer key.Close()
	err = key.SetDWordValue("JPEGImportQuality", 100)
	if err != nil {
		return "", err
	}

	return "", nil
}
