package uid

// ParseError is returned when parsing a UUID fails.
type ParseError struct{ error }

// Error implements errors.Error.
func (e ParseError) Error() string {
	const msg = "failed to parse UUID: "
	if e.error != nil {
		return msg + e.error.Error()
	}
	return msg + "malformed or unsupported"
}

// Unwrap implements errors.Unwrap.
func (e ParseError) Unwrap() error { return e.error }
