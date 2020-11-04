package server

import (
	"net/http"

	"github.com/Demacr/image_previewer/internal/cacher"
)

type RootHandler struct {
	fillHandler *FillHandler
}

func newRootHandler(fc cacher.Cache) *RootHandler {
	return &RootHandler{
		fillHandler: newFillHandler(fc),
	}
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	switch head {
	case "fill":
		h.fillHandler.ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}
