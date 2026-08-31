package uid

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"slices"
)

// UUID is a UUID as defined by RFC...
// Underlying array is unexported for immutability. UUID is comparable using `==`.
// The zero value is Nil UUID.
//
//nolint:recvcheck // only unserializers should (temporarily) have pointers.
type UUID struct{ b [16]byte }

// Version returns u's version.
func (u UUID) Version() Version { return Version(u.b[6] >> 4) } //nolint:mnd // lob

// Variant is u's variant.
func (u UUID) Variant() Variant {
	if u.Version() == VersionMax {
		return VariantMax
	}
	return Variant(u.b[8] >> 6) //nolint:mnd // lob
}

// Bytes returns a copy of u's raw bytes.
func (u UUID) Bytes() []byte { return u.b[:] } // copy

// MarshalBinary implements encoding.BinaryMarshaler. Never returns errors.
func (u UUID) MarshalBinary() ([]byte, error) { return u.b[:], nil }

// UnmarshalBinary implements encoding.BinaryUnmarshaler. Returns ErrInvalid for malformed input.
func (u *UUID) UnmarshalBinary(b []byte) error {
	if id, ok := Parse(string(b)); ok {
		*u = id
		return nil
	}
	return ErrInvalid
}

// canonLen is the length of the canonical representation.
const canonLen = len(NilCanonical)

// String implements fmt.Stringer. Returns canonical RFC-4122 representation.
func (u UUID) String() string {
	var buf [canonLen]byte // stack: encodeTo keeps no reference, so only the string below escapes
	u.encodeTo(buf[:])
	return string(buf[:])
}

// MarshalText implements encoding.TextMarshaler. Never returns errors.
func (u UUID) MarshalText() ([]byte, error) {
	buf := make([]byte, canonLen)
	u.encodeTo(buf)
	return buf, nil
}

// AppendText implements encoding.TextAppender. Never returns errors.
func (u UUID) AppendText(b []byte) ([]byte, error) { return u.appendText(b), nil }

// UnmarshalText implements encoding.TextUnmarshaler. Returns ErrInvalid for malformed input.
func (u *UUID) UnmarshalText(b []byte) error {
	if id, ok := Parse(string(b)); ok {
		*u = id
		return nil
	}
	return ErrInvalid
}

// MarshalJSON implements encoding/json.Marshaler. Never returns errors.
func (u UUID) MarshalJSON() ([]byte, error) { return []byte(`"` + u.String() + `"`), nil }

// UnmarshalJSON implements encoding/json.Unmarshaler. Returns ErrInvalid for malformed input.
func (u *UUID) UnmarshalJSON(b []byte) error {
	if id, ok := Parse(string(b)); ok {
		*u = id
		return nil
	}
	return ErrInvalid
}

// Compact32 returns NCName Base32 representation.
func (u UUID) Compact32() string {
	b := u.shifted()
	b[15] >>= 1
	//nolint:mnd // lob
	return string(u.Version()+65) + base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b[:])[0:25]
}

// Compact64 returns NCName Base64 representation.
func (u UUID) Compact64() string {
	b := u.shifted()
	b[15] >>= 2
	//nolint:mnd // lob
	return string(u.Version()+65) + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b[:])[0:21]
}

// Nil constructs a Nil UUID (all 0).
func Nil() UUID { return UUID{bytesNil /*copy*/} }

// IsNil returns true when u is the Nil UUID.
func (u UUID) IsNil() bool { return u.b == bytesNil } // compare to zero array is highly optimized

// Max constructs a Max UUID (all F).
func Max() UUID { return UUID{bytesMax /*copy*/} }

// IsMax returns true when u is the Max UUID.
func (u UUID) IsMax() bool { return u.b == bytesMax }

// b2h maps a byte to its two lowercase hex digits, the encoding counterpart of c2h.
//
//nolint:gochecknoglobals // lookup table
var b2h = func() [256][2]byte {
	const digits = "0123456789abcdef"
	var t [256][2]byte
	for i := range t {
		t[i][0], t[i][1] = digits[i>>4], digits[i&0x0f]
	}
	return t
}()

// encodeTo writes u's canonical representation over the first canonLen bytes of buf. canonOffsets places each byte,
// so the four dashes fall in the gaps it skips.
func (u UUID) encodeTo(buf []byte) {
	_ = buf[canonLen-1] // hoist the bounds check out of the loop
	buf[8], buf[13], buf[18], buf[23] = '-', '-', '-', '-'
	for i, x := range canonOffsets {
		buf[x], buf[x+1] = b2h[u.b[i]][0], b2h[u.b[i]][1]
	}
}

// appendText appends u's canonical representation to b. Grow-then-extend reserves the bytes encodeTo covers in full.
func (u UUID) appendText(b []byte) []byte {
	i := len(b)
	b = slices.Grow(b, canonLen)[:i+canonLen]
	u.encodeTo(b[i:]) // reslice after growing: growth may have moved the backing array
	return b
}

//nolint:nonamedreturns,mnd // golf, locality of behavior
func (u UUID) shifted() (out [16]byte) {
	ints := [4]uint32{
		binary.BigEndian.Uint32(u.b[0:4]),
		binary.BigEndian.Uint32(u.b[4:8]),
		binary.BigEndian.Uint32(u.b[8:12]),
		binary.BigEndian.Uint32(u.b[12:16]),
	}
	variant := (ints[2] & 0xf0000000) >> 24
	ints[1] = (ints[1] & 0xffff0000) | ((ints[1] & 0x00000fff) << 4) | (ints[2] & 0x0fffffff >> 24)
	ints[2] = (ints[2]&0x00ffffff)<<8 | ints[3]>>24
	ints[3] = (ints[3] << 8) | variant
	binary.BigEndian.PutUint32(out[0:4], ints[0])
	binary.BigEndian.PutUint32(out[4:8], ints[1])
	binary.BigEndian.PutUint32(out[8:12], ints[2])
	binary.BigEndian.PutUint32(out[12:16], ints[3])
	return
}

// Compare implements slices.SortFunc for the UUID type. v7 UUIDs sort by embedded time (unix_ts_ms and rand_a);
// random bits break ties so distinct UUIDs never compare equal.
func Compare(a, b UUID) int { return bytes.Compare(a.b[:], b.b[:]) }
