package client

import (
	"net/http"
	"strings"
)

func (c *Client) Get(path string) *Request {
	return c.req("GET", path)
}

func (c *Client) Post(path string) *Request {
	return c.req("POST", path)
}

func (c *Client) Put(path string) *Request {
	return c.req("PUT", path)
}

func (c *Client) Delete(path string) *Request {
	return c.req("DELETE", path)
}

func (c *Client) Patch(path string) *Request {
	return c.req("PATCH", path)
}

func (c *Client) Custom(method, path string) *Request {
	return c.req(method, path)
}

func (c *Client) req(method, path string) *Request {
	fullURL := path
	if c.baseURL != "" && !strings.HasPrefix(path, "http") {
		fullURL = c.baseURL + "/" + strings.TrimLeft(path, "/")
	}

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return &Request{client: c, err: err}
	}

	for k, v := range c.headers {
		for _, val := range v {
			req.Header.Add(k, val)
		}
	}

	return &Request{client: c, req: req}
}
