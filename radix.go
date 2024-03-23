package radix

import "net/http"

type Handler func(c *Context)

type engine struct {
	router *router
}

func New() *engine {
	return &engine{
		router: newRouter(),
	}
}

func (engine *engine) addRouter(method, pattern string, handler Handler) {
	engine.router.addRouter(method, pattern, handler)
}

func (engine *engine) GET(pattern string, handler Handler) {
	engine.addRouter("GET", pattern, handler)
}

func (engine *engine) POST(pattern string, handler Handler) {
	engine.addRouter("POST", pattern, handler)
}

func (engine *engine) Start() {
	info("Welcome Radix, service running on 127.0.0.1:8080")
	http.ListenAndServe(":8080", engine)
}

func (engine *engine) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := newContext(writer, req)
	engine.router.handle(ctx)
}
