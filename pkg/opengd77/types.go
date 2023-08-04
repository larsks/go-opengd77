package opengd77

import (
	"bytes"
	"fmt"
)

type (
	BCDFrequency    uint32
	BCDTone         uint16
	BCD16           uint16
	PaddedName      [16]byte
	PaddedNameShort [8]byte
)

func (v BCDFrequency) String() string {
	return fmt.Sprintf("%0.4f", float64(FromBCD(int(v)))/100000.0)
}

func (v BCDFrequency) Float() float64 {
	return float64(FromBCD(int(v))) / 100000.0
}

func (v BCDTone) Float() float64 {
	return float64(FromBCD(int(v))) / 10.0
}

func (v BCDTone) String() string {
	if v == 0xffff {
		return "-"
	} else {
		return fmt.Sprintf("%0.1f", float64(FromBCD(int(v)))/10.0)
	}
}

func (v BCD16) String() string {
	return fmt.Sprintf("%0d", FromBCD(int(v)))
}

func (v PaddedName) String() string {
	return string(bytes.TrimRight(v[:], "\xff"))
}

func (v PaddedNameShort) String() string {
	return string(bytes.TrimRight(v[:], "\xff"))
}
