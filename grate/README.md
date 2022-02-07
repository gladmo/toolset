# grate
Grate is simple rate limit middleware for Go.

## Usage
```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gladmo/toolset/grate"
)

func main() {
	// new rate limiter with ttl
	lim := grate.NewRateLimiter(1, 1, time.Second*10)

	var wg sync.WaitGroup

	// Limit for each key
	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(key string) {
			defer wg.Done()

			// each key allot a *rate.Limiter
			err := lim.Wait(key, context.Background())
			if err != nil {
				panic(err)
			}

			fmt.Println(fmt.Sprintf("idx:%s, now: %s", key, time.Now()))

		}(fmt.Sprint(i % 3))
	}

	wg.Wait()
}
```

### gin middleware
```go
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
```