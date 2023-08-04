package opengd77

import (
	"bytes"
)

const (
	TA_TX_OFF  int = 0
	TA_TX_APRS     = 1
	TA_TX_TEXT     = 2
	TA_TX_BOTH     = 3

	// LibreDMRFlag1
	CODEPLUG_CHANNEL_LIBREDMR_FLAG1_OPTIONAL_DMRID = 0x80
	CODEPLUG_CHANNEL_LIBREDMR_FLAG1_NO_BEEP        = 0x40
	CODEPLUG_CHANNEL_LIBREDMR_FLAG1_NO_ECO         = 0x20
	CODEPLUG_CHANNEL_LIBREDMR_FLAG1_OUT_OF_BAND    = 0x10

	// flag2
	CODEPLUG_CHANNEL_FLAG2_TIMESLOT_TWO = 0x40

	// flag3
	CODEPLUG_CHANNEL_FLAG3_DISABLE_ALL_LEDS = 0x20

	// flag4
	CODEPLUG_CHANNEL_FLAG4_SQUELCH   = 0x01
	CODEPLUG_CHANNEL_FLAG4_BW_25K    = 0x02
	CODEPLUG_CHANNEL_FLAG4_RX_ONLY   = 0x04
	CODEPLUG_CHANNEL_FLAG4_ALL_SKIP  = 0x10
	CODEPLUG_CHANNEL_FLAG4_ZONE_SKIP = 0x20
	CODEPLUG_CHANNEL_FLAG4_VOX       = 0x40
	CODEPLUG_CHANNEL_FLAG4_POWER     = 0x80
)

type (
	LibreDMRFlag1 struct {
		DMRId     bool
		NoBeep    bool
		NoEco     bool
		OutOfBand bool
	}

	ChannelFlag4 struct {
		Squelch  bool
		BW25K    bool
		RxOnly   bool
		AllSkip  bool
		ZoneSkip bool
		Vox      bool
		Power    bool
	}

	Channel struct {
		Name          PaddedName
		RxFreq        BCDFrequency
		TxFreq        BCDFrequency
		Mode          byte
		LibreDMRPower byte // was RxRefFreq
		TxRefFreq     byte
		TimeOutTimer  byte
		TotRekey      byte
		AdmitCriteria byte
		RssiThreshold byte
		ScanlistIndex byte
		RxTone        BCDTone
		TxTone        BCDTone
		VoiceEmphasis byte
		TxSig         byte
		LibreDMRFlag1 byte // was UnmuteRule

		// The following three bytes are the optional DMR ID
		// if CODEPLUG_CHANNEL_LIBREDMR_FLAG1_OPTIONAL_DMRID
		// is set.
		RxSig        byte
		ArtsInterval byte
		Encrypt      byte

		RxColor         byte
		RxGrouplist     byte
		TxColor         byte
		EmergencySystem byte
		ContactNumber   uint16
		Flag1           byte // lower 6 bits TA Tx control
		Flag2           byte // bits...0x40 = TS
		Flag3           byte // bits... 0x20 = DisableAllLeds
		Flag4           byte // bits... 0x80 = Power, 0x40 = Vox, 0x20 = ZoneSkip (AutoScan), 0x10 = AllSkip (LoneWoker), 0x08 = AllowTalkaround, 0x04 = OnlyRx, 0x02 = Channel width, 0x01 = Squelch
		VFOOffset       uint16
		VFOFlag         byte // uppder 4 bits are the step frequency (2.5, 5, 6.25, 10, 12.5, 25, 30, 50)
		Squelch         byte
	}
)

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

func ChannelFlag4FromInt(v int) *ChannelFlag4 {
	return &ChannelFlag4{
		Squelch:  v&CODEPLUG_CHANNEL_FLAG4_SQUELCH != 0,
		BW25K:    v&CODEPLUG_CHANNEL_FLAG4_BW_25K != 0,
		RxOnly:   v&CODEPLUG_CHANNEL_FLAG4_RX_ONLY != 0,
		AllSkip:  v&CODEPLUG_CHANNEL_FLAG4_ALL_SKIP != 0,
		ZoneSkip: v&CODEPLUG_CHANNEL_FLAG4_ZONE_SKIP != 0,
		Vox:      v&CODEPLUG_CHANNEL_FLAG4_VOX != 0,
		Power:    v&CODEPLUG_CHANNEL_FLAG4_POWER != 0,
	}
}

func LibreDMRFlag1FromInt(v int) *LibreDMRFlag1 {
	return &LibreDMRFlag1{
		DMRId:     v&CODEPLUG_CHANNEL_LIBREDMR_FLAG1_OPTIONAL_DMRID != 0,
		NoBeep:    v&CODEPLUG_CHANNEL_LIBREDMR_FLAG1_NO_BEEP != 0,
		NoEco:     v&CODEPLUG_CHANNEL_LIBREDMR_FLAG1_NO_ECO != 0,
		OutOfBand: v&CODEPLUG_CHANNEL_LIBREDMR_FLAG1_OUT_OF_BAND != 0,
	}
}
