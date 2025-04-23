package constants

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"golang.org/x/image/bmp"
)

type imageExtension struct {
	name          string
	extensionList []string
	decode        func(r io.Reader) (image.Image, error)
	_             struct{}
}

func (i imageExtension) Decode(r io.Reader) (image.Image, error) {
	return i.decode(r)
}

type imageExtensionList []imageExtension

func (i imageExtensionList) Includes(extension string) bool {
	_, ok := i.Find(extension)
	return ok
}

func (i imageExtensionList) Find(extension string) (imageExtension, bool) {
	extension = strings.ToLower(extension)
	for _, x := range i {
		for _, y := range x.extensionList {
			if y == extension {
				return x, true
			}
		}
	}
	return imageExtension{}, false
}

var ImageExtensionList imageExtensionList = imageExtensionList{
	{
		name: "bmp",
		extensionList: []string{
			".bmp",
		},
		decode: bmp.Decode,
	},
	{
		name: "gif",
		extensionList: []string{
			".gif",
		},
		decode: gif.Decode,
	},
	{
		name: "jpeg",
		extensionList: []string{
			".jpeg",
			".jpg",
		},
		decode: jpeg.Decode,
	},
	{
		name: "png",
		extensionList: []string{
			".png",
		},
		decode: png.Decode,
	},
}
