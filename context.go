package harugo

import "net/http"

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Params map[string]string
}
