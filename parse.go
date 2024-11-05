package uid

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"unicode"
)

// Parse attempts to parse `src` into a UUID and returns the parsed UUID and `true` on success.
// On failure, Parse returns the Nil UUID and `false`.
//
//nolint:mnd // locality of behavior
func Parse(src string) (UUID, error) {
	out := UUID{}
	var err error
	switch len(src) {
	case 36:
		err = parseCanonical(&out.b, src)
	case 16: // raw [16]byte slice
		copy(out.b[:], []byte(src))
	case 26:
		err = parseCompact32(&out.b, src)
	case 22:
		err = parseCompact64(&out.b, src)
		// case 45: // urn not yet supported
		// 	err = parseCanon(b, src[9:]) // strip off urn:uuid:
	default:
		return out, ParseError{}
	}
	if err != nil {
		out.b = bytesNil
	}
	return out, err
}

// Must is a helper to wrap the parser in a panic-on-error handler. Useful for testing.
func Must(u UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return u
}

func parseCanonical(tgt *[16]byte, src string) error {
	v, ok := versionCanonicalTable[src[14]]
	if !ok {
		return ParseError{}
	}
	// s/c nil
	if v == VersionNil {
		if src != NilCanonical {
			return ParseError{}
		}
		*tgt = bytesNil
		return nil
	}
	// s/c max
	if v == VersionMax {
		if src != MaxCanonical {
			return ParseError{}
		}
		*tgt = bytesMax
		return nil
	}
	// not Nil or Max, decode without dashes.
	if _, err := hex.Decode((*tgt)[:], []byte(strings.ReplaceAll(src, "-", ""))); err != nil {
		return ParseError{err}
	}
	return checkVariant(tgt)
}

//nolint:gochecknoglobals // shared locality of behavior
var b32decoder = base32.StdEncoding.WithPadding(base32.NoPadding)

func parseCompact32(tgt *[16]byte, src string) error {
	src = strings.ToUpper(src)
	v, ok := versionsNCNameTable[rune(src[0])]
	if !ok {
		// unsupported version
		return ParseError{}
	}
	// s/c nil
	if v == VersionNil {
		if src != NilCompact32 {
			return ParseError{}
		}
		*tgt = bytesNil
		return nil
	}
	// s/c max
	if v == VersionMax {
		if src != MaxCompact32 {
			return ParseError{}
		}
		*tgt = bytesMax
		return nil
	}
	// not Nil or Max, decode with padding v4/v7
	_, err := b32decoder.Decode((*tgt)[:], []byte(src + "A")[1:])
	if err != nil {
		return ParseError{err}
	}
	tgt[15] <<= 1 // unshift bookend
	unshift(tgt, uint32(v))
	return checkVariant(tgt)
}

func parseCompact64(tgt *[16]byte, src string) error {
	runes := []rune(src)
	v, ok := versionsNCNameTable[runes[0]]
	if !ok {
		return ParseError{}
	}
	// s/c nil
	if v == VersionNil {
		if src != NilCompact64 {
			return ParseError{}
		}
		*tgt = bytesNil
		return nil
	}
	// s/c max
	if v == VersionMax {
		if src != MaxCompact64 {
			return ParseError{}
		}
		*tgt = bytesMax
		return nil
	}
	// not Nil or Max, decode with padding
	runes[21] = unicode.ToUpper(runes[21])
	_, err := base64.RawURLEncoding.Decode((*tgt)[:], []byte(string(runes) + "A")[1:])
	if err != nil {
		return ParseError{err}
	}
	tgt[15] <<= 2
	unshift(tgt, uint32(v))
	return checkVariant(tgt)
}

func checkVariant(src *[16]byte) error {
	if variant(*src) != Variant9562 {
		return ParseError{}
	}
	return nil
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
