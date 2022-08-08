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

	node := root.Get("/user/get/:id")
	fmt.Println(node)
	node = root.Get("/user/create/hello")
	fmt.Println(node)
	node = root.Get("/user/create/aaa")
	fmt.Println(node)
	node = root.Get("/info/:id/aaa")
	fmt.Println(node)
}
