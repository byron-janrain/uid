package uid

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"slices"
	"strings"
	"unicode"
)

// Parse attempts to parse `src` into a UUID and returns the parsed UUID and `true` on success.
// On failure, Parse returns the Nil UUID and `false`.
//
//nolint:mnd // locality of behavior
func Parse(src string) (UUID, bool) {
	ln := len(src)
	switch ln {
	case 38: // canonical JSON encoded or non-canonical boundaries.
		src = src[1 : ln-1]
		fallthrough
	case 36:
		return parseCanonical(src)
	case 16:
		return parseBytes([]byte(src))
	case 28: // json encoded ncname32
		src = src[1 : ln-1]
		fallthrough
	case 26:
		return parseCompact32(src)
	case 24: // json encoded ncname64
		src = src[1 : ln-1]
		fallthrough
	case 22:
		return parseCompact64(src)
	}
	return UUID{}, false
}

//nolint:cyclop // parsers...
func canonicalV(s string) Version {
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return versionBad // malformed
	}
	varR := rune(s[19])
	switch rune(s[14]) {
	case '4':
		switch varR {
		case '8', '9', 'A', 'a', 'b', 'B':
			return Version4
		}
	case '7':
		switch varR {
		case '8', '9', 'A', 'a', 'b', 'B':
			return Version7
		}
	case '0':
		if s == NilCanonical {
			return VersionNil
		}
	case 'f', 'F':
		if strings.ToLower(s) == MaxCanonical {
			return VersionMax
		}
	}
	return versionBad
}

func bytesV(b []byte) Version {
	vrsn := Version(b[6] >> 4) //nolint:mnd // lob
	switch vrsn {              //nolint:exhaustive // golf
	case Version4, Version7:
		if Variant(b[8]>>6) == Variant9562 { //nolint:mnd // lob
			return vrsn
		}
	case VersionNil:
		if slices.Equal(b, bytesNil[:]) {
			return VersionNil
		}
	case VersionMax:
		if slices.Equal(b, bytesMax[:]) {
			return VersionMax
		}
	}
	return versionBad
}

func ncn64V(s string) Version {
	varR := rune(s[21])
	switch rune(s[0]) {
	case 'E':
		switch varR {
		case 'I', 'J', 'K', 'L':
			return Version4
		}
	case 'H':
		switch varR {
		case 'I', 'J', 'K', 'L':
			return Version7
		}
	case 'A':
		if s == NilCompact64 {
			return VersionNil
		}
	case 'P':
		if s == MaxCompact64 {
			return VersionMax
		}
	}
	return versionBad
}

func ncn32V(s string) Version {
	varR := rune(s[25])
	switch rune(s[0]) {
	case 'E', 'e':
		switch varR {
		case 'i', 'I', 'j', 'J', 'k', 'K', 'l', 'L':
			return Version4
		}
	case 'H', 'h':
		switch varR {
		case 'i', 'I', 'j', 'J', 'k', 'K', 'l', 'L':
			return Version7
		}
	case 'A', 'a':
		if strings.ToUpper(s) == NilCompact32 {
			return VersionNil
		}
	case 'P', 'p':
		if strings.ToUpper(s) == MaxCompact32 {
			return VersionMax
		}
	}
	return versionBad
}

//nolint:gochecknoglobals // ref
var canonOffsets = [16]byte{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34}

func parseCanonical(src string) (UUID, bool) {
	switch canonicalV(src) { //nolint:exhaustive // golf
	case versionBad:
		return UUID{}, false
	case VersionNil:
		return UUID{bytesNil}, true
	case VersionMax:
		return UUID{bytesMax}, true
	}
	// valid version+variant but not nil/max
	tgt := [16]byte{}
	for i, x := range canonOffsets {
		if !c2b(&tgt[i], src[x], src[x+1]) {
			return UUID{}, false
		}
	}
	return UUID{tgt}, true
}

func parseBytes(b []byte) (UUID, bool) {
	switch bytesV(b) { //nolint:exhaustive // golf
	case versionBad:
		return UUID{}, false
	case VersionNil:
		return UUID{bytesNil}, true
	case VersionMax:
		return UUID{bytesMax}, true
	}
	// valid version+variant but not nil/max
	var out UUID
	copy(out.b[:], b)
	return out, true
}

func parseCompact32(src string) (UUID, bool) {
	v := ncn32V(src)
	if v == versionBad {
		return UUID{}, false
	}
	// s/c nil
	if v == VersionNil {
		return UUID{bytesNil}, true
	}
	// s/c max
	if v == VersionMax {
		return UUID{bytesMax}, true
	}
	// not Nil or Max, decode with padding v4/v7
	var out UUID
	_, err := b32decoder.Decode(out.b[:], []byte(strings.ToUpper(src) + "A")[1:])
	if err != nil {
		return UUID{}, false
	}
	out.b[15] <<= 1 // unshift bookend
	unshift(&out.b, uint32(v))
	return out, true
}

func parseCompact64(src string) (UUID, bool) {
	v := ncn64V(src)
	if v == versionBad {
		return UUID{}, false
	}
	// s/c nil
	if v == VersionNil {
		return UUID{bytesNil}, true
	}
	// s/c max
	if v == VersionMax {
		return UUID{bytesMax}, true
	}
	// not Nil or Max, decode with padding
	runes := []rune(src)
	runes[21] = unicode.ToUpper(runes[21])
	var out UUID
	_, err := base64.RawURLEncoding.Decode(out.b[:], []byte(string(runes) + "A")[1:])
	if err != nil {
		return UUID{}, false
	}
	out.b[15] <<= 2
	unshift(&out.b, uint32(v))
	return out, true
}

//nolint:mnd // locality of behavior
func unshift(tgt *[16]byte, version uint32) {
	version &= 0xf
	ints := [4]uint32{
		binary.BigEndian.Uint32(tgt[0:4]),
		binary.BigEndian.Uint32(tgt[4:8]),
		binary.BigEndian.Uint32(tgt[8:12]),
		binary.BigEndian.Uint32(tgt[12:16]),
	}
	variant := (ints[3] & 0xf0) << 24
	ints[3] >>= 8
	ints[3] |= (ints[2] & 0xff) << 24
	ints[2] >>= 8
	ints[2] |= ((ints[1] & 0xf) << 24) | variant
	ints[1] = (ints[1] & 0xffff0000) | (version << 12) | ((ints[1] >> 4) & 0xfff)
	binary.BigEndian.PutUint32(tgt[0:4], ints[0])
	binary.BigEndian.PutUint32(tgt[4:8], ints[1])
	binary.BigEndian.PutUint32(tgt[8:12], ints[2])
	binary.BigEndian.PutUint32(tgt[12:16], ints[3])
}

// c2b sets v to the possibly hex decoded value of characters c1 and c2. returns false if v is invalid.
//
//nolint:mnd // lob
func c2b(b *byte, c1, c2 byte) bool {
	c1, c2 = c2h[c1], c2h[c2]
	*b = (c1 << 4) | c2
	return c1|c2 != 0xff
}

//nolint:gochecknoglobals // ref
var (
	// map ascii -> hex (0xff is invalid). lookup table method borrowed from google/uuid.
	c2h = [256]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, // valid '0'-'9'
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xa, 0xb, 0xc, 0xd, 0xe, 0xf, // valid 'A'-'F'
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xa, 0xb, 0xc, 0xd, 0xe, 0xf, // valid 'a'-'f'
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
	b32decoder = base32.StdEncoding.WithPadding(base32.NoPadding)
)
