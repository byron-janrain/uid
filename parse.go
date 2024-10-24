package uid

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
	"unicode"
)

var b32decoder = base32.StdEncoding.WithPadding(base32.NoPadding)

func Parse[T string | []rune | []byte | [16]byte](src T) (id UUID, err error) {
	b := &id.b
	switch t := any(src).(type) {
	case string, []byte, []rune:
		err = parseString(b, t.(string))
	case [16]byte:
		// err = parseArray(b, t)
		// case []byte, []rune:
		// 	err = parseSlice(b, t)
	}
	return id, err
}

func parseString(tgt *[16]byte, src string) error {
	switch len(src) {
	case 26: // compact32
		return parseCompact32(tgt, src)
	case 22: // compact64
		return parseCompact64(tgt, src)
	case 36: // canonical
		return parseCanon(tgt, src)
	case 45: // urn
		return parseCanon(tgt, src[9:]) // strip off urn:uuid:
	}
	return errors.New("unrecognized UUID format or version")
}

// func parseSlice(tgt *[16]byte, src []byte) error {
// 	const rawlen = 16
// 	if len(src) == rawlen {
// 		var a [rawlen]byte
// 		copy(a[:], src[:rawlen])
// 		return parseArray(tgt, a)
// 	}
// 	// if it's not 16 bytes, it's a string form
// 	return parseString(tgt, string(src))
// }

// func parseArray(tgt *[16]byte, src [16]byte) error {
// 	*tgt = src
// 	return nil
// }

func parseCanon(tgt *[16]byte, src string) error {
	if _, err := hex.Decode((*tgt)[:], []byte(strings.ReplaceAll(src, "-", ""))); err != nil {
		panic(err)
	}
	return nil
}

func parseCompact32(tgt *[16]byte, src string) error {
	src = strings.ToUpper(src) // canonicalize
	version, ok := versionsNCNameTable[rune(src[0])]
	if !ok {
		panic("fix this")
	}
	pad := "A"
	if version == VersionMax {
		pad = "P"
	}
	src += pad
	_, err := b32decoder.Decode((*tgt)[:], []byte(src)[1:])
	if err != nil {
		panic(src[1:] + " " + err.Error())
	}
	tgt[15] <<= 1 // unshift bookend
	unshift(tgt, uint32(version))
	return nil
}

func parseCompact64(tgt *[16]byte, src string) error {
	runes := []rune(src)
	version, ok := versionsNCNameTable[runes[0]]
	if !ok {
		panic("change to error")
	}
	runes[21] = unicode.ToUpper(runes[21])
	pad := "A"
	if uint8(version) == VersionMax {
		pad = "P"
	}
	src = string(runes) + pad
	_, err := base64.RawURLEncoding.Decode((*tgt)[:], []byte(src)[1:])
	if err != nil {
		panic(err)
	}
	tgt[15] <<= 2
	unshift(tgt, uint32(version))
	return nil
}

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
