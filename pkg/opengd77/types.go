package opengd77

import (
	"fmt"
	"strings"
)

type (
	MemoryBlock struct {
		Kind        int
		FileOffset  int
		RadioOffset int
		Length      int
	}

	BCDFrequency uint32
	BCDTone      uint16
	PaddedName   [16]byte

	Channel struct {
		Name            PaddedName
		RxFreq          BCDFrequency
		TxFreq          BCDFrequency
		Mode            byte
		MaybePower      byte
		TxRefFreq       byte
		TimeOutTimer    byte
		TotRekey        byte
		AdmitCriteria   byte
		RssiThreshold   byte
		ScanlistIndex   byte
		RxTone          BCDTone
		TxTone          BCDTone
		VoiceEmphasis   byte
		TxSig           byte
		UnmuteRule      byte
		RxSig           byte
		ArtsInterval    byte
		Encrypt         byte
		RxColor         byte
		RxGrouplist     byte
		TxColor         byte
		EmergencySystem byte
		ContactNumber   uint16
		Flag1           byte
		Flag2           byte
		Flag3           byte
		Flag4           byte
		VFOOffset       uint16
		VFOFlag         byte
		Squelch         byte
	}
)

func (v BCDFrequency) String() string {
	return fmt.Sprintf("%0.4f", float64(FromBCD(int(v), 8))/10000.0)
}

func (v BCDTone) String() string {
	return fmt.Sprintf("%0.1f", float64(FromBCD(int(v), 4))/10.0)
}

func (v PaddedName) String() string {
	return strings.TrimRight(string(v[:]), "\xff")
}
