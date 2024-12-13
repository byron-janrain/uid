package uid_test

import (
	"testing"

	"github.com/byron-janrain/uid"
	"github.com/stretchr/testify/assert"
)

func TestToPythonShort(t *testing.T) {
	check := func(u uid.UUID, expected string) { assert.Exactly(t, expected, uid.ToPythonShort(u)) }
	check(uid.Max(), uid.MaxPythonShort)
	check(uid.Nil(), uid.NilPythonShort)
	sut, ok := uid.Parse("3b1f8b40-222c-4a6e-b77e-779d5a94e21c")
	assert.True(t, ok)
	check(sut, "CXc85b4rqinB7s5J52TRYb")
}

func TestFromPythonShortHappy(t *testing.T) {
	check := func(input, expected string) {
		sut, ok := uid.FromPythonShort(input)
		assert.True(t, ok)
		actual := sut.String()
		assert.Exactly(t, expected, actual)
	}
	check(uid.MaxPythonShort, uid.MaxCanonical)
	check(uid.NilPythonShort, uid.NilCanonical)
	check("CXc85b4rqinB7s5J52TRYb", "3b1f8b40-222c-4a6e-b77e-779d5a94e21c")
	check(" CXc85b4rqinB7s5J52TRYb\t", "3b1f8b40-222c-4a6e-b77e-779d5a94e21c")
}

func TestFromPythonShortBads(t *testing.T) {
	shouldFail := func(badInput string) {
		u, ok := uid.FromPythonShort(badInput)
		assert.Exactly(t, uid.Nil(), u)
		assert.False(t, ok)
	}
	shouldFail("")                                           // empty
	shouldFail("tooshort")                                   // too short
	shouldFail("thisinputislongerthan22runes")               // too long
	shouldFail("02222" + "22222" + "22222" + "22222" + "22") // right length, bad runes
}
