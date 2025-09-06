# httpx/client

simple abstraction to manage http requests like a chad.

## examples

#### simple

```go
// we will rename the package to "httpx" since
// it's the only thing we're using here.
import httpx "github.com/nxtgo/httpx/client"

client := httpx.New().
    BaseURL("https://httpbin.org").
    Header("User-Agent", "httpx-demo")
```

#### get

```go
resp, err := client.Get("/get").
    Query("hello", "world").
    Do()
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

body, _ := io.ReadAll(resp.Body)
fmt.Println(string(body))
```

#### get json into struct

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var user User
err := client.Get("/json").
    JSON(&user)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("user: %+v\n", user)
```
