package messages

import (
	"sync"
	"time"
)

// TallyInterval is the fastest live tallies are sent to a room
const TallyInterval = 250 * time.Millisecond

var throttleMux sync.Mutex
var throttleOpen = make(map[string]bool)
var throttlePending = make(map[string]func())

// throttled sends at once if nothing was sent for this key in the last interval, otherwise
// it keeps only the latest send and runs it when the interval ends. Send functions carry
// their own copied data, so they never need the game lock.
func throttled(key string, interval time.Duration, send func()) {
	throttleMux.Lock()
	if throttleOpen[key] {
		throttlePending[key] = send
		throttleMux.Unlock()
		return
	}
	throttleOpen[key] = true
	throttleMux.Unlock()

	send()
	time.AfterFunc(interval, func() { flushThrottled(key, interval) })
}

func flushThrottled(key string, interval time.Duration) {
	throttleMux.Lock()
	send, pending := throttlePending[key]
	delete(throttlePending, key)
	if !pending {
		delete(throttleOpen, key)
		throttleMux.Unlock()
		return
	}
	throttleMux.Unlock()

	send()
	time.AfterFunc(interval, func() { flushThrottled(key, interval) })
}
