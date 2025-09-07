package router

import (
	"net/http"
	"strconv"

	"github.com/nxtgo/httpx/radix"
)

// Handler is the type for http Handlers used in this Router.
// it receives the response writer, the request, and the route Params.
type Handler func(w http.ResponseWriter, r *http.Request, p Params)

// Params stores path parameters extracted from the route.
type Params map[string]string

// Middleware type.
type Middleware func(Handler) Handler

// string returns the string value of a parameter by key.
func (p Params) String(key string) string {
	return p[key]
}

// int converts a parameter to int by key, returns error if conversion fails.
func (p Params) Int(key string) (int, error) {
	return strconv.Atoi(p[key])
}

// Router is the main http router structure.
// it stores a radix tree for each http method, and optional notfound and methodnotallowed Handlers.
type Router struct {
	trees       map[string]*radix.Router[Handler]
	middlewares []Middleware
}

// New creates a new empty Router with initialized method trees.
func New() *Router {
	return &Router{
		trees: map[string]*radix.Router[Handler]{},
	}
}

// Use registers a new middleware. (global)
func (r *Router) Use(middle Middleware) {
	r.middlewares = append(r.middlewares, middle)
}

// Handle registers a Handler for a specific method and path.
func (r *Router) Handle(method, path string, h Handler) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	tree, ok := r.trees[method]
	if !ok {
		tree = radix.NewRouter[Handler]()
		r.trees[method] = tree
	}
	tree.AddRoute(path, h)
}

func (r *Router) Get(path string, h Handler) {
	r.Handle(http.MethodGet, path, h)
}

func (r *Router) Post(path string, h Handler) {
	r.Handle(http.MethodPost, path, h)
}

func (r *Router) Put(path string, h Handler) {
	r.Handle(http.MethodPut, path, h)
}

func (r *Router) Delete(path string, h Handler) {
	r.Handle(http.MethodDelete, path, h)
}

// servehttp makes the Router implement http.Handler.
// it looks up the request path in the correct method tree and calls the Handler if found.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	tree, ok := r.trees[req.Method]
	if !ok {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hptr, p := tree.Lookup(req.URL.Path)
	if hptr == nil || *hptr == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	(*hptr)(w, req, p)
}
