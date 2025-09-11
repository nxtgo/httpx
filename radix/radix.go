package radix

import (
	"strings"
)

// single node in the radix tree
type Node[H any] struct {
	segment  string
	children map[string]*Node[H]
	handler  H
	isParam  bool
	param    string
}

type Router[H any] struct {
	root   *Node[H]
	static map[string]H
}

// creates a new empty router
func NewRouter[H any]() *Router[H] {
	return &Router[H]{
		root:   &Node[H]{children: make(map[string]*Node[H])},
		static: make(map[string]H),
	}
}

// adds a new route
func (r *Router[H]) AddRoute(path string, handler H) {
	if !strings.Contains(path, ":") {
		r.static[path] = handler
		return
	}

	segments := splitPath(path)
	current := r.root

	for _, seg := range segments {
		key := seg
		isParam := false
		param := ""

		if strings.HasPrefix(seg, ":") {
			key = ":"
			isParam = true
			param = seg[1:]
		}

		if current.children == nil {
			current.children = make(map[string]*Node[H])
		}

		child, exists := current.children[key]
		if !exists {
			child = &Node[H]{
				segment:  seg,
				children: make(map[string]*Node[H]),
				isParam:  isParam,
				param:    param,
			}
			current.children[key] = child
		}

		current = child
	}

	current.handler = handler
}

// wiggly wiggly
func (r *Router[H]) Lookup(path string) (handler H, params map[string]string) {
	// check static map first
	if h, ok := r.static[path]; ok {
		return h, nil
	}

	// dynamic routes
	segments := splitPath(path)
	current := r.root
	var p map[string]string

	for _, seg := range segments {
		if child, ok := current.children[seg]; ok {
			current = child
		} else if paramChild, ok := current.children[":"]; ok {
			current = paramChild
			if p == nil {
				p = make(map[string]string)
			}
			p[current.param] = seg
		} else {
			var zero H
			return zero, nil
		}
	}

	return current.handler, p
}

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}
