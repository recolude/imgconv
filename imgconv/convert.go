package imgconv

import (
	"image"
	"image/png"
	"io"

	"github.com/ftrvxmtrx/tga"
	"github.com/nfnt/resize"
	"golang.org/x/image/tiff"
)

func resizeIfBigger(img image.Image, resizeVal int) image.Image {
	if resizeVal > 0 && (img.Bounds().Dx() > resizeVal || img.Bounds().Dy() > resizeVal) {
		return resize.Resize(uint(resizeVal), uint(resizeVal), img, resize.NearestNeighbor)
	}
	return img
}

func ConvertTGA(tgaIn io.Reader, pngOut io.Writer, resizeVal int) error {
	img, err := tga.Decode(tgaIn)
	if err != nil {
		return err
	}
	img = resizeIfBigger(img, resizeVal)

	return png.Encode(pngOut, img)
}

func ConvertTIFF(tiffIn io.Reader, pngOut io.Writer, resizeVal int) error {
	img, err := tiff.Decode(tiffIn)
	if err != nil {
		return err
	}
	img = resizeIfBigger(img, resizeVal)

	return png.Encode(pngOut, img)
}

func ResizePNG(pngIn io.Reader, pngOut io.Writer, resizeVal int) error {
	img, err := png.Decode(pngIn)
	if err != nil {
		return err
	}
	img = resizeIfBigger(img, resizeVal)

	return png.Encode(pngOut, img)
}
