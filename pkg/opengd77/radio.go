package opengd77

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"

	"go.bug.st/serial"
)

const (
	// Max number of bytes can can be transferred using readChunk
	MAX_TRANSFER_SIZE = 32

	MEMTYPE_FLASH         = 1
	MEMTYPE_EEPROM        = 2
	MEMTYPE_MCU_ROM       = 5
	MEMTYPE_DISPLAY_BUFFE = 6
	MEMTYPE_SOUND_BUFFER  = 7
	MEMTYPE_VOICE_BUFFER  = 8
)

type (
	// A MemoryBlock maps between a block of memory in the radio
	// and the corresponding location in the dump file
	MemoryBlock struct {
		Kind        int
		FileOffset  int
		RadioOffset int
		Length      int
	}

	// Radio represents a radio running OpenGD77
	Radio struct {
		serialDevice  string
		serialMode    *serial.Mode
		serialTimeout time.Duration
		port          io.ReadWriteCloser
	}

	// Represents an option that can be passed to the Radio constructor
	RadioOption func(*Radio) error
)

var (
	// Location of configuration data in the radio
	ConfigBlocks = []MemoryBlock{
		{MEMTYPE_EEPROM, 0x00e0, 0x00e0, 0x6000 - 0xe0},
		{MEMTYPE_EEPROM, 0x7500, 0x7500, 0xb000 - 0x7500},
		{MEMTYPE_FLASH, 0xb000, 0x7b000, 0x13e60},
		{MEMTYPE_FLASH, 0x1ee60, 0x0000, 0x11a0},
	}
)

// Use a custom object as the serial port. This can be used
// to pass in an already open serial.Port object or to pass in
// any other io.ReadWriteCloser for testing.
func WithSerialPort(port io.ReadWriteCloser) RadioOption {
	return func(radio *Radio) error {
		radio.port = port
		return nil
	}
}

// Set a custom read timeout on the serial port.
func WithReadTimeout(timeout time.Duration) RadioOption {
	return func(radio *Radio) error {
		radio.serialTimeout = timeout
		return nil
	}
}

// Set a custom serial port configuration (data rate,
// parity, stop bits, etc)
func WithSerialMode(mode *serial.Mode) RadioOption {
	return func(radio *Radio) error {
		radio.serialMode = mode
		return nil
	}
}

// Create and return a new Radio object
func NewRadio(serialport string, options ...RadioOption) (*Radio, error) {
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	radio := &Radio{
		serialDevice:  serialport,
		serialMode:    mode,
		serialTimeout: 3 * time.Second,
	}

	for _, opt := range options {
		if err := opt(radio); err != nil {
			return nil, err
		}
	}

	return radio, nil
}

// Open the serial port and set the port attribute on the Radio object.
// This will not be necessary if you've passed in a custom port
// using WithSerialPort.
func (radio *Radio) Open() error {
	port, err := serial.Open(radio.serialDevice, radio.serialMode)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", radio.serialDevice, err)
	}
	if err := port.SetReadTimeout(radio.serialTimeout); err != nil {
		return fmt.Errorf("failed to set serial port timeout: %w", err)
	}
	radio.port = port

	return nil
}

func (radio *Radio) Close() error {
	if radio.port != nil {
		return radio.port.Close()
	}

	return nil
}

func (radio *Radio) sendSimpleCommand(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("invalid null command")
	}

	cmd := []byte{'C'}
	cmd = append(cmd, data...)

	nb, err := radio.port.Write(cmd)
	if err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}
	if nb != len(cmd) {
		return fmt.Errorf("failed to send command: short write")
	}

	buf := make([]byte, 1)
	nb, err = radio.port.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	if nb != 1 || buf[0] != byte('-') {
		return fmt.Errorf("unexpected response from radio")
	}

	return nil
}

func (radio *Radio) sendExtendedCommand(subcommand byte) error {
	return radio.sendSimpleCommand([]byte{6, subcommand})
}

