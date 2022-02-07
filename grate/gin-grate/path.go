package gin_grate

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/gladmo/toolset/grate"
)

// PathLimiter gin request path rate limiter
func PathLimiter(limit *grate.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		if !limit.Allow(c.Request.RequestURI) {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		
		c.Next()
	}
}
