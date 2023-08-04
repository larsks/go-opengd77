package opengd77

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaddedName(t *testing.T) {
	// https://tip.golang.org/ref/spec#Conversions_from_slice_to_array_or_array_pointer
	var val PaddedName = [16]byte(bytes.Repeat([]byte{0xff}, 16))
	copy(val[:], []byte("TEST"))
	have := val.String()
	want := "TEST"

	assert.Len(t, val, 16)
	assert.Len(t, have, 4)
	assert.Equal(t, want, have)
}

func TestBCDFrequencyString(t *testing.T) {
	val := BCDFrequency(ToBCD(14543000))

	assert.Equal(t, "145.4300", val.String())
}

func TestBCDFrequencyFloat(t *testing.T) {
	val := BCDFrequency(ToBCD(14543000))

	assert.Equal(t, 145.43, val.Float())
}

func TestNewChannel(t *testing.T) {
	ch := NewChannel()
	var want PaddedName = [16]byte(bytes.Repeat([]byte{0xff}, 16))
	assert.Equal(t, want, ch.Name)
}

func TestSetName(t *testing.T) {
	ch := NewChannel()
	ch.SetName("TEST")
	var want PaddedName = [16]byte(bytes.Repeat([]byte{0xff}, 16))
	copy(want[:], []byte("TEST"))

	assert.Equal(t, want, ch.Name)
}

func TestSetRxFrequency(t *testing.T) {
	ch := NewChannel()
	ch.SetRxFreq(145.43)

	assert.Equal(t, 145.43, ch.GetRxFreq())
}

func TestSetTxFrequency(t *testing.T) {
	ch := NewChannel()
	ch.SetTxFreq(145.43)

	assert.Equal(t, 145.43, ch.GetTxFreq())
}

func TestSetRxTone(t *testing.T) {
	ch := NewChannel()
	ch.SetRxTone(88.5)
	assert.Equal(t, 88.5, ch.GetRxTone())
}

func TestSetTxTone(t *testing.T) {
	ch := NewChannel()
	ch.SetTxTone(88.5)
	assert.Equal(t, 88.5, ch.GetTxTone())
}

func TestChannelFlag4(t *testing.T) {
	want := PackedChannelFlag4(CODEPLUG_CHANNEL_FLAG4_RX_ONLY | CODEPLUG_CHANNEL_FLAG4_ALL_SKIP)
	unpacked := want.Unpack()
	assert.True(t, unpacked.RxOnly)
	assert.True(t, unpacked.AllSkip)
	assert.False(t, unpacked.Vox)
	assert.False(t, unpacked.BW25K)
	assert.False(t, unpacked.Power)
	assert.False(t, unpacked.Squelch)
	assert.False(t, unpacked.ZoneSkip)

	packed := unpacked.Pack()
	assert.Equal(t, want, packed)
}

func TestLibreDMRFlag1(t *testing.T) {
	want := PackedLibreDMRFlag1(byte(CODEPLUG_CHANNEL_LIBREDMR_FLAG1_OPTIONAL_DMRID | CODEPLUG_CHANNEL_LIBREDMR_FLAG1_NO_ECO))
	unpacked := want.Unpack()
	assert.True(t, unpacked.DMRId)
	assert.True(t, unpacked.NoEco)
	assert.False(t, unpacked.NoBeep)
	assert.False(t, unpacked.OutOfBand)

	packed := unpacked.Pack()
	assert.Equal(t, want, packed)
}
