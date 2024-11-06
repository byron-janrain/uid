package uid_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/byron-janrain/uid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ref4    = "01867b2c-a0dd-459c-98d7-89e545538d6c"
	ref4b32 = "EAGDHWLFA3VM4RV4J4VCVHDLMJ"
	ref4b64 = "EAYZ7LKDdWcjXieVFU41sJ"
	ref7    = "0191e843-b452-7ac4-b853-8ee3953a28af"
	ref7b32 = "HAGI6QQ5UKKWEQU4O4OKTUKFPL"
	ref7b64 = "HAZHoQ7RSrEhTjuOVOiivL"
)

func TestCommonAccessors(t *testing.T) {
	// v4 max is a good surrogate since none of it's values are zero AND unshifting has an edge case
	id := uid.Max()
	// check methods
	assert.Exactly(t, uid.VersionMax, id.Version()) // 0xf
	assert.Exactly(t, uint8(3), id.Variant())       // 0x3
	// check marshal-unmarshal equivalences
	var id2 uid.UUID
	// text
	txt, err := id.MarshalText()
	require.NoError(t, err)
	assert.Exactly(t, id.String(), string(txt))
	require.NoError(t, id2.UnmarshalText(txt))
	assert.Exactly(t, id, id2)
	// json
	data, err := id.MarshalJSON()
	require.NoError(t, err)
	require.NoError(t, id2.UnmarshalJSON(data))
	assert.Exactly(t, id, id2)
	// binary
	data, err = id.MarshalBinary()
	require.NoError(t, err)
	assert.Exactly(t, id.Bytes(), data)
	require.NoError(t, id2.UnmarshalBinary(data))
	assert.Exactly(t, id, id2)
	// check compact forms in text unmarshaling
	for _, txt := range []string{id.Compact32(), id.Compact64()} {
		var id2 uid.UUID
		// txt
		err := id2.UnmarshalText([]byte(txt))
		require.NoError(t, err)
		assert.Exactly(t, id, id2)
		// json
		data, err := json.Marshal(txt)
		require.NoError(t, err)
		err = json.Unmarshal(data, &id2)
		require.NoError(t, err)
		assert.Exactly(t, id, id2)
	}
}

func TestBytesImmutable(t *testing.T) {
	id := uid.Max()   // use max for non-zero values
	id.Bytes()[0] = 0 // should be copy
	assert.Exactly(t, id.Bytes(), uid.Max().Bytes())
}

func TestUnmarshalBinaryFail(t *testing.T) {
	var id uid.UUID
	err := id.UnmarshalBinary([]byte{})
	require.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "failed to unmarshal binary: "))
}

func TestUnmarshalTextFail(t *testing.T) {
	var id uid.UUID
	err := id.UnmarshalText([]byte{})
	require.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "failed to unmarshal text: "))
}

func TestUnmarshalJSONFail(t *testing.T) {
	var id uid.UUID
	err := id.UnmarshalJSON([]byte{})
	require.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "failed to unmarshal json: "))
}

func TestNil(t *testing.T) {
	assert.Exactly(t, uid.VersionNil, uid.Nil().Version()) // version
	assert.True(t, uid.Nil().IsNil())                      // identity
	// binary equivalence
	assert.Exactly(t, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, uid.Nil().Bytes())
	// negatives
	assert.False(t, uid.Nil().IsMax())
	assert.True(t, uid.Nil().Time().IsZero())
	// stringifications
	assert.Exactly(t, uid.NilCanonical, uid.Nil().String())
	assert.Exactly(t, uid.NilCompact32, uid.Nil().Compact32())
	assert.Exactly(t, uid.NilCompact64, uid.Nil().Compact64())
	actual, err := uid.Nil().MarshalJSON()
	require.NoError(t, err)
	assert.Exactly(t, uid.NilJSON, string(actual))
}

func TestMax(t *testing.T) {
	assert.Exactly(t, uid.VersionMax, uid.Max().Version()) // version
	// identity
	assert.True(t, uid.Max().IsMax())
	assert.False(t, uid.Max().IsNil())
	assert.True(t, uid.Max().Time().IsZero())
	// binary equivalence
	expected := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	assert.Exactly(t, expected, uid.Max().Bytes())
	// stringifications
	assert.Exactly(t, uid.MaxCanonical, uid.Max().String())
	assert.Exactly(t, uid.MaxCompact32, uid.Max().Compact32())
	assert.Exactly(t, uid.MaxCompact64, uid.Max().Compact64())
	actual, err := uid.Max().MarshalJSON()
	require.NoError(t, err)
	assert.Exactly(t, uid.MaxJSON, string(actual))
}

func TestCompare(t *testing.T) {
	assert.Exactly(t, 0, uid.Compare(uid.Nil(), uid.Nil()))
	assert.Exactly(t, 0, uid.Compare(uid.Max(), uid.Max()))
	id1 := uid.NewV7()
	time.Sleep(time.Microsecond)
	id2 := uid.NewV7()
	assert.Exactly(t, -1, uid.Compare(id1, id2))
}

func TestSanityCollisions(t *testing.T) {
	const count = 2 * 1000 * 1000 // fill 2,000,000 ids to ensure we don't have a glaring issue
	ids := map[uid.UUID]bool{}    // doubles as check that uid.UUID can be used as map key
	// v4
	for range count {
		id := uid.NewV4()
		_, ok := ids[id]
		assert.False(t, ok)
		ids[id] = true
	}
	// v7 in the same set
	for range count {
		id := uid.NewV7()
		_, ok := ids[id]
		assert.False(t, ok)
		ids[id] = true
	}
}
