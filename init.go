package uid

import (
	crand "crypto/rand"
	"math/rand/v2"
	"time"
)

//nolint:gochecknoglobals // manipulatable via functions in export_test
var (
	cryptoRead = crand.Read
	rng        *rand.ChaCha8
	now        = time.Now
)

//nolint:gochecknoinits // indirect for testability
func init() { _init() }

// use real crypto/rand to initialize the ChaCha8 generator.
func _init() {
	var seed [32]byte
	if _, err := cryptoRead(seed[:]); err != nil {
		panic("unable to initialize seed from crypto/rand") // untestable
	}
	rng = rand.NewChaCha8(seed)
}
