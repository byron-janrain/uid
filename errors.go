package uid

// ParseError is a sentinel error.
type ParseError struct{}

// Error implements errors.Error.
func (e ParseError) Error() string { return "" }
