# httpx/cache [wip]

cache control for http servers.

## example

```go
package main

import (
    "net/http"
    cache "github.com/nxtgo/httpx/cache"
)

mux := http.NewServeMux()

mux.Handle("/user", 
    cache.WithETag(
        cache.WithLastModified(
            http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Write([]byte("User profile"))
            }),
            func(r *http.Request) time.Time { return time.Now().Add(-time.Hour) },
        ),
        func(r *http.Request) string { return `"v1-user"` },
    ),
)

http.ListenAndServe(":8080", mux)
```

