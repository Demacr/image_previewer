package server

import (
	"net/http"

	domain "github.com/Demacr/image_previewer/internal/domain/previewer"
)

type RootHandler struct {
	fillHandler *FillHandler
}

func newRootHandler(p domain.Previewer) *RootHandler {
	return &RootHandler{
		fillHandler: newFillHandler(p),
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
