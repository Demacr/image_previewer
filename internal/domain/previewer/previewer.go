package previewer

import (
	"bytes"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math"

	"github.com/Demacr/image_previewer/internal/cacher"
	domain "github.com/Demacr/image_previewer/internal/domain/types"
	"github.com/Demacr/image_previewer/internal/http/client"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
)

type Previewer interface {
	GetImage(url string, width int, height int) (*domain.DownloadedImage, error)
}

type previewer struct {
	fc cacher.Cache
}

func NewPreviewer(fc cacher.Cache) Previewer {
	return &previewer{
		fc: fc,
	}
}

func (p *previewer) GetImage(url string, width int, height int) (*domain.DownloadedImage, error) {
	nameBytes := sha256.Sum256([]byte(fmt.Sprintf("%v%v%v", url, width, height)))
	key := base32.StdEncoding.EncodeToString(nameBytes[:])
	result, err := p.fc.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "error during getting image from disk cache")
	}
	if result != nil {
		log.Println("get image from cache")
		return result, nil
	}

	result, err = client.GetImage(url)
	if err != nil {
		return nil, err
	}

	result.Buffer, err = Preview(result.Buffer, width, height)
	if err != nil {
		return nil, errors.Wrap(err, "error during image previewing")
	}

	if err = p.fc.Set(key, result); err != nil {
		return nil, errors.Wrap(err, "error during saving previewed image")
	}
	result, err = p.fc.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "error during getting image from disk cache (new)")
	}
	return result, nil
}

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
