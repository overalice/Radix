package radix

import (
	"net/http"
	"reflect"
	"strconv"
)

type Handler func(ctx *Context)

type Data map[string]interface{}

func Response(data interface{}) Data {
	return Data{
		"code": http.StatusOK,
		"msg":  "success",
		"data": data,
	}
}

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

func (engine *engine) PUT(pattern string, handler Handler) {
	engine.addRouter("PUT", pattern, handler)
}

func (engine *engine) DELETE(pattern string, handler Handler) {
	engine.addRouter("DELETE", pattern, handler)
}

func (engine *engine) REST(pattern string, keys ...interface{}) {
	var key string
	var identidy string
	if len(keys) == 0 {
		key = "id"
	} else {
		key, ok := keys[0].(string)
		if !ok {
			fault("Expected type of key: string, find %s", reflect.TypeOf(key).String())
			return
		}
	}

	engine.POST(pattern, func(ctx *Context) {
		body := ctx.PostBody()
		switch body[key].(type) {
		case float64:
			identidy = strconv.FormatFloat(body[key].(float64), 'f', -1, 64)
		case string:
			identidy = body[key].(string)
		default:
			response := Response("")
			response["code"] = http.StatusInternalServerError
			response["msg"] = "Expected type of identidy: int or string, find " + reflect.TypeOf(body[key]).String()
			ctx.JSON(response)
			return
		}
		err := SaveJSON(pattern+"-"+identidy+".json", body)
		response := Response(body)
		if err != nil {
			response["code"] = http.StatusInternalServerError
			response["msg"] = "Failed to save data"
		}
		ctx.JSON(response)
	})
	engine.DELETE(pattern, func(ctx *Context) {
		identidy := ctx.GetQuery(key)
		err := RemoveJSON(pattern + "-" + identidy + ".json")
		response := Response("")
		if err != nil {
			response["code"] = http.StatusInternalServerError
			response["msg"] = "Failed to remove data"
		} else {
			response["msg"] = "successfully delete " + key + ": " + identidy
		}
		ctx.JSON(response)
	})
	engine.GET(pattern, func(ctx *Context) {
		data := make(Data)
		identidy = ctx.GetQuery(key)
		err := ReadJSON(pattern+"-"+identidy+".json", data)
		response := Response(data)
		if err != nil {
			response["code"] = http.StatusInternalServerError
			response["msg"] = "Failed to get data"
		}
		ctx.JSON(response)
	})
	engine.PUT(pattern, func(ctx *Context) {
		body := ctx.PostBody()
		switch body[key].(type) {
		case float64:
			identidy = strconv.FormatFloat(body[key].(float64), 'f', -1, 64)
		case string:
			identidy = body[key].(string)
		default:
			response := Response("")
			response["code"] = http.StatusInternalServerError
			response["msg"] = "Expected type of key: int or string, find " + reflect.TypeOf(body[key]).String()
			ctx.JSON(response)
			return
		}
		data := make(Data)
		err := ReadJSON(pattern+"-"+identidy+".json", data)

		for key, value := range body {
			data[key] = value
		}
		response := Response(data)
		if err != nil {
			response["code"] = http.StatusInternalServerError
			response["msg"] = "Failed to get data"
		} else {
			err = SaveJSON(pattern+"-"+identidy+".json", data)
			if err != nil {
				response["code"] = http.StatusInternalServerError
				response["msg"] = "Failed to save data"
			}
		}
		ctx.JSON(response)
	})
}

func (engine *engine) Start(addrs ...interface{}) {
	var addr string
	if len(addrs) == 0 {
		addr = "8080"
	} else {
		var ok bool
		addr, ok = addrs[0].(string)
		if !ok {
			fault("Expected type of addr: string, find %s", reflect.TypeOf(addrs[0]).String())
			return
		}
	}
	info("Welcome Radix, service running on 127.0.0.1:%s", addr)
	http.ListenAndServe(":"+addr, engine)
}

func (engine *engine) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := newContext(writer, req)
	engine.router.handle(ctx)
}
