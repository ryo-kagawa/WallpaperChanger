package main

import (
	"bytes"
	"fmt"
	"image"
	"io/fs"
	"math"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/ryo-kagawa/WallpaperChanger/configs"
	"github.com/ryo-kagawa/WallpaperChanger/constants"
	"github.com/ryo-kagawa/WallpaperChanger/utils"
	"github.com/ryo-kagawa/WallpaperChanger/utils/window"
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
				fmt.Println("not registers image extension: " + path)
				return nil
			}
			filePathList = append(filePathList, path)
			return nil
		},
	)

	resultImage := createResultBaseImage(config.RectangleList)
	// アルファ値を0xFFとすることで24bit画像として出力されるようにする
	for i := 3; i < len(resultImage.Pix); i += 4 {
		resultImage.Pix[i] = 0xFF
	}

	// 壁紙に配置する画像を生成
	for _, x := range config.RectangleList {
		targetFilePath := filePathList[rand.Intn(len(filePathList))]
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

		ratio := min(
			float64(x.Width)/float64(imageData.Bounds().Dx()),
			float64(x.Height)/float64(imageData.Bounds().Dy()),
		)
		dx := min(uint64(math.Ceil(float64(imageData.Bounds().Dx())*ratio)), x.Width)
		dy := min(uint64(math.Ceil(float64(imageData.Bounds().Dy())*ratio)), x.Height)
		offsetPoint := image.Pt(int((x.Width-uint64(dx))/2), int((x.Height-uint64(dy))/2))
		startPoint := image.Pt(int(x.X), int(x.Y)).Add(offsetPoint)

		draw.CatmullRom.Scale(
			resultImage,
			image.Rectangle{
				Min: startPoint,
				Max: startPoint.Add(image.Pt(int(dx), int(dy))),
			},
			imageData,
			imageData.Bounds(),
			draw.Over,
			nil,
		)
	}

	err = window.SetWallPaper(resultImage)
	if err != nil {
		fmt.Println(err)
	}
}

func createResultBaseImage(rectangleList []configs.Rectangle) *image.RGBA {
	// 最終的な画像サイズ
	var width uint64
	var height uint64
	for _, x := range rectangleList {
		width = max(x.X+x.Width, width)
		height = max(x.Y+x.Height, height)
	}

	return image.NewRGBA(
		image.Rectangle{
			Min: image.Point{},
			Max: image.Pt(int(width), int(height)),
		},
	)
}
