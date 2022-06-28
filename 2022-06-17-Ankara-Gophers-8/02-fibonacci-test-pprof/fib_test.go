package main

import "testing"

func fibonacci(n uint64) uint64 {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
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
