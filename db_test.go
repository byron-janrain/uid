package uid_test

import (
	"testing"

	"github.com/hoodie-ninja/uid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValue(t *testing.T) {
	id := uid.NewV4()
	v, err := id.Value()
	require.NoError(t, err)
	assert.Equal(t, id.String(), v)
}

func TestScan(t *testing.T) {
	id := uid.NewV4()
	// canonical string, canonical []byte, and raw 16 bytes all round-trip
	for _, src := range []any{id.String(), []byte(id.String()), id.Bytes()} {
		var got uid.UUID
		require.NoError(t, got.Scan(src))
		assert.Exactly(t, id, got)
	}
	// SQL NULL >> error
	var n uid.UUID
	require.Error(t, n.Scan(nil))
	// unsupported type and malformed text >> ErrInvalid
	var bad uid.UUID
	require.ErrorIs(t, bad.Scan(true), uid.ErrInvalid)
	require.ErrorIs(t, bad.Scan("not-a-uuid"), uid.ErrInvalid)
}
