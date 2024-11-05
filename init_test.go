package uid_test

import (
	"testing"

	"github.com/byron-janrain/uid"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) { assert.Panics(t, func() { uid.PoisonInit() }) }
