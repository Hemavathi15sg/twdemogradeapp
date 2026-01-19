package middleware

import (
	"net/http"
	"time"
)

// CacheStatus tracks cache hits and misses for debugging
type CacheStatus string

const (
	CacheHit    CacheStatus = "HIT"
	CacheMiss   CacheStatus = "MISS"
	CacheUpdate CacheStatus = "UPDATE"
	CacheDelete CacheStatus = "DELETE"
)

// SetCacheHeader adds debugging cache header to response
func SetCacheHeader(w http.ResponseWriter, status CacheStatus, ttl ...time.Duration) {
	headerValue := string(status)
	if len(ttl) > 0 && ttl[0] > 0 {
		headerValue += "; TTL=" + ttl[0].String()
	}
	w.Header().Set("X-Cache", headerValue)
}
