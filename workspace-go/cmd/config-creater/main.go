package main

import (
	"fmt"
	"os"

	"github.com/ryo-kagawa/WallpaperChanger/subcommand"
	"github.com/ryo-kagawa/go-utils/commandline"
)

func main() {
	result, err := commandline.Execute(
		Command{},
		os.Args[1:],
		subcommand.Version{},
	)
	if result != "" {
		fmt.Fprint(os.Stdout, result)
	}
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
