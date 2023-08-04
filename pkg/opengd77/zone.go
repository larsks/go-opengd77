package opengd77

import "bytes"

const (
	CODEPLUG_ADDR_EX_ZONE_BASIC             = 0x8000
	CODEPLUG_ADDR_EX_ZONE_INUSE_PACKED_DATA = 0x8010
	CODEPLUG_EX_ZONE_INUSE_PACKED_DATA_SIZE = 32
	CODEPLUG_ADDR_EX_ZONE_LIST              = 0x8030
	CODEPLUG_ZONE_MAX_COUNT                 = 250
)

type (
	Zone struct {
		Name     PaddedName
		Channels [80]uint16
	}
)

func NewZone() *Zone {
	return &Zone{
		Name: [16]byte(bytes.Repeat([]byte{0xff}, 16)),
	}
}
