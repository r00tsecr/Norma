package tool

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebFetchDefaultUserAgent(t *testing.T) {
	got := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got <- r.UserAgent()
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte("ok"))
	}))
	t.Cleanup(srv.Close)
	raw, _ := json.Marshal(map[string]any{"url": srv.URL})
	res, err := NewWebFetch(WebFetchConfig{}).Call(context.Background(), raw, &ToolContext{})
	if err != nil || res.IsError {
		t.Fatalf("webfetch: err=%v result=%q", err, res.Flatten())
	}
	select {
	case ua := <-got:
		if ua != DefaultWebFetchUserAgent {
			t.Fatalf("User-Agent = %q, want DefaultWebFetchUserAgent", ua)
		}
	default:
		t.Fatal("server received no request")
	}
}

func TestWebFetchCustomUserAgent(t *testing.T) {
	const want = "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.76(0x18004c34) NetType/WIFI Language/zh_CN miniProgram"
	got := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got <- r.UserAgent()
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte("ok"))
	}))
	t.Cleanup(srv.Close)
	raw, _ := json.Marshal(map[string]any{"url": srv.URL})
	res, err := NewWebFetch(WebFetchConfig{
		UserAgent: func(*http.Request) string { return want },
	}).Call(context.Background(), raw, &ToolContext{})
	if err != nil || res.IsError {
		t.Fatalf("webfetch: err=%v result=%q", err, res.Flatten())
	}
	select {
	case ua := <-got:
		if ua != want {
			t.Fatalf("User-Agent = %q, want custom", ua)
		}
	default:
		t.Fatal("server received no request")
	}
}
