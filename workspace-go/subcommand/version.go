package subcommand

import (
	"fmt"
	"runtime/debug"

	"github.com/ryo-kagawa/go-utils/commandline"
)

var (
	name    = "name"
	version = "develop"
)

type Version struct{}

var _ = (commandline.SubCommand)(Version{})

func (Version) Execute([]string) (string, error) {
	result := ""
	info, _ := debug.ReadBuildInfo()
	result += fmt.Sprintf("%s %s\n", name, version)
	result += fmt.Sprintf("%s\n", info.GoVersion)
	for _, dep := range info.Deps {
		result += fmt.Sprintf("%s %s\n", dep.Path, dep.Version)
	}
	return result, nil
}

func (Version) Name() string {
	return "version"
}
