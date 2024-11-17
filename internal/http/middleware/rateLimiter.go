package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type clientData struct {
	timestamps []time.Time
	mu         sync.Mutex
}

var (
	clientRateLimiters = make(map[string]*clientData)
	rateLimiterMutex   sync.Mutex
)

func RateLimitMiddleware(keyFunc func(c *gin.Context) string, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)
		if key == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Rate limiting key missing"})
			return
		}

		rateLimiterMutex.Lock()
		clientLimiter, exists := clientRateLimiters[key]
		if !exists {
			clientLimiter = &clientData{timestamps: []time.Time{}}
			clientRateLimiters[key] = clientLimiter
		}
		rateLimiterMutex.Unlock()

		clientLimiter.mu.Lock()
		defer clientLimiter.mu.Unlock()

		now := time.Now()
		threshold := now.Add(-window)
		validTimestamps := []time.Time{}
		for _, timestamp := range clientLimiter.timestamps {
			if timestamp.After(threshold) {
				validTimestamps = append(validTimestamps, timestamp)
			}
		}
		clientLimiter.timestamps = validTimestamps

		if len(clientLimiter.timestamps) >= maxRequests {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		clientLimiter.timestamps = append(clientLimiter.timestamps, now)

		c.Next()
	}
}
