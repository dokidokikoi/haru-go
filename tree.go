// 前缀树
package harugo

import (
	"strings"
)

type TreeNode struct {
	name       string
	children   []*TreeNode
	routerName string
	isEnd      bool
	isWild     bool
}

func (t *TreeNode) Put(path string) {
	root := t
	strs := strings.Split(path, "/")
next:
	for i, name := range strs {
		if i == 0 {
			continue
		}

		children := t.children
		for _, node := range children {
			if node.name == name {
				t = node
				continue next
			}
		}

		node := &TreeNode{
			name:       name,
			children:   make([]*TreeNode, 0),
			routerName: t.routerName + "/" + name,
			isEnd:      true,
			isWild:     strings.Contains(name, ":"),
		}
		t.children = append(children, node)
		t.isEnd = false
		t = node
	}
	t = root
}

func (t *TreeNode) Get(path string, ctx *Context) *TreeNode {
	strs := strings.Split(path, "/")
next:
	for i, name := range strs {
		if i == 0 {
			continue
		}

		children := t.children
		for _, node := range children {
			if node.name == name || node.name == "*" || strings.Contains(node.name, ":") {
				t = node
				if node.isWild {
					ctx.Params[node.name[1:]] = name
				}
				if i == len(strs)-1 {
					return node
				}
				continue next
			}
		}

		for _, node := range children {
			if node.name == "**" {
				return node
			}
		}
	}
	return nil
}
