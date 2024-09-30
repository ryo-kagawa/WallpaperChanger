package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ryo-kagawa/WallpaperChanger/subcommand"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/kernel32"
	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/windef"
	"github.com/ryo-kagawa/go-utils/commandline"
	"golang.org/x/sys/windows"
)

const configFileName = "config.yaml"

func main() {
	result, err := commandline.Execute(
		Command{},
		os.Args[1:],
		subcommand.Version{},
	)
	if result != "" || err != nil {
		if kernel32.AttachConsole(kernel32.ATTACH_PARENT_PROCESS) == windef.TRUE {
			stdoutfd, _ := windows.Open("CONOUT$", windows.O_RDWR, 0)
			if result != "" {
				windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, stdoutfd)
				os.Stdout = os.NewFile(uintptr(stdoutfd), "/dev/stdout")
				fmt.Fprintln(os.Stdout, result)
			}
			if err != nil {
				windows.SetStdHandle(windows.STD_ERROR_HANDLE, stdoutfd)
				os.Stderr = os.NewFile(uintptr(stdoutfd), "/dev/stderr")
				fmt.Fprintln(os.Stderr, err)
			}
		} else {
			kernel32.AllocConsole()
			stdoutfd, _ := windows.Open("CONOUT$", windows.O_RDWR, 0)
			windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, stdoutfd)
			os.Stdout = os.NewFile(uintptr(stdoutfd), "/dev/stdout")
			fmt.Fprintln(os.Stdout, result)
			if result != "" {
				fmt.Fprintln(os.Stdout, result)
			}
			if err != nil {
				windows.SetStdHandle(windows.STD_ERROR_HANDLE, stdoutfd)
				os.Stderr = os.NewFile(uintptr(stdoutfd), "/dev/stderr")
				fmt.Fprintln(os.Stderr, err)
			}
			stdinfd, _ := windows.Open("CONIN$", windows.O_RDWR, 0)
			windows.SetStdHandle(windows.STD_INPUT_HANDLE, stdinfd)
			os.Stdin = os.NewFile(uintptr(stdinfd), "/dev/stdin")
			fmt.Fprintln(os.Stdout, "Enterキーで終了します")
			bufio.NewScanner(os.Stdin).Scan()
		}
	}
}
