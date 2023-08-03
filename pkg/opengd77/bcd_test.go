package opengd77

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromBCD(t *testing.T) {
	val := 0x14543000
	want := 14543000
	have := FromBCD(val)

	assert.Equal(t, want, have, "converting from bcd failed")
}

func TestToBCD(t *testing.T) {
	val := 14543000
	want := 0x14543000
	have := ToBCD(val)

	assert.Equal(t, want, have, "converting to bcd failed")
}
