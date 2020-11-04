package server

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Demacr/image_previewer/internal/cacher"
)

type FillHandler struct {
	fc cacher.Cache
}

func newFillHandler(fc cacher.Cache) *FillHandler {
	return &FillHandler{fc: fc}
}

func (h *FillHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var widthStr, heightStr string
	widthStr, r.URL.Path = shiftPath(r.URL.Path)
	heightStr, r.URL.Path = shiftPath(r.URL.Path)

	width, err := strconv.Atoi(widthStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageURL := strings.TrimPrefix(r.URL.Path, "/")
	image, err := h.fc.GetImage(imageURL, width, height)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getimage:", err)
		return
	}

	for key, value := range image.Header {
		if key != "Content-Length" {
			for _, v := range value {
				w.Header().Add(key, v)
			}
		}
	}
	n, err := io.Copy(w, image.Buffer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error during copying image:", err)
		return
	}
	defer image.Buffer.Close()

	log.Printf("finished request, wrote %v bytes\n", n)
}
