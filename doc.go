package uid

import (
	crand "crypto/rand"
	"math/rand/v2"
	"time"
)

// ref
const (
	// well-known strings
	MaxCanonical = "ffffffff-ffff-ffff-ffff-ffffffffffff" // case insensitive
	MaxCompact32 = "P777777777777777777777777P"           // case insensitive
	MaxCompact64 = "P____________________P"               // case sensitive
	NilCanonical = "00000000-0000-0000-0000-000000000000"
	NilCompact32 = "AAAAAAAAAAAAAAAAAAAAAAAAAA" // case insensitive
	NilCompact64 = "AAAAAAAAAAAAAAAAAAAAAA"     // case sensitive
	// versions
	VersionNil = uint8(0)
	Version4   = uint8(4)
	Version7   = uint8(7)
	VersionMax = uint8(15)
	// variants
	Variant9562 = uint8(2)
)

// wtb const arrays
var (
	bytesNil            = [16]byte{}
	bytesMax            = [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	versionsNCNameTable = map[rune]uint8{
		'A': VersionNil, // nil
		// 'B': 1, // v1 not supported yet
		// 'C': 2, // v2 not supported yet
		// 'D': 3, // v3 not supported yet
		'E': Version4, // v4
		// 'F': 5, // v5 not supported yet
		// 'G': 6, // v6 not supported yet
		'H': Version7, // v7
		// 'I': 8, // v8 not supported yet
		'P': VersionMax, // max
	}
)

// pkg runtime
var rng *rand.ChaCha8 // manipulatable via functions in export_test
var now = time.Now    // manipulatable via functions in export_test

// use crypto/rand to initialize the ChaCha8 generator
func init() {
	var seed [32]byte
	n, err := crand.Read(seed[:])
	if err != nil || n != 32 {
		panic("unable to initialize seed from crypto/rand")
	}
	rng = rand.NewChaCha8(seed)
}
