package main

import (
	"os"
	"path/filepath"

	"github.com/ryo-kagawa/WallpaperChanger/configs"
	"github.com/ryo-kagawa/WallpaperChanger/utils"
	"github.com/ryo-kagawa/WallpaperChanger/utils/window"
	"github.com/ryo-kagawa/go-utils/commandline"
	"gopkg.in/yaml.v3"
)

const configFileName = "config.yaml"

type Command struct{}

var _ = (commandline.RootCommand)(Command{})

func (Command) Execute([]string) (string, error) {
	directoryPath, err := window.GetImageDirectoryPath()
	if err != nil {
		return "", err
	}
	rectangleList, err := window.GetMonitorRectangleList()
	if err != nil {
		return "", err
	}
	config := configs.Config{
		ImagePath:     directoryPath,
		RectangleList: rectangleList,
	}

	buf, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	exeFileDirectory, err := utils.GetExeFileDirectory()
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filepath.Join(exeFileDirectory, configFileName), buf, 0777)
	if err != nil {
		return "", err
	}

	return "", nil
}
