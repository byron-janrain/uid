package uid

// NewV4 returns a new v4 UUID.
func NewV4() UUID {
	var b [16]byte
	_, _ = rng.Read(b[:]) //nolint:errcheck // does not return errors
	// version, variant
	b[6], b[8] = (b[6]&0x0f)|0x40, (b[8]&0x3f)|0x80 //nolint:mnd // lob

	return UUID{b}
}
