package main

import (
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/gladmo/toolset/grate"
	gingrate "github.com/gladmo/toolset/grate/gin-grate"
)

func main() {
	r := gin.Default()
	
	// use gin grate ip limiter
	r.Use(gingrate.IPLimiter(grate.NewRateLimiter(10, 1, time.Minute)))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	
	err := r.Run("0.0.0.0:8910")
	if err != nil {
		panic(err)
	}
}
