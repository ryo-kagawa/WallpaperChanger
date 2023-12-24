package main

import (
	"bytes"
	"fmt"
	"image"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/ryo-kagawa/WallpaperChanger/configs"
	"github.com/ryo-kagawa/WallpaperChanger/constants"
	"github.com/ryo-kagawa/WallpaperChanger/model"
	"github.com/ryo-kagawa/WallpaperChanger/utils"
	"github.com/ryo-kagawa/WallpaperChanger/utils/window"
)

const configFileName = "config.yaml"

func main() {
	exeFileDirectory, err := utils.GetExeFileDirectory()
	if err != nil {
		fmt.Println(err)
		return
	}
	config, err := configs.LoadConfig(filepath.Join(exeFileDirectory, configFileName))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 対象となる画像ファイルパス一覧
	var filePathList []string = []string{}
	filepath.WalkDir(
		config.ImagePath,
		func(path string, d fs.DirEntry, err error) error {
			// ディレクトリ判定は省略できない
			if d.IsDir() {
				return nil
			}
			if !constants.ImageExtensionList.Includes(filepath.Ext(path)) {
				fmt.Println("not registers image extension: " + path)
				return nil
			}
			filePathList = append(filePathList, path)
			return nil
		},
	)

	// 壁紙に配置する画像を生成
	var imageList []model.ImageData = make([]model.ImageData, 0, len(config.RectangleList))
	for i := 0; i < len(config.RectangleList); i++ {
		targetFilePath := filePathList[uint64(rand.Int63n(int64(len(filePathList))))]
		buffer, err := os.ReadFile(targetFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 事前にチェック済みなのでok判定は不要
		decoder, _ := constants.ImageExtensionList.Find(filepath.Ext(targetFilePath))
		imageData, err := decoder.Decode(bytes.NewBuffer(buffer))
		if err != nil {
			fmt.Println(err)
			return
		}
		imageData = imageData.Resize(config.RectangleList[i].Width, config.RectangleList[i].Height)
		imageList = append(imageList, imageData)
	}

	// 最終的な画像サイズ
	var width uint64
	var height uint64
	for _, x := range config.RectangleList {
		width = max(x.X+x.Width, width)
		height = max(x.Y+x.Height, height)
	}

	// ファイル生成
	var resultImage *image.RGBA = image.NewRGBA(
		image.Rectangle{
			Min: image.Point{},
			Max: image.Pt(int(width), int(height)),
		},
	)

	// 壁紙にする画像を作成
	for i, x := range imageList {
		image := x.GetImage().(*image.RGBA)
		for imageX := 0; imageX < image.Rect.Dx(); imageX++ {
			for imageY := 0; imageY < image.Rect.Dy(); imageY++ {
				offset := (int(config.RectangleList[i].Y)+imageY)*resultImage.Stride + (int(config.RectangleList[i].X)+imageX)*4
				imageOffset := imageY*image.Stride + imageX*4
				resultImage.Pix[offset] = image.Pix[imageOffset]
				resultImage.Pix[offset+1] = image.Pix[imageOffset+1]
				resultImage.Pix[offset+2] = image.Pix[imageOffset+2]
				// アルファ値を0xFFとすることで24bit画像として出力されるようにする
				resultImage.Pix[offset+3] = 0xFF
			}
		}
	}

	err = window.SetWallPaper(resultImage)
	if err != nil {
		fmt.Println(err)
	}
}
