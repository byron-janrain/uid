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
func Parse[T []byte | string](src T) (UUID, bool) {
	out, ok := UUID{}, false
	switch len(src) {
	case 38: // json encoded canonical
		src = src[1 : len(src)-1]
		fallthrough
	case 36:
		ok = parseCanonical(&out.b, string(src))
	case 16:
		ok = copy(out.b[:], src) == 16
	case 28: // json encoded ncname32
		src = src[1 : len(src)-1]
		fallthrough
	case 26:
		ok = parseCompact32(&out.b, string(src))
	case 24: // json encoded ncname64
		src = src[1 : len(src)-1]
		fallthrough
	case 22:
		ok = parseCompact64(&out.b, string(src))
	}
	if !ok {
		return UUID{}, false
	}
	return out, true
}

//nolint:gochecknoglobals // localit
var offsets = [8]byte{0, 4, 9, 14, 19, 24, 28, 32}

//nolint:cyclop // parser...
func parseCanonical(tgt *[16]byte, src string) bool {
	if src[8] != '-' || src[13] != '-' || src[18] != '-' || src[23] != '-' {
		return false
	}
	switch src[14] {
	case '0':
		if src == NilCanonical {
			*tgt = bytesNil
			return true
		}
		return false
	case 'f', 'F':
		if strings.ToLower(src) == MaxCanonical {
			*tgt = bytesMax
			return true
		}
		return false
	case '4', '7': // decode
	default:
		return false
	}
	for x, i := range offsets {
		n, err := hex.Decode((*tgt)[x*2:(x*2)+2], []byte(src[i:i+4])) //nolint:mnd // lob
		if err != nil || n != 2 {
			return false
		}
	}
	return checkVariant(tgt)
}

//nolint:gochecknoglobals // shared locality of behavior
var b32decoder = base32.StdEncoding.WithPadding(base32.NoPadding)

func parseCompact32(tgt *[16]byte, src string) bool {
	src = strings.ToUpper(src)
	v, ok := versionsNCNameTable[rune(src[0])]
	if !ok {
		return false // unsupported version
	}
	// s/c nil
	if v == VersionNil {
		if src != NilCompact32 {
			return false
		}
		*tgt = bytesNil
		return true
	}
	// s/c max
	if v == VersionMax {
		if src != MaxCompact32 {
			return false
		}
		*tgt = bytesMax
		return true
	}
	// not Nil or Max, decode with padding v4/v7
	_, err := b32decoder.Decode((*tgt)[:], []byte(src + "A")[1:])
	if err != nil {
		return false
	}
	tgt[15] <<= 1 // unshift bookend
	unshift(tgt, uint32(v))
	return checkVariant(tgt)
}

func parseCompact64(tgt *[16]byte, src string) bool {
	runes := []rune(src)
	v, ok := versionsNCNameTable[runes[0]]
	if !ok {
		return false
	}
	// s/c nil
	if v == VersionNil {
		if src != NilCompact64 {
			return false
		}
		*tgt = bytesNil
		return true
	}
	// s/c max
	if v == VersionMax {
		if src != MaxCompact64 {
			return false
		}
		*tgt = bytesMax
		return true
	}
	// not Nil or Max, decode with padding
	runes[21] = unicode.ToUpper(runes[21])
	_, err := base64.RawURLEncoding.Decode((*tgt)[:], []byte(string(runes) + "A")[1:])
	if err != nil {
		return false
	}
	tgt[15] <<= 2
	unshift(tgt, uint32(v))
	return checkVariant(tgt)
}

func checkVariant(src *[16]byte) bool { return variant(*src) == Variant9562 }

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
