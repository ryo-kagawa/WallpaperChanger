package window

import (
	"image"
	"os"
	"path/filepath"

	"github.com/ryo-kagawa/WallpaperChanger/utils/windows/user32"
	"golang.org/x/image/bmp"
)

const (
	outputFilePath = "./wallpaper.bmp"
)

func SetWallPaper(img image.Image) error {
	writeFile := func(img image.Image) error {
		file, err := os.Create(outputFilePath)
		defer file.Close()
		if err != nil {
			return err
		}
		return bmp.Encode(file, img)
	}

	// ファイル出力
	err := writeFile(img)
	if err != nil {
		return err
	}

	filePath, err := filepath.Abs(outputFilePath)
	if err != nil {
		return err
	}

	user32.SystemParametersInfo(filePath)

	return nil
}
