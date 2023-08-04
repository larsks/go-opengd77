package opengd77

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackedTimeslot(t *testing.T) {
	want := PackedTimeslotOverride(0x02)
	unpacked := want.Unpack()
	assert.True(t, unpacked.HasOverride)
	assert.Equal(t, unpacked.Timeslot, 1)

	packed := unpacked.Pack()
	assert.Equal(t, want, packed)
}
