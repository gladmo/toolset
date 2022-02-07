package gin_grate

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/gladmo/toolset/grate"
)

// RouteLimiter gin request route rate limiter
func RouteLimiter(limit *grate.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		if !limit.Allow(c.FullPath()) {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		
		c.Next()
	}
}
