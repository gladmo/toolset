package gin_grate

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/gladmo/toolset/grate"
)

// BasicAuthLimiter gin request basic auth user rate limiter
func BasicAuthLimiter(limit *grate.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		username, _, ok := c.Request.BasicAuth()
		if ok && !limit.Allow(username) {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		
		c.Next()
	}
}
