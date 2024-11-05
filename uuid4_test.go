package uid_test

import (
	"encoding/json"
	"testing"

	"github.com/byron-janrain/uid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestV4Equal(t *testing.T) {
	reset := uid.ReseedPRNG()
	id1 := uid.NewV4()
	reset()
	defer uid.ReseedPRNG()()
	id2 := uid.NewV4()
	assert.Exactly(t, id1, id2)
}

func TestV4(t *testing.T) {
	defer uid.ReseedPRNG()()
	id := uid.NewV4()
	// identity
	assert.Exactly(t, uid.Version4, id.Version())
	assert.False(t, id.IsMax())
	assert.False(t, id.IsNil())
	assert.True(t, id.Time().IsZero())
	// string
	assert.Exactly(t, "d9877ece-6d36-4aac-9a6f-419ec627c76b", id.String())
	// bytes
	expected := []byte{0xd9, 0x87, 0x7e, 0xce, 0x6d, 0x36, 0x4a, 0xac, 0x9a, 0x6f, 0x41, 0x9e, 0xc6, 0x27, 0xc7, 0x6b}
	assert.Exactly(t, expected, id.Bytes())
	// unmarshals
	s := id.String()
	id2 := uid.UUID{}
	require.NoError(t, id2.UnmarshalBinary(id.Bytes()))
	assert.Exactly(t, id, id2)
	require.NoError(t, id2.UnmarshalText([]byte(s)))
	assert.Exactly(t, id, id2)
	jsb, err := json.Marshal(s)
	require.NoError(t, err)
	require.NoError(t, id2.UnmarshalJSON(jsb))
	assert.Exactly(t, id, id2)
}
