package main

import (
	"fmt"
	"time"

	"github.com/atharvadani9/rate-limiter/ratelimiter"
)

func main() {
	rl := ratelimiter.New(3, 5*time.Second)
	defer rl.Stop()

	fmt.Printf("rate limiter created with limit 3 and window duration 5s\n")
	allowed := 0
	denied := 0
	for i := range 10 {
		if rl.Allow("key1") {
			allowed++
			fmt.Printf("allowed %d\n", i)
		} else {
			denied++
			fmt.Printf("denied %d\n", i)
		}
		time.Sleep(time.Second)
		s := rl.Status("key1")
		fmt.Printf("status: remaining=%d, resets at=%s\n", s.Remaining, s.ResetAt.Format(time.RFC3339))
	}
	fmt.Printf("allowed %d requests\n", allowed)
	fmt.Printf("denied %d requests\n", denied)
}
