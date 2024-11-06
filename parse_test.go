package uid_test

import (
	"strings"
	"testing"

	"github.com/byron-janrain/uid"
	"github.com/stretchr/testify/assert"
)

func TestParseBads(t *testing.T) {
	check := func(in string) {
		id, ok := uid.Parse(in)
		assert.False(t, ok)
		assert.Exactly(t, uid.Nil(), id)
	}
	// bad canonical nil
	check("00000000-0000-0000-0000-000000000001")
	// bad canonical max
	check("ffffffff-ffff-ffff-ffff-fffffffffffe")
	// bad canonical format
	check(strings.Replace(ref4, "-", "_", 1))
	// bad canonical encoding
	check("00000000-0000-4000-8000-g00000000000")
	// bad canonical version
	check(strings.ReplaceAll(ref4, "4", "e"))
	// bad canonical v4 variant
	check("ffffffff-ffff-4fff-efff-ffffffffffff")
	// bad compact32 encoding
	check("E077777777777777777777777P")
	// bad compact32 version
	check("Q777777777777777777777777P")
	// bad compact32 nil
	check("ABAAAAAAAAAAAAAAAAAAAAAAAA")
	// bad compact32 version
	check("P777777777777777777777777Q")
	// bad compact64 encoding
	check("E+___________________P")
	// bad compact64 version
	check("Q____________________P")
	// bad compact64 variant
	check("P____________________Q")
	// bad compact64 nil
	check("ABAAAAAAAAAAAAAAAAAAAA")
}

func TestSamples(t *testing.T) {
	tested := 0
	for _, sample := range sampleData {
		c, ok := uid.Parse(sample.Canonical)
		assert.True(t, ok)
		b32, ok := uid.Parse(sample.B32)
		assert.True(t, ok)
		b64, ok := uid.Parse(sample.B64)
		assert.True(t, ok)
		ids := [3]uid.UUID{c, b32, b64}
		// ensure equivalent parsing
		assert.Exactly(t, ids[0], ids[1])
		assert.Exactly(t, ids[1], ids[2])
		for _, id := range ids {
			assert.Exactly(t, uid.Version4, id.Version())
			assert.Exactly(t, uid.Variant9562, id.Variant())
			assert.Exactly(t, sample.Canonical, id.String())
			assert.True(t, strings.EqualFold(sample.B32, id.Compact32()))
			assert.Exactly(t, sample.B64, id.Compact64())
		}
		tested++
	}
	assert.Exactly(t, 1000, tested)
}

func TestMagicIDs(t *testing.T) {
	check := func(input string, expected uid.UUID) {
		actual, ok := uid.Parse(input)
		assert.True(t, ok)
		assert.Exactly(t, expected, actual)
	}
	// nil
	check(uid.NilCanonical, uid.Nil())
	check(uid.NilCompact32, uid.Nil())
	check(uid.NilCompact64, uid.Nil())
	// max
	check(uid.MaxCanonical, uid.Max())
	check(uid.MaxCompact32, uid.Max())
	check(uid.MaxCompact64, uid.Max())
}
