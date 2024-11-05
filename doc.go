/*
Package uid ...

UUID V7 uses Method 3 (Replace Leftmost Random Bits with Increased Clock Precision) to implement single-node
monotonicity.
*/
package uid

// MaxCanonical is the canonical RFC9562 "Max" UUID.
const MaxCanonical = "ffffffff-ffff-ffff-ffff-ffffffffffff"

// MaxCompact32 is the canonical NCName Compact Base32 "Max" UUID.
const MaxCompact32 = "P777777777777777777777777P"

// MaxCompact64 is the canonical NCName Compact Base64 "Max" UUID.
const MaxCompact64 = "P____________________P"

// MaxJSON is the canonical JSON "Max" UUID.
const MaxJSON = `"ffffffff-ffff-ffff-ffff-ffffffffffff"`

// NilCanonical is the canonical RFC9562 "Nil" UUID.
const NilCanonical = "00000000-0000-0000-0000-000000000000"

// NilCompact32 is the canonical NCName Compact Base32 "Nil" UUID.
const NilCompact32 = "AAAAAAAAAAAAAAAAAAAAAAAAAA"

// NilCompact64 is the canonical NCName Compact Base64 "Nil" UUID.
const NilCompact64 = "AAAAAAAAAAAAAAAAAAAAAA"

// NilJSON is the canonical JSON "Nil" UUID.
const NilJSON = `"00000000-0000-0000-0000-000000000000"`

// VersionNil is the Nil UUID version.
const VersionNil = Version(0) // 0x0

// Version4 is the version of random UUIDs.
const Version4 = Version(4) // 0x4

// Version7 is the version of time-sortable UUIDs.
const Version7 = Version(7) // 0x7

// VersionMax is the Max UUID version.
const VersionMax = Version(15) // 0xf

// Variant9562 is the variant that RFC9562 defines for the types therein.
const Variant9562 = uint8(2) // 0x2

// Version is the RFC9562 UUID Version.
type Version uint8

//nolint:gochecknoglobals // wtb const builtins
var (
	bytesNil = [16]byte{}
	bytesMax = [16]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
	versionsNCNameTable = map[rune]Version{
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
	versionCanonicalTable = map[byte]Version{
		'0': VersionNil, '4': Version4, '7': Version7, 'f': VersionMax, 'F': VersionMax,
	}
)
