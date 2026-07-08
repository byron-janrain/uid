package uid

// ErrInvalid is the constant sentinel returned by all unserializers for malformed input.
const ErrInvalid = anError("uid: invalid")

type anError string

// Error implements error.
func (e anError) Error() string { return string(e) }
