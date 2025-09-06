package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

func (r *Request) Query(key, value string) *Request {
	if r.err != nil {
		return r
	}
	q := r.req.URL.Query()
	q.Set(key, value)
	r.req.URL.RawQuery = q.Encode()
	return r
}

func (r *Request) Header(key, value string) *Request {
	if r.err != nil {
		return r
	}
	r.req.Header.Set(key, value)
	return r
}

func (r *Request) JSONBody(v any) *Request {
	if r.err != nil {
		return r
	}
	data, err := json.Marshal(v)
	if err != nil {
		r.err = err
		return r
	}
	r.req.Body = io.NopCloser(bytes.NewReader(data))
	r.req.Header.Set("Content-Type", "application/json")
	return r
}

func (r *Request) FormBody(data map[string]string) *Request {
	if r.err != nil {
		return r
	}
	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	r.req.Body = io.NopCloser(strings.NewReader(values.Encode()))
	r.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}
