package radix

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"text/template"
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
	*RouterGroup
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

type RouterGroup struct {
	prefix      string
	middlewares []Handler
	engine      *engine
}

func New() *engine {
	engine := &engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *engine {
	engine := New()
	engine.Use(Reconvery)
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (engine *engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (group *RouterGroup) Use(middlewares ...Handler) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRouter(method, pattern string, handler Handler) {
	group.engine.router.addRouter(method, group.prefix+pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler Handler) {
	group.addRouter("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler Handler) {
	group.addRouter("POST", pattern, handler)
}

func (group *RouterGroup) PUT(pattern string, handler Handler) {
	group.addRouter("PUT", pattern, handler)
}

func (group *RouterGroup) DELETE(pattern string, handler Handler) {
	group.addRouter("DELETE", pattern, handler)
}

func (group *RouterGroup) REST(pattern string, keys ...interface{}) {
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

	group.POST(pattern, func(ctx *Context) {
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
	group.DELETE(pattern, func(ctx *Context) {
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
	group.GET(pattern, func(ctx *Context) {
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
	group.PUT(pattern, func(ctx *Context) {
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
		addr = config["port"]
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
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			ctx.handlers = append(ctx.handlers, group.middlewares...)
		}
	}
	ctx.engine = engine
	engine.router.handle(ctx)
}
