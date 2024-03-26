package radix

import "net/http"

func Reconvery(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			ctx.String("Internal Server Error: %s", err)
		}
	}()
	ctx.Next()
}
