package uid

import "encoding/binary"

// NewV4 returns a new v4 UUID.
func NewV4() UUID {
	var b [16]byte
	binary.LittleEndian.PutUint64(b[0:8], rand64())
	binary.LittleEndian.PutUint64(b[8:16], rand64())
	// version, variant
	b[6], b[8] = (b[6]&0x0f)|0x40, (b[8]&0x3f)|0x80 //nolint:mnd // lob

	return UUID{b}
}
