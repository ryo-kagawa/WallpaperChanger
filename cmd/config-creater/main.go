package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ryo-kagawa/WallpaperChanger/configs"
	"github.com/ryo-kagawa/WallpaperChanger/utils"
	"github.com/ryo-kagawa/WallpaperChanger/utils/window"
	"gopkg.in/yaml.v3"
)

const configFileName = "config.yaml"

func main() {
	directoryPath, err := window.GetImageDirectoryPath()
	if err != nil {
		fmt.Print(err)
		return
	}
	rectangleList, err := window.GetMonitorRectangleList()
	if err != nil {
		fmt.Print(err)
		return
	}
	config := configs.Config{
		ImagePath:     directoryPath,
		RectangleList: rectangleList,
	}

	buf, err := yaml.Marshal(config)
	if err != nil {
		fmt.Print(err)
		return
	}

	exeFileDirectory, err := utils.GetExeFileDirectory()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(filepath.Join(exeFileDirectory, configFileName), buf, 0777)
	if err != nil {
		fmt.Print(err)
		return
	}
}
