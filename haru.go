package harugo

import (
	"log"
	"net/http"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request)

type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	routerGroup := &routerGroup{
		name:          name,
		handleFuncMap: make(map[string]HandleFunc),
	}
	r.routerGroups = append(r.routerGroups, routerGroup)
	return routerGroup
}

func (r *routerGroup) Add(key string, handleFunc HandleFunc) {
	r.handleFuncMap[key] = handleFunc
}

type routerGroup struct {
	name          string
	handleFuncMap map[string]HandleFunc
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{},
	}
}

func (e *Engine) Run() {
	for _, group := range e.routerGroups {
		for key, val := range group.handleFuncMap {
			http.HandleFunc("/"+group.name+key, val)
		}
	}
	if err := http.ListenAndServe(":8111", nil); err != nil {
		log.Fatal(err)
	}
}
