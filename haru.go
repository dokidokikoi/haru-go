package harugo

import (
	"fmt"
	"log"
	"net/http"
)

const (
	ANY = "ANY"
)

type HandleFunc func(ctx *Context)

type router struct {
	routerGroups map[string]*routerGroup
}

type routerGroup struct {
	name          string
	handleFuncMap map[string]map[string]HandleFunc
	treeNode      *TreeNode
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
		name:          name,
		handleFuncMap: make(map[string]map[string]HandleFunc),
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

func (r *routerGroup) handle(key string, handleFunc HandleFunc, method string) {
	uri := r.name + key
	if _, ok := r.handleFuncMap[key]; !ok {
		r.handleFuncMap[uri] = make(map[string]HandleFunc)
	}
	if _, ok := r.handleFuncMap[key][method]; ok {
		panic("存在重复路由")
	}
	r.treeNode.Put(uri)
	r.handleFuncMap[uri][method] = handleFunc
}

func (r *routerGroup) Any(key string, handleFunc HandleFunc) {
	r.handle(key, handleFunc, ANY)
}

func (r *routerGroup) Get(key string, handleFunc HandleFunc) {
	r.handle(key, handleFunc, http.MethodGet)
}

func (r *routerGroup) Post(key string, handleFunc HandleFunc) {
	r.handle(key, handleFunc, http.MethodPost)
}

func (r *routerGroup) Put(key string, handleFunc HandleFunc) {
	r.handle(key, handleFunc, http.MethodPut)
}

func (r *routerGroup) Delete(key string, handleFunc HandleFunc) {
	r.handle(key, handleFunc, http.MethodDelete)
}

func (r *routerGroup) Patch(key string, handleFunc HandleFunc) {
	r.handle(key, handleFunc, http.MethodPatch)
}

func (r *routerGroup) Options(key string, handleFunc HandleFunc) {
	r.handle(key, handleFunc, http.MethodOptions)
}

func (r *routerGroup) Head(key string, handleFunc HandleFunc) {
	r.handle(key, handleFunc, http.MethodHead)
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

// 路由匹配
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methed := r.Method
	for _, group := range e.routerGroups {
		ctx := &Context{w, r, make(map[string]string)}
		node := group.treeNode.Get(r.RequestURI, ctx)
		if node != nil && node.isEnd {
			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				handle(ctx)
				return
			}
			// 根据 method 匹配
			handle, ok = group.handleFuncMap[node.routerName][methed]
			if !ok {
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintf(w, "%s %s not allowed\n", r.RequestURI, methed)
				return
			} else {
				handle(ctx)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s %s not found\n", r.RequestURI, methed)
}

// 启动服务器
func (e *Engine) Run(addr string) {
	// 所有路由都走 engine.serveHTTP 方法
	http.Handle("/", e)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
