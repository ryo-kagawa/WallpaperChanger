package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/ryo-kagawa/WallpaperChanger/configs"
	"github.com/ryo-kagawa/WallpaperChanger/constants"
	"github.com/ryo-kagawa/WallpaperChanger/model"
	"github.com/ryo-kagawa/go-utils/conditional"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	"golang.org/x/sys/windows"
)

const filePath = "config.yaml"

func main() {
	config, err := configs.LoadConfig(filePath)
	if err != nil {
		fmt.Println(err)
	}

	// 対象となる画像ファイルパス一覧
	var filePathList []string = []string{}
	filepath.Walk(
		config.ImagePath,
		func(path string, info fs.FileInfo, err error) error {
			// ディレクトリ判定は省略できない
			if info.IsDir() {
				return nil
			}
			if !constants.ImageExtensionList.Includes(filepath.Ext(path)) {
				return nil
			}
			filePathList = append(filePathList, path)
			return nil
		},
	)

	var imageList []model.ImageData = make([]model.ImageData, 0, len(config.ImageList))
	for i := 0; i < len(config.ImageList); i++ {
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
		imageData = imageData.Resize(config.ImageList[i].W, config.ImageList[i].H)
		imageList = append(imageList, imageData)
	}

	// 最終的な画像サイズ
	var width uint64
	var height uint64
	for _, x := range config.ImageList {
		w := x.X + x.W
		h := x.Y + x.H
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
					X: int(config.ImageList[i].X),
					Y: int(config.ImageList[i].Y),
				},
				Max: image.Point{
					X: int(config.ImageList[i].X + config.ImageList[i].W),
					Y: int(config.ImageList[i].Y + config.ImageList[i].H),
				},
			},
			x.GetImage(),
			image.Point{},
			draw.Over,
		)
	}

	// ファイル出力
	file, _ := os.Create("./background.bmp")
	defer file.Close()
	err = bmp.Encode(file, resultImage)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Win32APIの参考
	// https://docs.microsoft.com/ja-jp/windows/win32/api/winuser/nf-winuser-systemparametersinfow
	user32, err := windows.LoadDLL("user32.dll")
	if err != nil {
		fmt.Println(err)
		return
	}
	systemParametersInfo, err := user32.FindProc("SystemParametersInfoW")
	if err != nil {
		fmt.Println(err)
		return
	}
	var SPI_SETDESKWALLPAPER int = 0x0014
	var SPIF_UPDATEINIFILE int = 0x01
	// var SPIF_SENDCHANGE int = 0x02
	filePath, err := filepath.Abs(file.Name())
	if err != nil {
		fmt.Println(err)
		return
	}
	fileNameUTF16, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// TODO: 通常の方法とSPIF_UPDATEINIFILEのみの2回行った方が良いか？
	systemParametersInfo.Call(
		uintptr(SPI_SETDESKWALLPAPER),
		uintptr(0x0000),
		uintptr(unsafe.Pointer(fileNameUTF16)),
		// SPIF_SENDCHANGEを指定した場合に壁紙が正常に更新されない
		// uintptr(SPIF_UPDATEINIFILE|SPIF_SENDCHANGE),
		uintptr(SPIF_UPDATEINIFILE),
	)
}
