package uid

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func PoisonInit() {
	cryptoRead = func(_ []byte) (int, error) { return 0, assert.AnError }
	_init()
}

// ReseedPRNG seeds the underlying ChaCha8 with zero to make it testably predictable.
func ReseedPRNG() func() {
	old, err := rng.MarshalBinary()
	if err != nil {
		panic(err)
	}
	rng.Seed([32]byte{})
	return func() {
		if err := rng.UnmarshalBinary(old); err != nil {
			panic(err)
		}
	}
}

// SetNowFunc replaces the internal time.Now for unit testing returns a deferrable that undoes this change.
func SetNowFunc(f func() time.Time) func() {
	now = f
	return func() {
		now = time.Now
	}
}
