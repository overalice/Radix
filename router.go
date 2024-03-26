package radix

import (
	"net/http"
)

type router struct {
	roots    map[string]*treeNode
	handlers map[string]Handler
}

type treeNode struct {
	path     string
	fullPath string
	isEnd    bool
	indices  string
	children []*treeNode
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*treeNode),
		handlers: make(map[string]Handler),
	}
}

func (router *router) addRouter(method, pattern string, handler Handler) {
	if router.roots[method] == nil {
		router.roots[method] = &treeNode{}
	}
	key := method + "-" + pattern
	if node := router.roots[method].search(pattern); node != nil {
		warn("Failed to add router: router %s already exists", key)
	} else {
		info("Added router: %6s - %s", method, pattern)
		router.roots[method].insert(pattern)
		router.handlers[key] = handler
	}
}

func (router *router) handle(ctx *Context) {
	key := ctx.Method + "-" + ctx.Path
	if router.roots[ctx.Method] == nil {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.String("404 NOT FOUND\t%s", key)
		ctx.Next()
		return
	}

	if node := router.roots[ctx.Method].search(ctx.Path); node != nil {
		ctx.handlers = append(ctx.handlers, router.handlers[ctx.Method+"-"+node.fullPath])
	} else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.SetStatusCode(http.StatusNotFound)
			ctx.String("404 NOT FOUND\t%s", key)
		})
	}
	ctx.Next()
}

func (node *treeNode) insert(pattern string) {
	if node.path == "" && len(node.children) == 0 {
		node.path = pattern
		node.fullPath = pattern
		node.isEnd = true
		return
	}
	fullPattern := pattern
loop:
	for {
		prefix := ""
		for i := 0; i < len(node.path) && i < len(pattern) && node.path[i] == pattern[i]; i++ {
			prefix += string(pattern[i])
		}
		if len(prefix) < len(node.path) {
			child := &treeNode{
				path:     node.path[len(prefix):],
				fullPath: node.fullPath,
				isEnd:    node.isEnd,
				indices:  node.indices,
				children: node.children,
			}

			node.children = []*treeNode{child}
			node.indices = string(node.path[len(prefix)])
			node.isEnd = false
			node.fullPath = node.fullPath[:len(node.fullPath)-(len(node.path)-len(prefix))]
			node.path = prefix
		}
		if len(prefix) < len(pattern) {
			pattern = pattern[len(prefix):]
			for i := 0; i < len(node.indices); i++ {
				if pattern[0] == node.indices[i] {
					node = node.children[i]
					continue loop
				}
			}

			child := &treeNode{path: pattern, fullPath: fullPattern, isEnd: true}
			node.children = append(node.children, child)
			node.indices += string(pattern[0])
			return
		}
		node.isEnd = true
		return
	}
}

func (node *treeNode) search(pattern string) *treeNode {
loop:
	for {
		var i int
		for i = 0; i < len(node.path); i++ {
			if len(pattern) == i {
				return nil
			}
			if node.path[i] == '*' {
				return node
			}
			if node.path[i] != pattern[i] {
				return nil
			}
		}
		pattern = pattern[i:]
		if pattern == "" && node.isEnd {
			return node
		}
		for i = 0; i < len(node.indices); i++ {
			if pattern[0] == node.indices[i] {
				node = node.children[i]
				continue loop
			}
		}
		for i = 0; i < len(node.indices); i++ {
			if node.indices[i] == '*' {
				return node
			}
		}
		return nil
	}
}
