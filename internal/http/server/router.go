package server

import (
	"net/http"

	domain "github.com/Demacr/image_previewer/internal/domain/previewer"
)

type Router struct {
	rootHandler *RootHandler
}

func NewRouter(p domain.Previewer) *Router {
	return &Router{rootHandler: newRootHandler(p)}
}

func (router *Router) RootHandler() http.Handler {
	return router.rootHandler
}
