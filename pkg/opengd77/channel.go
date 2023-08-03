package opengd77

import (
	"bytes"
	"fmt"
)

type (
	BCDFrequency uint32
	BCDTone      uint16
	PaddedName   [16]byte

	Channel struct {
		Name            PaddedName
		RxFreq          BCDFrequency
		TxFreq          BCDFrequency
		Mode            byte
		RxRefFreq       byte
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

func (v PaddedName) String() string {
	return string(bytes.TrimRight(v[:], "\xff"))
}

func NewChannel() *Channel {
	return &Channel{
		Name: [16]byte(bytes.Repeat([]byte{0xff}, 16)),
	}
}

func (ch *Channel) GetRxFreq() float64 {
	return ch.RxFreq.Float()
}

func (ch *Channel) GetTxFreq() float64 {
	return ch.TxFreq.Float()
}

func (ch *Channel) SetRxFreq(freq float64) {
	ch.RxFreq = BCDFrequency(ToBCD(int(freq * 100000.0)))
}

func (ch *Channel) SetTxFreq(freq float64) {
	ch.TxFreq = BCDFrequency(ToBCD(int(freq * 100000.0)))
}

func (ch *Channel) GetName() string {
	return ch.Name.String()
}

func (ch *Channel) SetName(name string) {
	var chname PaddedName = [16]byte(bytes.Repeat([]byte{0xff}, 16))
	copy(chname[:len(name)], []byte(name)[:16])

	ch.Name = chname
}

func (ch *Channel) SetRxTone(tone float64) {
	ch.RxTone = BCDTone(ToBCD(int(tone * 10.0)))
}

func (ch *Channel) SetTxTone(tone float64) {
	ch.TxTone = BCDTone(ToBCD(int(tone * 10.0)))
}

func (ch *Channel) GetRxTone() float64 {
	return ch.RxTone.Float()
}

func (ch *Channel) GetTxTone() float64 {
	return ch.TxTone.Float()
}
