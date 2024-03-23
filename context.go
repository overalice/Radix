package radix

import (
	"fmt"
	"net/http"
)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Path   string
	Method string

	statusCode int
}

func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: writer,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (ctx *Context) SetStatusCode(statusCode int) {
	ctx.statusCode = statusCode
	ctx.Writer.WriteHeader(statusCode)
}

func (ctx *Context) SetHeader(key, value string) {
	ctx.Writer.Header().Set(key, value)
}

func (ctx *Context) String(format string, values ...interface{}) {
	if ctx.statusCode == 0 {
		ctx.SetStatusCode(http.StatusOK)
	}
	ctx.SetHeader("Context-Type", "text/plain")
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
