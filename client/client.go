package client

import (
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	headers http.Header
	client  *http.Client
}

type Request struct {
	client *Client
	req    *http.Request
	err    error
}

func New() *Client {
	return &Client{
		client:  http.DefaultClient,
		headers: make(http.Header),
	}
}

func (c *Client) BaseURL(url string) *Client {
	c.baseURL = url
	return c
}

func (c *Client) Header(k, v string) *Client {
	c.headers.Add(k, v)
	return c
}

func (c *Client) Headers(headers map[string]string) *Client {
	for key, value := range headers {
		c.headers.Add(key, value)
	}
	return c
}

func (c *Client) Timeout(d time.Duration) *Client {
	c.client.Timeout = d
	return c
}

func (c *Client) HTTPClient(hc *http.Client) *Client {
	c.client = hc
	return c
}
