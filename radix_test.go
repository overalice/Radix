package radix

import (
	"testing"
)

func TestServe(t *testing.T) {
	r := Default()
	r.GET("/", func(ctx *Context) {
		ctx.String("Welcome Radix!")
	})
	r.GET("/hi/:name", func(ctx *Context) {
		ctx.String("hi:" + ctx.Param("name"))
	})
	r.GET("/book/:name/get", func(ctx *Context) {
		ctx.String("GET Book:" + ctx.Param("name"))
	})
	r.GET("/book/:name/delete", func(ctx *Context) {
		ctx.String("DELETE Book:" + ctx.Param("name"))
	})
	r.GET("/:name/doing/:thing", func(ctx *Context) {
		ctx.String(ctx.Param("name") + ", you're doing " + ctx.Param("thing"))
	})
	r.GET("/assets/*filename", func(ctx *Context) {
		ctx.String(ctx.Param("filename"))
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/hi", func(ctx *Context) {
			ctx.String("v1 hi ")
		})
		v1.GET("/panic", func(ctx *Context) {
			arr := []string{"123"}
			ctx.String("%s", arr[99])
		})
	}
	r.REST("/person")
	r.Start()
}
