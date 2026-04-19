package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/atharvadani9/rate-limiter/ratelimiter"
)

func main() {
	rl := ratelimiter.New(3, 5*time.Second)
	defer rl.Stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", rl.Middleware(mux))
}
