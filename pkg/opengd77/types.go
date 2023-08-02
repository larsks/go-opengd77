package opengd77

type (
	Channel struct {
		Name            [16]byte
		RxFreq          [4]byte
		TxFreq          [4]byte
		Mode            byte
		rxRefFreq       byte
		txRefFreq       byte
		TimeOutTimer    byte
		totRekey        byte
		admin           byte
		rssiThreshold   byte
		scanlistIndex   byte
		rxTone          [2]byte
		txTone          [2]byte
		voiceEmphasis   byte
		txSig           byte
		unmuteRule      byte
		rxSig           byte
		artsInterval    byte
		encrypt         byte
		RxColor         byte
		RxGrouplist     byte
		TxColor         byte
		EmergencySystem byte
		ConcatNumber    uint16
		Flag1           byte
		Flag2           byte
		Flag3           byte
		Flag4           byte
		VFOOffset       uint16
		VFOFlag         byte
		Squelch         byte
	}
)
