package hit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestClient_Do(t *testing.T) {
	const wantHits, wantErrors = 10, 0
	var gotHits atomic.Int64

	handler := func(_ http.ResponseWriter, _ *http.Request) {
		gotHits.Add(1)
	}

	server := newTestServer(t, handler)
	defer server.Close()
	request := newRequest(t, http.MethodGet, server.URL)

	c := &Client{}
	sum := c.Do(context.Background(), request, wantHits)
	if got := gotHits.Load(); got != wantHits {
		t.Errorf("hits=%d; want %d", got, wantHits)
	}
	if got := sum.Requests; got != wantHits {
		t.Errorf("requests=%d; want %d", got, wantHits)
	}
	if got := sum.Errors; got != wantErrors {
		t.Errorf("errors=%d; want %d", got, wantErrors)
	}
}

func newTestServer(tb testing.TB, h http.HandlerFunc) *httptest.Server {
	tb.Helper()
	s := httptest.NewServer(h)
	tb.Cleanup(s.Close)
	return s
}

func newRequest(tb testing.TB, method string, url string) *http.Request {
	tb.Helper()
	request, err := http.NewRequest(method, url, http.NoBody)
	if err != nil {
		tb.Fatalf("NewRequest err=%q; want nil", err)
	}
	return request
}
