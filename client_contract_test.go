package aisandboxsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

const (
	testAPIKey       = "test-key"
	tracePostHeader  = "trace-post-1"
	traceUpperHeader = "trace-upper-case"
)

func TestClientContractGetRetriesOnTransient5xx(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/agents" {
			http.Error(w, "unexpected route", http.StatusBadRequest)
			return
		}

		next := atomic.AddInt32(&calls, 1)
		if next == 1 {
			http.Error(w, "temporary failure", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"payload":{"items":[]}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, testAPIKey)
	client.MaxRetries = 2

	result, err := client.ListAgents("")
	if err != nil {
		t.Fatalf("ListAgents returned error: %v", err)
	}

	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("GET calls = %d, want 2", got)
	}
	if ok, _ := result["ok"].(bool); !ok {
		t.Fatalf("expected ok=true, got: %#v", result["ok"])
	}
}

func TestClientContractPostDoesNotRetryOnHttp5xx(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/agents/123/chat" {
			http.Error(w, "unexpected route", http.StatusBadRequest)
			return
		}
		atomic.AddInt32(&calls, 1)
		w.Header().Set("x-trace-id", tracePostHeader)
		http.Error(w, "write failed", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, testAPIKey)
	client.MaxRetries = 2

	_, err := client.SendChat(123, map[string]any{"channelId": "c-1", "message": "hello"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected ApiError, got %T", err)
	}
	if apiErr.Status != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", apiErr.Status, http.StatusServiceUnavailable)
	}
	if apiErr.TraceID != tracePostHeader {
		t.Fatalf("trace id = %q, want %q", apiErr.TraceID, tracePostHeader)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("POST calls = %d, want 1", got)
	}
}

type flakyRoundTripper struct {
	calls int
}

func (f *flakyRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	f.calls++
	if f.calls == 1 {
		return nil, io.ErrUnexpectedEOF
	}

	body := io.NopCloser(strings.NewReader(`{"ok":true,"payload":{"messageId":"m-1"}}`))
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     headers,
		Body:       body,
	}, nil
}

func TestClientContractPostRetriesOnNetworkFault(t *testing.T) {
	rt := &flakyRoundTripper{}
	client := NewClient("https://api.example.test", testAPIKey)
	client.MaxRetries = 2
	client.HTTPClient = &http.Client{Transport: rt}

	result, err := client.SendChat(123, map[string]any{"channelId": "c-1", "message": "hello"})
	if err != nil {
		t.Fatalf("SendChat returned error: %v", err)
	}
	if rt.calls != 2 {
		t.Fatalf("network calls = %d, want 2", rt.calls)
	}

	raw, _ := json.Marshal(result)
	if !strings.Contains(string(raw), "\"ok\":true") {
		t.Fatalf("unexpected payload: %s", string(raw))
	}
}

func TestClientContractGetTraceIDHeaderCaseInsensitive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Trace-Id", traceUpperHeader)
		http.Error(w, "failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, testAPIKey)
	client.MaxRetries = 0

	_, err := client.GetAgentDetail(1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected ApiError, got %T", err)
	}
	if apiErr.TraceID != traceUpperHeader {
		t.Fatalf("trace id = %q, want %q", apiErr.TraceID, traceUpperHeader)
	}
	if !strings.Contains(apiErr.Error(), fmt.Sprintf("trace_id=%s", apiErr.TraceID)) {
		t.Fatalf("error string should include trace_id, got: %s", apiErr.Error())
	}
}
