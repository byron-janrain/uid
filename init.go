package uid

import (
	"math/rand/v2"
	"time"
)

//nolint:gochecknoglobals // manipulatable via functions in export_test
var (
	// rand64 draws from math/rand/v2's per-CPU ChaCha8 generators: error-free, lock-free, and runtime-seeded from
	// OS entropy.
	rand64 = rand.Uint64
	now    = time.Now
)
