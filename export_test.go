package uid

import "time"

// ReseedPRNG seeds the underlying ChaCha8 with zero to make is testably predictable.
func ReseedPRNG() { rng.Seed([32]byte{}) }

// SetNowFunc replaces the internal time.Now for unit testing returns a function that undoes this change.
func SetNowFunc(f func() time.Time) func() {
	now = f
	return func() {
		now = time.Now
	}
}
