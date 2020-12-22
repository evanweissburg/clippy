package ratelimit

import (
	"time"
)

const (
	maxAccess   = 5
	refreshTime = 30 * time.Second
)

var accessCount = make(map[string]int)

func makeRefreshTicker(refreshSecs time.Duration) {
	ticker := time.NewTicker(refreshSecs)
	for {
		select {
		case <-ticker.C:
			accessCount = make(map[string]int)
		}
	}
}

func init() {
	go makeRefreshTicker(refreshTime)
}

// RequestAccess allows address-based access through rate limiter.
// Returns whether the address is still able to make queries.
func RequestAccess(address string) bool {
	accessCount[address]++
	return accessCount[address] <= maxAccess
}
