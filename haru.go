package harugo

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	ANY = "ANY"
)

type MiddlewareFunc func(handleFunc HandleFunc) HandleFunc

type HandleFunc func(ctx *Context)

type router struct {
	routerGroups map[string]*routerGroup
}

type routerGroup struct {
	name              string
	handleFuncMap     map[string]map[string]HandleFunc
	middlewareFuncMap map[string]map[string][]MiddlewareFunc
	treeNode          *TreeNode
	middlewares       []MiddlewareFunc
}

func (r *routerGroup) Use(middlewareFunc ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewareFunc...)
}

func (r *routerGroup) Methodhandle(h HandleFunc, ctx *Context, routerMiddlewares []MiddlewareFunc) {
	// 叠加中间件
	// 路由中间件
	for _, middlewareFunc := range routerMiddlewares {
		h = middlewareFunc(h)
	}

	// 全局中间件
	for _, middlewareFunc := range r.middlewares {
		h = middlewareFunc(h)
	}

	h(ctx)
}

// 路由分组
func (r *router) Group(name string) *routerGroup {
	if name == "/" {
		name = ""
	}
	if group, ok := r.routerGroups[name]; ok {
		return group
	}

	routerGroup := &routerGroup{
		name:              name,
		handleFuncMap:     make(map[string]map[string]HandleFunc),
		middlewareFuncMap: make(map[string]map[string][]MiddlewareFunc),
		treeNode: &TreeNode{
			name:       "",
			children:   make([]*TreeNode, 0),
			routerName: "",
			isEnd:      true,
		},
	}
	r.routerGroups[name] = routerGroup
	return routerGroup
}

// func (r *routerGroup) Add(key string, handleFunc HandleFunc) {
// 	r.handleFuncMap[key] = handleFunc
// }

func (r *routerGroup) handle(key string,
	handleFunc HandleFunc,
	method string,
	middlewares ...MiddlewareFunc) {

	if _, ok := r.handleFuncMap[key]; !ok {
		r.handleFuncMap[key] = make(map[string]HandleFunc)
	}

	if _, ok := r.middlewareFuncMap[key]; !ok {
		r.middlewareFuncMap[key] = make(map[string][]MiddlewareFunc)
	}

	if _, ok := r.handleFuncMap[key][method]; ok {
		panic("存在重复路由")
	}
	r.treeNode.Put(key)
	r.handleFuncMap[key][method] = handleFunc
	r.middlewareFuncMap[key][method] = append(r.middlewareFuncMap[key][method], middlewares...)
}

func (r *routerGroup) Any(key string,
	handleFunc HandleFunc,
	middlewares ...MiddlewareFunc) {
	r.handle(key, handleFunc, ANY, middlewares...)
}

func (r *routerGroup) Get(key string,
	handleFunc HandleFunc,
	middlewares ...MiddlewareFunc) {
	r.handle(key, handleFunc, http.MethodGet, middlewares...)
}

func (r *routerGroup) Post(key string,
	handleFunc HandleFunc,
	middlewares ...MiddlewareFunc) {
	r.handle(key, handleFunc, http.MethodPost, middlewares...)
}

func (r *routerGroup) Put(key string,
	handleFunc HandleFunc,
	middlewares ...MiddlewareFunc) {
	r.handle(key, handleFunc, http.MethodPut, middlewares...)
}

func (r *routerGroup) Delete(key string,
	handleFunc HandleFunc,
	middlewares ...MiddlewareFunc) {
	r.handle(key, handleFunc, http.MethodDelete, middlewares...)
}

func (r *routerGroup) Patch(key string,
	handleFunc HandleFunc,
	middlewares ...MiddlewareFunc) {
	r.handle(key, handleFunc, http.MethodPatch, middlewares...)
}

func (r *routerGroup) Options(key string,
	handleFunc HandleFunc,
	middlewares ...MiddlewareFunc) {
	r.handle(key, handleFunc, http.MethodOptions, middlewares...)
}

func (r *routerGroup) Head(key string,
	handleFunc HandleFunc,
	middlewares ...MiddlewareFunc) {
	r.handle(key, handleFunc, http.MethodHead, middlewares...)
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{
			routerGroups: make(map[string]*routerGroup),
		},
	}
}

func SubstringLast(str string, substr string) string {
	index := strings.Index(str, substr)
	if index < 0 {
		return ""
	}

	return str[index+len(substr):]
}

func (e *Engine) HttpRequestHandle(w http.ResponseWriter, r *http.Request) {
	methed := r.Method
	for _, group := range e.routerGroups {
		ctx := &Context{w, r, make(map[string]string)}
		routerName := SubstringLast(r.RequestURI, group.name)
		node := group.treeNode.Get(routerName, ctx)
		if node != nil && node.isEnd {
			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				group.Methodhandle(handle, ctx, group.middlewareFuncMap[node.routerName][ANY])
				return
			}
			// 根据 method 匹配
			handle, ok = group.handleFuncMap[node.routerName][methed]
			if !ok {
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintf(w, "%s %s not allowed\n", r.RequestURI, methed)
				return
			} else {
				group.Methodhandle(handle, ctx, group.middlewareFuncMap[node.routerName][methed])
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s %s not found\n", r.RequestURI, methed)
}

// 路由匹配
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.HttpRequestHandle(w, r)
}

// 启动服务器
func (e *Engine) Run(addr string) {
	// 所有路由都走 engine.serveHTTP 方法
	http.Handle("/", e)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
