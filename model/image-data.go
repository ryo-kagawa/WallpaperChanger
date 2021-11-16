package model

import (
	"image"
	"math"

	"github.com/ryo-kagawa/go-utils/conditional"
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
	var ratio float64 = conditional.Float64(ratioW < ratioH, ratioW, ratioH)
	var dx uint64 = uint64(math.Min(math.Ceil(float64(i.image.Bounds().Dx())*ratio), float64(w)))
	var dy uint64 = uint64(math.Min(math.Ceil(float64(i.image.Bounds().Dy())*ratio), float64(h)))
	var offsetStart image.Point = image.Point{
		X: int(math.Floor(float64((w - dx) / 2))),
		Y: int(math.Floor(float64((h - dy) / 2))),
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
