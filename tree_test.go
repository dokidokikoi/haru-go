package harugo

import (
	"fmt"
	"testing"
)

func TestTreeNode(t *testing.T) {
	root := &TreeNode{name: "/", children: make([]*TreeNode, 0)}
	root.Put("/user/get/:id")
	root.Put("/user/create/hello")
	root.Put("/user/create/aaa")
	root.Put("/info/:id/aaa")

	ctx := &Context{
		Params: make(map[string]string),
	}

	node := root.Get("/user/get/:id", ctx)
	fmt.Println(node)
	node = root.Get("/user/create/hello", ctx)
	fmt.Println(node)
	node = root.Get("/user/create/aaa", ctx)
	fmt.Println(node)
	node = root.Get("/info/:id/aaa", ctx)
	fmt.Println(node)
}
