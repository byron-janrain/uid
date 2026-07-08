package uid

import (
	"encoding/binary"
	"math/bits"
	"strings"
)

const (
	fiftySeven     = 57
	pythonShortLen = 22
	b57encRef      = "23456789" + "ABCDEFGH" + "JKLMN" + "PQRSTUVWXYZ" + "abcdefghijk" + "mnopqrstuvwxyz"
)

//nolint:gochecknoglobals // wtb const arrays
var (
	// b57dec maps a byte to its base57 value, 0xff for invalid.
	b57dec = func() [256]byte {
		var t [256]byte
		for i := range t {
			t[i] = 0xff
		}
		for i, c := range b57encRef {
			t[c] = byte(i)
		}
		return t
	}()
)

// ToPythonShort returns the Python ShortUUID encoding of u. See https://pypi.org/project/shortuuid.
func ToPythonShort(u UUID) string {
	hi, lo := binary.BigEndian.Uint64(u.b[0:8]), binary.BigEndian.Uint64(u.b[8:16])
	var out [pythonShortLen]byte
	for i := pythonShortLen - 1; i >= 0; i-- {
		var r uint64
		hi, r = hi/fiftySeven, hi%fiftySeven
		lo, r = bits.Div64(r, lo, fiftySeven)
		out[i] = b57encRef[r]
	}
	return string(out[:])
}

// FromPythonShort parses a UUID from Python ShortUUID encoded ps.
func FromPythonShort(ps string) (UUID, bool) {
	ps = strings.TrimSpace(ps)
	if len(ps) != pythonShortLen {
		return UUID{}, false
	}
	var hi, lo uint64
	for i := range pythonShortLen {
		v := b57dec[ps[i]]
		if v == 0xff { //nolint:mnd // lob
			return UUID{}, false
		}
		// (hi,lo) = (hi,lo)*57 + v
		hh, hl := bits.Mul64(hi, fiftySeven)
		lh, ll := bits.Mul64(lo, fiftySeven)
		var c uint64
		lo, c = bits.Add64(ll, uint64(v), 0)
		hi, c = bits.Add64(hl, lh+c, 0) // lh <= 56, no carry from lh+c
		if hh|c != 0 {
			return UUID{}, false // 57^22 > 2^128: value does not fit a UUID
		}
	}
	var out UUID
	binary.BigEndian.PutUint64(out.b[0:8], hi)
	binary.BigEndian.PutUint64(out.b[8:16], lo)
	return out, true
}
