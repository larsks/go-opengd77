package opengd77

type (
	DeviceInfo struct {
		MinUHFFreq          BCD16
		MaxUHFFreq          BCD16
		MinVHFFreq          BCD16
		MaxVHFFreq          BCD16
		LastProgrammingTime [6]byte
		_                   uint16
		Model               PaddedNameShort
		SerialNumber        PaddedName
		CpsSoftwareVersion  PaddedNameShort
		HardwareVersion     PaddedNameShort
		FirmwareVersion     PaddedNameShort
		DspFirmwareVersion  [24]byte
		_                   [8]byte
	}
)
