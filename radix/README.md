# httpx/radix

*completely* agnostic radix tree structure implementation.

## example

```go
router := radix.NewRouter[func()]()

router.AddRoute("/users/:id", func() { fmt.Println("handler for /users/:id") })
router.AddRoute("/posts/:postId/comments/:commentId", func() { fmt.Println("handler for /posts/:postId/comments/:commentId") })
router.AddRoute("/about", func() { fmt.Println("handler for /about") })

// ...
```
