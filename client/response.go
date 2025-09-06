package client

import (
	"encoding/json"
	"io"
	"net/http"
)

func (r *Request) Do() (*http.Response, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.client.client.Do(r.req)
}

func (r *Request) Bytes() ([]byte, error) {
	resp, err := r.Do()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (r *Request) String() (string, error) {
	b, err := r.Bytes()
	return string(b), err
}

func (r *Request) JSON(v any) error {
	resp, err := r.Do()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func (r *Request) Decode(f func(resp *http.Response) (any, error)) (any, error) {
	resp, err := r.Do()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return f(resp)
}
