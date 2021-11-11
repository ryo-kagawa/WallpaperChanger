package constants

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/ryo-kagawa/WallpaperChanger/model"
	"github.com/ryo-kagawa/go-utils/conditional"
	"golang.org/x/image/bmp"
)

type imageExtension struct {
	name   string
	decode func(r io.Reader) (image.Image, error)
	_      struct{}
}

func (i imageExtension) Decode(r io.Reader) (model.ImageData, error) {
	img, err := i.decode(r)
	if err != nil {
		return model.ImageData{}, err
	}
	return model.NewImageData(img), nil
}

type imageExtensionList []imageExtension

func (i imageExtensionList) Includes(extension string) bool {
	_, ok := i.Find(extension)
	return ok
}

func (i imageExtensionList) Find(extension string) (imageExtension, bool) {
	extension = conditional.String(extension[0] == '.', extension[1:], extension)
	for _, x := range i {
		if x.name == extension {
			return x, true
		}
	}
	return imageExtension{}, false
}

var ImageExtensionList imageExtensionList = imageExtensionList{
	{
		name:   "bmp",
		decode: bmp.Decode,
	},
	{
		name:   "jpeg",
		decode: jpeg.Decode,
	},
	{
		name:   "jpg",
		decode: jpeg.Decode,
	},
	{
		name:   "png",
		decode: png.Decode,
	},
}
