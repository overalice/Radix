package radix

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Path   string
	Method string

	statusCode int

	params map[string]string

	index    int
	handlers []Handler

	engine *engine
}

func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: writer,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		params: make(map[string]string),
		index:  -1,
	}
}

func (ctx *Context) Param(key string) string {
	value := ctx.params[key]
	return value
}

func (ctx *Context) Next() {
	ctx.index++
	size := len(ctx.handlers)
	for ; ctx.index < size; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
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
	ctx.SetHeader("Content-Type", "text/plain")
	if ctx.statusCode == 0 {
		ctx.SetStatusCode(http.StatusOK)
	}
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (ctx *Context) JSON(data Data) {
	ctx.SetHeader("Content-Type", "application/json")
	if statusCode, exist := data["code"]; exist {
		ctx.SetStatusCode(statusCode.(int))
	}
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(data); err != nil {
		ctx.SetStatusCode(500)
		ctx.String(err.Error())
	}
}

func (ctx *Context) HTML(name string, data interface{}) {
	ctx.SetHeader("Context-Type", "text/html")
	if ctx.statusCode == 0 {
		ctx.SetStatusCode(http.StatusOK)
	}
	if err := ctx.engine.htmlTemplates.ExecuteTemplate(ctx.Writer, name, data); err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.String(err.Error())
	}
}

func (ctx *Context) GetQuery(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) PostBody() map[string]interface{} {
	var params map[string]interface{}
	json.NewDecoder(ctx.Req.Body).Decode(&params)
	return params
}
