package aisandboxsdk

import (
	"strings"
	"time"
)

const (
	retryBaseDelay = 200 * time.Millisecond
	retryMaxDelay  = 2 * time.Second
)

// ShouldRetryTransientHTTP retries 429/5xx only for GET/HEAD to avoid duplicate write side effects.
func ShouldRetryTransientHTTP(method string, status int) bool {
	if status != 429 && (status < 500 || status > 599) {
		return false
	}
	m := strings.ToUpper(strings.TrimSpace(method))
	return m == "GET" || m == "HEAD"
}

func RetryDelay(attempt int) time.Duration {
	safeAttempt := attempt
	if safeAttempt < 1 {
		safeAttempt = 1
	}
	delay := retryBaseDelay * time.Duration(1<<(safeAttempt-1))
	if delay > retryMaxDelay {
		return retryMaxDelay
	}
	return delay
}
