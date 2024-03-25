package radix

import (
	"testing"
)

func TestServe(t *testing.T) {
	r := Default()
	r.GET("/index", func(ctx *Context) {
		ctx.String("Welcome Radix!")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/hi", func(ctx *Context) {
			arr := []string{"123"}
			ctx.String("%s", arr[99])
		})
	}
	r.REST("/person")
	r.Start()
}
