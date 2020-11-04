package server

import (
	"net/http"

	"github.com/Demacr/image_previewer/internal/cacher"
)

type Router struct {
	rootHandler *RootHandler
}

func NewRouter(fc cacher.Cache) *Router {
	return &Router{rootHandler: newRootHandler(fc)}
}

func (router *Router) RootHandler() http.Handler {
	return router.rootHandler
}
