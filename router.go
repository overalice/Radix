package radix

import (
	"net/http"
)

type router struct {
	handlers map[string]Handler
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]Handler),
	}
}

func (router *router) addRouter(method, pattern string, handler Handler) {
	key := method + "-" + pattern
	if _, exist := router.handlers[key]; exist {
		warn("Failed to add router: router %s already exists", key)
	} else {
		info("Added router: %s", key)
		router.handlers[key] = handler
	}
}

func (router *router) handle(ctx *Context) {
	key := ctx.Method + "-" + ctx.Path
	if handler, exist := router.handlers[key]; exist {
		handler(ctx)
	} else {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.String("404 NOT FOUND\t%s", key)
	}
}
