package model

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

type ImageData struct {
	image image.Image
	_     struct{}
}

func NewImageData(
	image image.Image,
) ImageData {
	return ImageData{
		image: image,
	}
}

func (i ImageData) GetImage() image.Image {
	return i.image
}

func (i ImageData) Resize(w, h uint64) ImageData {
	var result *image.RGBA = image.NewRGBA(
		image.Rectangle{
			Min: image.Point{},
			Max: image.Point{
				X: int(w),
				Y: int(h),
			},
		},
	)

	var ratioW float64 = float64(w) / float64(i.image.Bounds().Dx())
	var ratioH float64 = float64(h) / float64(i.image.Bounds().Dy())
	var ratio float64 = min(ratioW, ratioH)
	var dx uint64 = min(uint64(math.Ceil(float64(i.image.Bounds().Dx())*ratio)), w)
	var dy uint64 = min(uint64(math.Ceil(float64(i.image.Bounds().Dy())*ratio)), h)
	var offsetStart image.Point = image.Point{
		X: int((w - dx) / 2),
		Y: int((h - dy) / 2),
	}

	draw.CatmullRom.Scale(
		result,
		image.Rectangle{
			Min: offsetStart,
			Max: image.Point{
				X: offsetStart.X + int(dx),
				Y: offsetStart.Y + int(dy),
			},
		},
		i.image,
		i.image.Bounds(),
		draw.Over,
		nil,
	)
	return NewImageData(result)
}
