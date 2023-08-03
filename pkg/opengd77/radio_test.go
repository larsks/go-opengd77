package opengd77

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.bug.st/serial"
)

func TestNewRadio(t *testing.T) {
	radio, err := NewRadio("/dev/ttyS0")
	assert.Nil(t, err)
	assert.Equal(t, "/dev/ttyS0", radio.serialDevice)
}

func TestNewRadioWithTimeout(t *testing.T) {
	radio, err := NewRadio("/dev/ttyS0", WithReadTimeout(10*time.Second))
	assert.Nil(t, err)
	assert.Equal(t, 10*time.Second, radio.serialTimeout)
}

func TestNewRadioWithSerialMode(t *testing.T) {
	radio, err := NewRadio("/dev/ttyS0",
		WithSerialMode(&serial.Mode{BaudRate: 115200}))
	assert.Nil(t, err)
	assert.Equal(t, 115200, radio.serialMode.BaudRate)
}

func TestSendSimpleCommand(t *testing.T) {
	buf := NewTestBuffer([]byte{'-'})
	radio, err := NewRadio("/dev/ttyS0", WithSerialPort(buf))
	assert.Nil(t, err)
	err = radio.sendSimpleCommand([]byte{1, 2, 3, 4})
	assert.Nil(t, err)
	assert.Equal(t, buf.WriteBuffer.Bytes(), []byte{'C', 1, 2, 3, 4})
}

func TestSendSimpleCommandFails(t *testing.T) {
	buf := NewTestBuffer([]byte{})
	radio, err := NewRadio("/dev/ttyS0", WithSerialPort(buf))
	assert.Nil(t, err)
	err = radio.sendSimpleCommand([]byte{1, 2, 3, 4})
	assert.ErrorIs(t, err, io.EOF)
}

func TestSendExtendedCommand(t *testing.T) {
	buf := NewTestBuffer([]byte{'-'})
	radio, err := NewRadio("/dev/ttyS0", WithSerialPort(buf))
	assert.Nil(t, err)
	err = radio.sendExtendedCommand(1)
	assert.Nil(t, err)
	assert.Equal(t, []byte{'C', 6, 1}, buf.WriteBuffer.Bytes())
}

func TestReadChunk(t *testing.T) {
	buf := NewTestBuffer([]byte{'R', 0x00, 0x04, 1, 2, 3, 4})
	radio, err := NewRadio("/dev/ttyS0", WithSerialPort(buf))
	assert.Nil(t, err)
	data := make([]byte, 4)
	nb, err := radio.readChunk(1, 0, 4, data)
	assert.Nil(t, err)
	assert.Equal(t, nb, 4)
	assert.Equal(t, data, []byte{1, 2, 3, 4})
}
