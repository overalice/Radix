package radix

import (
	"net/http"
	"testing"
)

func TestServe(t *testing.T) {
	r := New()
	r.GET("/index", func(ctx *Context) {
		ctx.SetStatusCode(http.StatusOK)
		ctx.String("Welcome Radix!")
	})
	r.GET("/index", func(ctx *Context) {
		ctx.SetStatusCode(http.StatusOK)
		ctx.String("Welcome xxxx!")
	})
	r.Start()
}
