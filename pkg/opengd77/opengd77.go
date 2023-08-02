package opengd77

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

const (
	MAX_TRANSFER_SIZE = 32
	FLASH_BLOCK_SIZE  = 4096

	MEMTYPE_FLASH         = 1
	MEMTYPE_EEPROM        = 2
	MEMTYPE_MCU_ROM       = 5
	MEMTYPE_DISPLAY_BUFFE = 6
	MEMTYPE_SOUND_BUFFER  = 7
	MEMTYPE_VOICE_BUFFER  = 8
)

type (
	Radio struct {
		port serial.Port
	}

	RadioOption func(*Radio) error
)

var (
	ConfigBlocks = []MemoryBlock{
		{MEMTYPE_EEPROM, 0x00e0, 0x00e0, 0x6000 - 0xe0},
		{MEMTYPE_EEPROM, 0x7500, 0x7500, 0xb000 - 0x7500},
		{MEMTYPE_FLASH, 0xb000, 0x7b000, 0x13e60},
		{MEMTYPE_FLASH, 0x1ee60, 0x0000, 0x11a0},
	}

	ChannelBlocks = []int{
		0x3780,
		0x0b1b0,
		0x0cdc0,
		0x0e9d0,
		0x105e0,
		0x121f0,
		0x13e00,
		0x15a10,
	}
)

func WithReadTimeout(timeout time.Duration) RadioOption {
	return func(radio *Radio) error {
		return radio.port.SetReadTimeout(timeout)
	}
}

func NewRadio(serialport string, options ...RadioOption) (*Radio, error) {
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(serialport, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", serialport, err)
	}
	if err := port.SetReadTimeout(3 * time.Second); err != nil {
		return nil, err
	}

	radio := &Radio{
		port: port,
	}

	for _, opt := range options {
		if err := opt(radio); err != nil {
			return nil, err
		}
	}

	return radio, nil
}

func (radio *Radio) Close() error {
	return radio.port.Close()
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
