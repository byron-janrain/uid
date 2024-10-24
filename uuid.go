package uid

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"time"
)

// UUID is a UUID as defined by RFC...
// Underlying array is unexported for immutability. UUID is comparable using `==`.
// The zero value is == Nil UUID.
type UUID struct{ b [16]byte }

// Version returns u's version.
func (u UUID) Version() uint8 { return u.b[6] >> 4 }

func (u UUID) Variant() uint8 { return (u.b[8] >> 6) }

// Bytes returns a copy of u's raw bytes.
func (u UUID) Bytes() []byte { return u.b[:] } // copy

// String implements fmt.Stringer. Returns canonical RFC-4122 representation.
func (u UUID) String() string {
	buf := make([]byte, 36)
	buf[8], buf[13], buf[18], buf[23] = '-', '-', '-', '-'
	hex.Encode(buf[0:8], u.b[0:4])
	hex.Encode(buf[9:13], u.b[4:6])
	hex.Encode(buf[14:18], u.b[6:8])
	hex.Encode(buf[19:23], u.b[8:10])
	hex.Encode(buf[24:], u.b[10:])
	return string(buf[:])
}

// Compact returns draft-taylor-uuid-ncname-03 compact Base32 representation.
func (u UUID) Compact32() string {
	b := u.shifted()
	b[15] = b[15] >> 1
	return string(u.Version()+65) + base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b[:])[0:25]
}

func (u UUID) Compact64() string {
	b := u.shifted()
	b[15] = b[15] >> 2
	return string(u.Version()+65) + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b[:])[0:21]
}

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

// MarshalBinary implements encoding.BinaryMarshaler. Never returns errors.
func (u UUID) MarshalBinary() ([]byte, error) { return u.b[:], nil }

// MarshalText implements encoding.TextMarshaler. Never returns errors.
func (u UUID) MarshalText() ([]byte, error) { return []byte(u.String()), nil }

// MarshalJSON implements encoding/json.Marshaler. Never returns errors.
func (u UUID) MarshalJSON() ([]byte, error) { return []byte(`"` + u.String() + `"`), nil }

// UnmarshalBinary implement encoding.BinaryUnmarshaler.
func (u *UUID) UnmarshalBinary(b []byte) (err error) {
	*u, err = Parse(b)
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *UUID) UnmarshalText(b []byte) (err error) {
	*u, err = Parse(b)
	return
}

// UnmarshalJSON implements encoding/json.Unmarshaler.
func (u *UUID) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return
	}
	*u, err = Parse(s)
	return
}

func Must(u UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return u
}

// Nil constructs a Nil UUID (all 0).
func Nil() UUID { return UUID{bytesNil /*copy*/} }

// IsNil returns true when u is the Nil UUID.
func (u UUID) IsNil() bool { return u.b == bytesNil } // compare to zero array is highly optimized

// NewV4 returns a new v4 UUID.
func NewV4() UUID {
	var b [16]byte
	rng.Read(b[:])                                  //nolint:errcheck // does not return errors
	b[6], b[8] = (b[6]&0x0f)|0x40, (b[8]&0x3f)|0x80 // version, variant
	return UUID{b}
}

// Max constructs a Max UUID (all F).
func Max() UUID { return UUID{bytesMax /*copy*/} }

// IsMax returns true when u is the Max UUID.
func (u UUID) IsMax() bool { return u.b == bytesMax }

// NewV7 constructs a new v7 UUID.
func NewV7() UUID {
	var b [16]byte
	t := now().UnixMilli()
	if t < 0 {
		panic("v7 UUID does not support time before epoch")
	}
	b[0], b[1], b[2], b[3], b[4], b[5] = //first 6 bytes: bigendian uint64 without math.Big
		byte(t>>40), byte(t>>32), byte(t>>24), byte(t>>16), byte(t>>8), byte(t)
	rng.Read(b[6:])                                 //nolint:errcheck // never returns errors
	b[6], b[8] = (b[6]&0x0F)|0x70, (b[8]&0x3F)|0x80 // version, variant
	return UUID{b}
}

// Time returns the embedded timestamp of UUID. For non-V7 zero(time.Time) is returned. If you don't pre-check version
// use `.IsZero()` to ensure time is "real".
func (u UUID) Time() time.Time {
	if u.Version() != Version7 {
		return time.Time{}
	}
	ms := uint64(u.b[5]) | uint64(u.b[4])<<8 | uint64(u.b[3])<<16 | uint64(u.b[2])<<24 | uint64(u.b[1])<<32 |
		uint64(u.b[0])<<40
	return time.UnixMilli(int64(ms))
}
