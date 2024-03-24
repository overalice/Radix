package radix

import (
	"testing"
)

func TestServe(t *testing.T) {
	r := New()
	r.GET("/index", func(ctx *Context) {
		ctx.String("Welcome Radix!")
	})
	r.REST("/person")
	r.Start()
}
