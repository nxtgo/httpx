package router

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/nxtgo/httpx/radix"
)

// Handler is the type for http Handlers used in this Router.
type Handler func(w http.ResponseWriter, r *http.Request, p Params)

// Params stores path parameters extracted from the route.
type Params map[string]string

// Middleware type.
type Middleware func(Handler) Handler

// String returns the string value of a parameter by key.
func (p Params) String(key string) string {
	return p[key]
}

// Int converts a parameter to int by key, returns error if conversion fails.
func (p Params) Int(key string) (int, error) {
	return strconv.Atoi(p[key])
}

// Router is the main HTTP router structure.
type Router struct {
	trees       map[string]*radix.Router[Handler]
	static      map[string]map[string]Handler
	middlewares []Middleware
}

// New creates a new empty Router with initialized method trees.
func New() *Router {
	return &Router{
		trees:  make(map[string]*radix.Router[Handler]),
		static: make(map[string]map[string]Handler),
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

	if !strings.Contains(path, ":") {
		if _, ok := r.static[method]; !ok {
			r.static[method] = make(map[string]Handler)
		}
		r.static[method][path] = h
		return
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

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path

	// check static routes first
	if m, ok := r.static[method]; ok {
		if h, ok := m[path]; ok {
			h(w, req, nil)
			return
		}
	}

	tree, ok := r.trees[method]
	if !ok {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h, p := tree.Lookup(path)
	if h == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	h(w, req, p)
}
