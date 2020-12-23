package ratelimit

import (
	"fmt"
	"testing"
)

func TestTooManyQueries(t *testing.T) {
	address := "addr"
	for i := 0; i < 5; i++ {
		if !RequestAccess(address) {
			fmt.Printf("we blocked it with %d\n", i)
			t.Fatalf("Access blocked incorrectly for address %s\n", address)
		}
	}
	if RequestAccess(address) {
		t.Fatalf("Address granted incorrectly for address %s\n", address)
	}
}
