package opengd77

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

const (
	CODEPLUG_ADDR_USER_CALLSIGN    = 0x00E0
	CODEPLUG_ADDR_GENERAL_SETTINGS = 0x00E0
	CODEPLUG_ADDR_DEVICE_INFO      = 0x80
)

type (
	Codeplug struct {
		data []byte
	}
)

var (
	ChannelSize = int(unsafe.Sizeof(Channel{}))
	ZoneSize    = int(unsafe.Sizeof(Zone{}))

	// Location of channel storage in the codeplug
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

func NewCodeplug(r io.Reader) (*Codeplug, error) {
	cp := Codeplug{}
	buf := make([]byte, 8192)

	for {
		nb, err := r.Read(buf)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("failed to read codeplug: %w", err)
		}

		cp.data = append(cp.data, buf[:nb]...)
		if nb < len(buf) {
			break
		}
	}

	return &cp, nil
}

func (cp *Codeplug) Channel(n int) (*Channel, error) {
	if n <= 0 {
		return nil, fmt.Errorf("invalid channel number %d", n)
	}

	channel := NewChannel()

	// channels are 1-indexed
	n--

	block := n / 128
	n -= block * 128
	if block >= len(ChannelBlocks) {
		return nil, fmt.Errorf("channel number too large")
	}

	addr := ChannelBlocks[block] + 16 + (n * ChannelSize)
	err := binary.Read(bytes.NewReader(cp.data[addr:]), binary.LittleEndian, channel)

	return channel, err
}

func (cp *Codeplug) ChannelIter() func() (int, *Channel, bool) {
	n := 1

	return func() (int, *Channel, bool) {
		for {
			ch, err := cp.Channel(n)
			n++
			if err != nil {
				return n, nil, false
			}
			if ch.Name[0] == 0xff {
				continue
			}

			return n, ch, true
		}
	}
}

func (cp *Codeplug) Zone(n int) (*Zone, error) {
	if n <= 0 {
		return nil, fmt.Errorf("invalid zone number")
	}

	// zones are 1-indexed
	n--

	if n > CODEPLUG_ZONE_MAX_COUNT {
		return nil, fmt.Errorf("zone number too large")
	}

	zone := Zone{}
	addr := CODEPLUG_ADDR_EX_ZONE_LIST + (n * ZoneSize)
	if err := binary.Read(bytes.NewReader(cp.data[addr:]), binary.LittleEndian, &zone); err != nil {
		return nil, fmt.Errorf("failed to read zone from codeplug: %w", err)
	}

	return &zone, nil
}

func (cp *Codeplug) ZoneIter() func() (int, *Zone, bool) {
	n := 1

	return func() (int, *Zone, bool) {
		for {
			zone, err := cp.Zone(n)
			n++
			if err != nil {
				return n, nil, false
			}
			if zone.Name[0] == 0xff {
				continue
			}

			return n, zone, true
		}
	}
}

func (cp *Codeplug) Settings() (*Settings, error) {
	settings := Settings{}

	addr := CODEPLUG_ADDR_GENERAL_SETTINGS
	if err := binary.Read(bytes.NewReader(cp.data[addr:]), binary.LittleEndian, &settings); err != nil {
		return nil, fmt.Errorf("failed to read radio settings: %w", err)
	}

	return &settings, nil
}

func (cp *Codeplug) DeviceInfo() (*DeviceInfo, error) {
	info := DeviceInfo{}

	if err := binary.Read(bytes.NewReader(cp.data[CODEPLUG_ADDR_DEVICE_INFO:]), binary.LittleEndian, &info); err != nil {
		return nil, fmt.Errorf("failed to read device info: %w", err)
	}

	return &info, nil
}
