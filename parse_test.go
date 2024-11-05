package uid_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/byron-janrain/uid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBads(t *testing.T) {
	check := func(in, expectedErr string, isWrapped bool) {
		id, err := uid.Parse(in)
		assert.Panics(t, func() {
			_ = uid.Must(uid.Parse(in))
		})
		require.ErrorAs(t, err, &uid.ParseError{})
		require.EqualError(t, err, expectedErr)
		if isWrapped {
			require.Error(t, errors.Unwrap(err))
		} else {
			require.NoError(t, errors.Unwrap(err))
		}
		assert.Exactly(t, uid.Nil(), id)
	}
	defaultErrTxt := uid.ParseError{}.Error()
	// bad length
	check("", uid.ParseError{}.Error(), false)
	// bad canonical nil
	check("00000000-0000-0000-0000-000000000001", defaultErrTxt, false)
	// bad canonical max
	check("ffffffff-ffff-ffff-ffff-fffffffffffe", defaultErrTxt, false)
	// bad canonical encoding/format
	check(strings.Replace(ref4, "-", "_", 1), "failed to parse UUID: encoding/hex: invalid byte: U+005F '_'", true)
	// bad canonical version
	check(strings.ReplaceAll(ref4, "4", "e"), defaultErrTxt, false)
	// bad canonical v4 variant
	check("ffffffff-ffff-4fff-efff-ffffffffffff", defaultErrTxt, false)
	// bad compact32 encoding
	check("E077777777777777777777777P", "failed to parse UUID: illegal base32 data at input byte 0", true)
	// bad compact32 version
	check("Q777777777777777777777777P", defaultErrTxt, false)
	// bad compact32 nil
	check("ABAAAAAAAAAAAAAAAAAAAAAAAA", defaultErrTxt, false)
	// bad compact32 version
	check("P777777777777777777777777Q", defaultErrTxt, false)
	// bad compact64 encoding
	check("E+___________________P", "failed to parse UUID: illegal base64 data at input byte 0", true)
	// bad compact64 version
	check("Q____________________P", defaultErrTxt, false)
	// bad compact64 variant
	check("P____________________Q", defaultErrTxt, false)
	// bad compact64 nil
	check("ABAAAAAAAAAAAAAAAAAAAA", defaultErrTxt, false)
}

func TestSamples(t *testing.T) {
	tested := 0
	for _, sample := range sampleData {
		ids := [3]uid.UUID{
			uid.Must(uid.Parse(sample.Canonical)),
			uid.Must(uid.Parse(sample.B32)),
			uid.Must(uid.Parse(sample.B64)),
		}
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
		actual := uid.Must(uid.Parse(input))
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
