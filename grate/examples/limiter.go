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
