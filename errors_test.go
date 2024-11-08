package uid_test

import (
	"testing"

	"github.com/byron-janrain/uid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	require.Error(t, uid.ParseError{})
	assert.Empty(t, uid.ParseError{}.Error()) // yup it's on purpose
}
