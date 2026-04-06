package aisandboxsdk

import "testing"

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
