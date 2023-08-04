package opengd77

type (
	Settings struct {
		RadioName               [8]byte // callsign
		RadioId                 uint32  // DMD Id
		LibreDMRCodeplugVersion byte
		_                       [3]byte // was 4 bytes of reserved data but we use the first byte as the codeplug version
		ArsInitDelay            byte
		TxPreambleDuration      byte
		MonitorType             byte
		VoxSense                byte
		RxLowBatt               byte
		CallAlertDuration       byte
		RespTimer               byte
		ReminderTimer           byte
		GrpHang                 byte
		PrivateHang             byte
		Flag1                   byte
		Flag2                   byte
		Flag3                   byte
		Flag4                   byte
		_                       [2]byte
		ProgrammingPassword     [8]byte
	}
)
