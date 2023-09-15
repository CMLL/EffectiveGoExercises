package hit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"
)

type Client struct {
	C       int
	RPS     int
	Timeout time.Duration
}

type Option func(*Client)

func Concurrency(n int) Option {
	return func(c *Client) {
		c.C = n
	}
}

func Timeout(d time.Duration) Option {
	return func(c *Client) {
		c.Timeout = d
	}
}

type SendFunc func(r *http.Request) *Result

func (c *Client) Do(ctx context.Context, r *http.Request, n int) *Result {
	t := time.Now()
	sum := c.do(ctx, r, n)
	return sum.Finalize(time.Since(t))
}

func (c *Client) do(ctx context.Context, r *http.Request, n int) *Result {
	p := produce(ctx, n, func() *http.Request {
		return r.Clone(ctx)
	})
	if c.RPS > 0 {
		p = throttle(p, time.Second/time.Duration(c.RPS*c.concurrency()))
	}
	var (
		sum    Result
		client = c.client()
	)
	//defer client.CloseIdleConnections()
	for result := range split(p, c.concurrency(), c.send(client)) {
		sum.Merge(result)
	}
	return &sum
}

func (c *Client) send(client *http.Client) SendFunc {
	return func(r *http.Request) *Result {
		return Send(client, r)
	}

}

func (c *Client) client() *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: c.concurrency(),
		},
	}
}

func (c *Client) concurrency() int {
	if c.C > 0 {
		return c.C
	}
	return runtime.NumCPU()
}

func Send(client *http.Client, r *http.Request) *Result {
	t := time.Now()

	var (
		code  int
		bytes int64
	)
	response, err := client.Do(r)
	if err == nil {
		code = response.StatusCode
		bytes, err = io.Copy(io.Discard, response.Body)
		response.Body.Close()
	}

	return &Result{
		Duration: time.Since(t),
		Bytes:    bytes,
		Status:   code,
		Error:    err,
	}
}

func Do(ctx context.Context, url string, n int, opts ...Option) (*Result, error) {
	r, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("new http request: %w", err)
	}
	var c Client
	for _, option := range opts {
		option(&c)
	}
	return c.Do(ctx, r, n), nil

}
