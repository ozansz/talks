package main

import "testing"

func fibonacci(n uint64) uint64 {
	// Allocate n+1 uint64 values in heap, which is
	// (n+1)*8 bytes in total.
	cache := make([]uint64, n+1)

	cache[0] = 0
	cache[1] = 1

	var i uint64
	for i = 2; i <= n; i++ {
		cache[i] = cache[i-1] + cache[i-2]
	}

	return cache[n]
}

func TestFibonacci(t *testing.T) {
	tests := []struct {
		n    uint64
		want uint64
	}{
		{50, 12586269025},
		{10000000, 10047910021417012027}, // It's actually 1953282128707757731632014947596256332443542996591873396953405194571625257887015694766641987634150146128879524335220236084625510912019560233744... and so on.
	}
	for _, test := range tests {
		if got := fibonacci(test.n); got != test.want {
			t.Errorf("fibonacci(%d) = %d, want %d", test.n, got, test.want)
		}
	}
}