func (radio *Radio) ScreenShow() error {
	return radio.sendSimpleCommand([]byte{0})
}

func (radio *Radio) ScreenClear() error {
	return radio.sendSimpleCommand([]byte{1})
}

func (radio *Radio) ScreenPrint(x, y, size, align, invert byte, message string) error {
	data := []byte{2, x, y, size, align, invert}
	data = append(data, []byte(message)...)
	return radio.sendSimpleCommand(data)
}

func (radio *Radio) ScreenRender() error {
	return radio.sendSimpleCommand([]byte{3})
}

// Convenience function for showing a one-line message
func (radio *Radio) ShowMessage(size int, message string) error {
	var err error

	if err = radio.ScreenClear(); err == nil {
		if err = radio.ScreenShow(); err == nil {
			if err = radio.ScreenPrint(0, 16, byte(size), 1, 0, message); err == nil {
				err = radio.ScreenRender()
			}
		}
	}

	return err
}

func (radio *Radio) ScreenBacklight() error {
	return radio.sendSimpleCommand([]byte{4})
}

func (radio *Radio) ScreenClose() error {
	return radio.sendSimpleCommand([]byte{5})
}

func (radio *Radio) SaveAndReboot() error {
	return radio.sendExtendedCommand(0)
}

func (radio *Radio) Reboot() error {
	return radio.sendExtendedCommand(1)
}

func (radio *Radio) Save() error {
	return radio.sendExtendedCommand(2)
}

func (radio *Radio) FlashGreen() error {
	return radio.sendExtendedCommand(3)
}

func (radio *Radio) FlashRed() error {
	return radio.sendExtendedCommand(4)
}

func (radio *Radio) CodecInit() error {
	return radio.sendExtendedCommand(5)
}

func (radio *Radio) SoundInit() error {
	return radio.sendExtendedCommand(6)
}

func (radio *Radio) readChunk(mode, addr, length int, data []byte) (int, error) {
	if length > MAX_TRANSFER_SIZE {
		length = MAX_TRANSFER_SIZE
	}

	req := []byte{'R', byte(mode), 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint32(req[2:], uint32(addr))
	binary.BigEndian.PutUint16(req[6:], uint16(length))

	nb, err := radio.port.Write(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send command: %w", err)
	}
	if nb != len(req) {
		return 0, fmt.Errorf("failed to send command: short write")
	}

	res := make([]byte, 3)
	nb, err = radio.port.Read(res)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}
	if nb != len(res) {
		return 0, fmt.Errorf("failed to read response: short read")
	}

	reslen := binary.BigEndian.Uint16(res[1:])
	if int(reslen) != length {
		log.Printf("short read: requested %d bytes, offered %d", length, reslen)
	}

	nb, err = radio.port.Read(data)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}
	if nb != int(reslen) {
		log.Printf("short read: expected %d bytes, received %d", reslen, nb)
	}

	return nb, nil
}

func (radio *Radio) ReadMemory(mode, addr, length int, data []byte) error {
	remaining := length
	offset := 0
	for remaining > 0 {
		nb, err := radio.readChunk(mode, addr+offset, remaining, data[offset:])
		if err != nil {
			return fmt.Errorf("failed to read memory @ %d: %w", addr+offset, err)
		}

		remaining -= nb
		offset += nb
	}

	return nil
}

func (radio *Radio) ReadCodePlug(w io.WriteSeeker) error {
	for _, block := range ConfigBlocks {
		log.Printf("reading %d bytes from %x", block.Length, block.RadioOffset)
		buf := make([]byte, block.Length)
		if err := radio.ReadMemory(block.Kind, block.RadioOffset, block.Length, buf); err != nil {
			return fmt.Errorf("failed to read from radio: %w", err)
		}

		if _, err := w.Seek(int64(block.FileOffset), 0); err != nil {
			return fmt.Errorf("failed to set output position: %w", err)
		}

		if _, err := w.Write(buf); err != nil {
			return fmt.Errorf("failed to write data: %w", err)
		}
	}
	return nil
}
