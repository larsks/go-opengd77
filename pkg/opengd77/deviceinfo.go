package opengd77

type (
	DeviceInfo struct {
		MinUHFFreq          uint16
		MaxUHFFreq          uint16
		MinVHFFreq          uint16
		MaxVHFFreq          uint16
		LastProgrammingTime [6]byte
		_                   uint16
		Model               [8]byte
		Sn                  [16]byte
		CpsSoftwareVersion  [8]byte
		HardwareVersion     [8]byte
		FirmwareVersion     [8]byte
		DspFirmwareVersion  [24]byte
		_                   [8]byte
	}
)
