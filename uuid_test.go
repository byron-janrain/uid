package uid_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/byron-janrain/uid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	checkStrings := func(id uid.UUID, expectedCanonical, expectedCompact32, expectedCompact64 string,
		isNil, isMax, isTimeZero bool) {
		assert.Exactly(t, expectedCanonical, id.String())
		actualTxt, err := id.MarshalText()
		require.NoError(t, err)
		assert.Exactly(t, []byte(expectedCanonical), actualTxt)
		actualJS, err := id.MarshalJSON()
		require.NoError(t, err)
		expectedJS, err := json.Marshal(expectedCanonical)
		require.NoError(t, err)
		assert.Exactly(t, expectedJS, actualJS)
		require.NoError(t, err)
		assert.Exactly(t, expectedCompact32, id.Compact32())
		assert.Exactly(t, expectedCompact64, id.Compact64())
		assert.Exactly(t, isMax, id.IsMax())
		assert.Exactly(t, isNil, id.IsNil())
		assert.Exactly(t, isTimeZero, id.Time().IsZero())
	}
	// check static UUIDs (nil and max)
	checkStrings(uid.Nil(), uid.NilCanonical, uid.NilCompact32, uid.NilCompact64, true, false, true)
	checkStrings(uid.Max(), uid.MaxCanonical, uid.MaxCompact32, uid.MaxCompact64, false, true, true)
	// check nsname compact reference
	const ref4 = "01867b2c-a0dd-459c-98d7-89e545538d6c"
	refid := uid.Must(uid.Parse(ref4))
	checkStrings(refid, ref4, "EAGDHWLFA3VM4RV4J4VCVHDLMJ", "EAYZ7LKDdWcjXieVFU41sJ", false, false, true)
	const ref7 = "0191e843-b452-7ac4-b853-8ee3953a28af"
	refid = uid.Must(uid.Parse(ref7))
	checkStrings(refid, ref7, "HAGI6QQ5UKKWEQU4O4OKTUKFPL", "HAZHoQ7RSrEhTjuOVOiivL", false, false, false)
}

func TestSanityCheckCollision(t *testing.T) {
	const max = 1000000            // fill 1,000,000 ids to ensure we don't have a glaring dup/underflow issue or v4
	ids := map[uid.UUID]struct{}{} // doubles as check that uid.UUID can be used as map key
	// v4
	for i := 0; i < max; i++ {
		id := uid.NewV4()
		_, ok := ids[id]
		assert.False(t, ok)
		ids[id] = struct{}{}
	}
	// v7 in the same set
	for i := 0; i < max; i++ {
		id := uid.NewV7()
		_, ok := ids[id]
		assert.False(t, ok)
		ids[id] = struct{}{}
	}
}

func TestEqual(t *testing.T) {
	uid.ReseedPRNG()
	id1 := uid.NewV4()
	uid.ReseedPRNG()
	id2 := uid.NewV4()
	assert.Exactly(t, id1, id2)
}

func TestBytesImmutable(t *testing.T) {
	id1 := uid.NewV4()
	var id2 uid.UUID
	err := id2.UnmarshalBinary(id1.Bytes())
	require.NoError(t, err)
	assert.Exactly(t, id1, id2)
	assert.Exactly(t, id1.Bytes(), id2.Bytes())
	// ensure bytes copy returned
	b := id1.Bytes()
	b[6] = 0x0 // byte 6 has version so shouldn't be
	assert.NotEqual(t, 0x0, id1.Bytes()[6])
}

func TestMax(t *testing.T) {
	id := uid.Max()
	// version
	assert.Exactly(t, uid.VersionMax, id.Version())
	assert.True(t, id.IsMax())
	// binary equivalence
	assert.Exactly(t, []byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, id.Bytes())
	// negative tests
	assert.False(t, id.IsNil())
	assert.True(t, id.Time().IsZero())
}

func TestNil(t *testing.T) {
	id := uid.Nil()
	// version
	assert.Exactly(t, uid.VersionNil, id.Version())
	// identity
	assert.False(t, id.IsMax())
	assert.True(t, id.IsNil())
	// binary equivalence
	assert.Exactly(t, []byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, id.Bytes())
	// time
	assert.True(t, id.Time().IsZero())
	// stringifications
	assert.Exactly(t, uid.NilCanonical, id.String())
	assert.Exactly(t, uid.NilCompact32, uid.Nil().Compact32())
	assert.Exactly(t, uid.NilCompact64, uid.Nil().Compact64())
}

// func TestV4(t *testing.T) {
// 	uid.ReseedPRNG()
// 	id := uid.NewV4()
// 	// identity
// 	assert.Exactly(t, uid.Version4, id.Version())
// 	assert.False(t, id.IsMax())
// 	assert.False(t, id.IsNil())
// 	assert.True(t, id.Time().IsZero())
// 	// string
// 	assert.Exactly(t, "d9877ece-6d36-4aac-9a6f-419ec627c76b", id.String())
// 	// bytes
// 	assert.Exactly(t, []byte{
// 		0xd9, 0x87, 0x7e, 0xce, 0x6d, 0x36, 0x4a, 0xac, 0x9a, 0x6f, 0x41, 0x9e, 0xc6, 0x27, 0xc7, 0x6b}, id.Bytes())
// }

// func TestV7(t *testing.T) {
// 	uid.ReseedPRNG()
// 	id := uid.NewV7()
// 	// identity
// 	assert.Exactly(t, uid.Version7, id.Version())
// 	assert.False(t, id.IsMax())
// 	assert.False(t, id.IsNil())
// 	assert.False(t, id.Time().IsZero())
// 	// check version, seeded rng, and variant string, msts tested later
// 	assert.True(t, strings.HasSuffix(id.String(), "-7987-bece-6d368aac1a6f"))
// 	// check version, seeded rng, and variant bytes, msts tested later
// 	assert.Exactly(t, []byte{0x79, 0x87, 0xbe, 0xce, 0x6d, 0x36, 0x8a, 0xac, 0x1a, 0x6f}, id.Bytes()[6:])
// 	// check time is passing
// 	time.Sleep(100) //nolint:staticcheck //intended
// 	assert.True(t, time.Now().After(id.Time()))
// 	// sanity check time is near "now" (within 1s)
// 	assert.InDelta(t, time.Now().UnixMilli(), id.Time().UnixMilli(), 1000)
// }

func TestPreEpochPanic(t *testing.T) {
	// check that negative times panic (impossible outside this unit test)
	defer uid.SetNowFunc(func() time.Time { return time.Unix(-10, 0) })()
	assert.Panics(t, func() {
		_ = uid.NewV7()
	})
}

func TestEpochEdgeCase(t *testing.T) {
	defer uid.SetNowFunc(func() time.Time { return time.Unix(0, 0) })()
	assert.Exactly(t, int64(0), uid.NewV7().Time().UnixMilli())
}
