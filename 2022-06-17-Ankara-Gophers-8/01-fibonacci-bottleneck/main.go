package main

import "fmt"

func fibonacci(n uint64) uint64 {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	fmt.Printf("I'm doing some work here...\n")
	fmt.Printf("...and some more work here...\n")
	_ = fibonacci(45)
	fmt.Printf("...aaand I'm done here.\n")
}
