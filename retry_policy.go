package aisandboxsdk

import "strings"

// ShouldRetryTransientHTTP 僅對 GET/HEAD 在 429 或 5xx 時建議重試，避免寫入重複。
func ShouldRetryTransientHTTP(method string, status int) bool {
	if status != 429 && (status < 500 || status > 599) {
		return false
	}
	m := strings.ToUpper(strings.TrimSpace(method))
	return m == "GET" || m == "HEAD"
}
