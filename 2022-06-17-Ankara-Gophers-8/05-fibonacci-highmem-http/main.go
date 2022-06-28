package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"
)

func fibonacci(n uint64) uint64 {
	cache := make([]uint64, n+1)

	cache[0] = 0
	cache[1] = 1

	var i uint64
	for i = 2; i <= n; i++ {
		cache[i] = cache[i-1] + cache[i-2]
	}

	return cache[n]
}

func main() {
	go func() {
		// Just need to navigate to http://localhost:6060/debug/pprof
		// Or: go tool pprof http://localhost:6060/debug/pprof/heap and so on...
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	for {
		_ = fibonacci(10000000)
	}
}
