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

func (ctx *Context) GetQuery(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) PostBody() map[string]interface{} {
	var params map[string]interface{}
	json.NewDecoder(ctx.Req.Body).Decode(&params)
	return params
}
