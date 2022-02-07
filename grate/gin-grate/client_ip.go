package gin_grate

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/gladmo/toolset/grate"
)

// IPLimiter gin ip rate limiter
func IPLimiter(limit *grate.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		if !limit.Allow(c.ClientIP()) {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		
		c.Next()
	}
}
