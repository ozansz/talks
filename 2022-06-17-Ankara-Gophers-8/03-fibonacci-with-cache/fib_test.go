package main

import "testing"

// Allocate 51 uint64 values in stack, which is
// 408 bytes in total.
var cache [51]uint64

func fibonacci(n uint64) uint64 {
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
		{0, 0},
		{1, 1},
		{2, 1},
		{15, 610},
		{50, 12586269025},
	}
	for _, test := range tests {
		if got := fibonacci(test.n); got != test.want {
			t.Errorf("fibonacci(%d) = %d, want %d", test.n, got, test.want)
		}
	}
}
