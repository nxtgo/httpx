package main

import (
	"fmt"
	"net/http"

	"github.com/nxtgo/httpx/router"
)

func main() {
	// create a new router
	r := router.New()

	// register a global middleware
	r.Use(func(next router.Handler) router.Handler {
		return func(w http.ResponseWriter, req *http.Request, p router.Params) {
			fmt.Printf("request: %s %s\n", req.Method, req.URL.Path)
			next(w, req, p)
		}
	})

	// register a get route with a parameter
	r.Get("/users/:id", func(w http.ResponseWriter, req *http.Request, p router.Params) {
		fmt.Fprintf(w, "user id: %s\n", p["id"])
	})

	// register a static get route
	r.Get("/about", func(w http.ResponseWriter, req *http.Request, p router.Params) {
		fmt.Fprintln(w, "about page")
	})

	// start the http server
	fmt.Println("server running on :8080")
	http.ListenAndServe(":8080", r)
}
