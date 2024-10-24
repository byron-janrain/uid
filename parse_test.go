package uid_test

import (
	"strings"
	"testing"

	"github.com/byron-janrain/uid"
	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	id := uid.NewV4()
	assert.Exactly(t, id.String(), uid.Must(uid.Parse(id.String())).String())
}

func TestReferences(t *testing.T) {
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
