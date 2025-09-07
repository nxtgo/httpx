package cache

import (
	"fmt"
	"net/http"
	"time"
)

func WithETag(next http.Handler, gen func(*http.Request) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		etag := gen(r)

		if match := r.Header.Get("If-None-Match"); match == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.Header().Set("ETag", etag)
		next.ServeHTTP(w, r)
	})
}

func WithLastModified(next http.Handler, modTime func(*http.Request) time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastMod := modTime(r)

		if since := r.Header.Get("If-Modified-Since"); since != "" {
			if t, _ := time.Parse(http.TimeFormat, since); !lastMod.After(t) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		w.Header().Set("Last-Modified", lastMod.UTC().Format(http.TimeFormat))
		next.ServeHTTP(w, r)
	})
}

func NoCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
}

func Cached(w http.ResponseWriter, d time.Duration) {
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(d.Seconds())))
}
