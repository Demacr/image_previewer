package previewer

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math"

	"github.com/pkg/errors"
	"golang.org/x/image/draw"
)

func Preview(imageOrig io.ReadCloser, width int, heigth int) (io.ReadCloser, error) {
	imageDecoded, err := jpeg.Decode(imageOrig)
	if err != nil {
		return nil, errors.Wrap(err, "bad jpeg file")
	}
	defer imageOrig.Close()

	result := image.NewRGBA(image.Rect(0, 0, width, heigth))

	srcX := imageDecoded.Bounds().Dx()
	srcY := imageDecoded.Bounds().Dy()
	var r1 float64 = float64(srcX) / float64(srcY)
	var r2 float64 = float64(width) / float64(heigth)
	var scaledRect image.Rectangle

	switch {
	case math.Abs(r2-r1) < 0.0001:
		scaledRect = imageDecoded.Bounds()
	case r2 > r1:
		newY := int(float64(srcY) * r1 / r2)
		delta := (srcY - newY) / 2
		scaledRect = image.Rect(0, delta, srcX, srcY-delta)
	case r2 < r1:
		newX := int(float64(srcX) * r2 / r1)
		delta := (srcX - newX) / 2
		scaledRect = image.Rect(delta, 0, srcX-delta, srcY)
	}
	draw.CatmullRom.Scale(result, result.Bounds(), imageDecoded, scaledRect, draw.Src, nil)

	buf := bytes.Buffer{}

	if err = jpeg.Encode(&buf, result, nil); err != nil {
		return nil, errors.Wrap(err, "error during encoding")
	}
	resultBuf := ioutil.NopCloser(&buf)
	return resultBuf, nil
}
