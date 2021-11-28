package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/ryo-kagawa/WallpaperChanger/configs"
	"github.com/ryo-kagawa/WallpaperChanger/constants"
	"github.com/ryo-kagawa/WallpaperChanger/model"
	"github.com/ryo-kagawa/WallpaperChanger/utils"
	"github.com/ryo-kagawa/WallpaperChanger/utils/window"
	"github.com/ryo-kagawa/go-utils/conditional"
	"golang.org/x/image/draw"
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
				return nil
			}
			filePathList = append(filePathList, path)
			return nil
		},
	)

	// 壁紙に配置する画像を生成
	var imageList []model.ImageData = make([]model.ImageData, 0, len(config.RectangleList))
	for i := 0; i < len(config.RectangleList); i++ {
		rand.Seed(time.Now().UnixNano())
		targetFilePath := filePathList[uint64(rand.Int63n(int64(len(filePathList))))]
		buffer, err := ioutil.ReadFile(targetFilePath)
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
		w := x.X + x.Width
		h := x.Y + x.Height
		width = conditional.UInt64(width < w, w, width)
		height = conditional.UInt64(height < h, h, height)
	}

	// ファイル生成
	var resultImage *image.RGBA = image.NewRGBA(
		image.Rectangle{
			Min: image.Point{},
			Max: image.Pt(int(width), int(height)),
		},
	)
	for x := 0; x < resultImage.Rect.Dx(); x++ {
		for y := 0; y < resultImage.Rect.Dy(); y++ {
			resultImage.SetRGBA(
				x,
				y,
				color.RGBA{
					R: 0x00,
					G: 0x00,
					B: 0x00,
					A: 0xFF,
				},
			)
		}
	}
	for i, x := range imageList {
		draw.Draw(
			resultImage,
			image.Rectangle{
				Min: image.Point{
					X: int(config.RectangleList[i].X),
					Y: int(config.RectangleList[i].Y),
				},
				Max: image.Point{
					X: int(config.RectangleList[i].X + config.RectangleList[i].Width),
					Y: int(config.RectangleList[i].Y + config.RectangleList[i].Height),
				},
			},
			x.GetImage(),
			image.Point{},
			draw.Over,
		)
	}

	err = window.SetWallPaper(resultImage)
	if err != nil {
		fmt.Println(err)
	}
}
