package uid

import (
	"math/rand/v2"
	"time"
)

// ReseedPRNG swaps the PRNG for a zero-seeded ChaCha8 to make output testably predictable.
func ReseedPRNG() func() {
	old := rand64
	rand64 = rand.NewChaCha8([32]byte{}).Uint64
	return func() { rand64 = old }
}

// ResetV7Strict zeroes the strict-monotonicity state so tests driving fake clocks don't stall later real-clock tests.
func ResetV7Strict() {
	mux.Lock()
	lastMS, lastRA = 0, 0
	mux.Unlock()
}

// SetNowFunc replaces the internal time.Now for unit testing returns a deferrable that undoes this change.
func SetNowFunc(f func() time.Time) func() {
	now = f
	return func() {
		now = time.Now
	}
}
