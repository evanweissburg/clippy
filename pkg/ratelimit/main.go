package ratelimit

import (
	"sync"
	"time"
)

const (
	maxAccess   = 5
	refreshTime = 30 * time.Second
)

var accessCount = make(map[string]int)
var mu sync.Mutex

func makeRefreshTicker(refreshSecs time.Duration) {
	ticker := time.NewTicker(refreshSecs)
	for {
		select {
		case <-ticker.C:
			mu.Lock()
			accessCount = make(map[string]int)
			mu.Unlock()
		}
	}
}

func init() {
	go makeRefreshTicker(refreshTime)
}

// RequestAccess allows address-based access through rate limiter.
// Returns whether the address is still able to make queries.
func RequestAccess(address string) bool {
	mu.Lock()
	defer mu.Unlock()
	accessCount[address]++
	return accessCount[address] <= maxAccess
}
