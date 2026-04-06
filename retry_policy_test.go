package aisandboxsdk

import (
	"testing"
	"time"
)

func TestShouldRetryTransientHTTP(t *testing.T) {
	if !ShouldRetryTransientHTTP("GET", 503) {
		t.Fatal("GET 503 should retry")
	}
	if ShouldRetryTransientHTTP("POST", 503) {
		t.Fatal("POST 503 must not auto-retry")
	}
	if ShouldRetryTransientHTTP("GET", 404) {
		t.Fatal("GET 404 should not retry")
	}
}

func TestRetryDelay(t *testing.T) {
	if got := RetryDelay(1); got != 200*time.Millisecond {
		t.Fatalf("attempt 1 delay = %v, want 200ms", got)
	}
	if got := RetryDelay(2); got != 400*time.Millisecond {
		t.Fatalf("attempt 2 delay = %v, want 400ms", got)
	}
	if got := RetryDelay(3); got != 800*time.Millisecond {
		t.Fatalf("attempt 3 delay = %v, want 800ms", got)
	}
	if got := RetryDelay(4); got != 1600*time.Millisecond {
		t.Fatalf("attempt 4 delay = %v, want 1600ms", got)
	}
	if got := RetryDelay(5); got != 2*time.Second {
		t.Fatalf("attempt 5 delay = %v, want 2s", got)
	}
	if got := RetryDelay(8); got != 2*time.Second {
		t.Fatalf("attempt 8 delay = %v, want 2s", got)
	}
}
