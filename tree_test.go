package harugo

import (
	"fmt"
	"sync"
	"testing"
)

func TestTreeNode(t *testing.T) {
	root := &TreeNode{name: "/", children: make([]*TreeNode, 0)}
	root.Put("/hello/:id")
	root.Put("/user/create/hello")
	root.Put("/user/create/aaa")
	root.Put("/info/:id/aaa")

	ctx := &Context{
		Params: make(map[string]string),
	}

	node := root.Get("/info/hello/*", ctx)
	fmt.Println(node)
	node = root.Get("/user/create/hello", ctx)
	fmt.Println(node)
	node = root.Get("/user/create/aaa", ctx)
	fmt.Println(node)
	node = root.Get("/info/:id/aaa", ctx)
	fmt.Println(node)
}

type MyMutex struct {
	count int
	sync.Mutex
}

func TestXxx(t *testing.T) {
	var mu MyMutex
	mu.Lock()

	mu.count++
	mu.Unlock()
	var mu1 = mu
	mu1.Lock()
	mu1.count++
	mu1.Unlock()
	fmt.Println(mu.count, mu1.count)
}
