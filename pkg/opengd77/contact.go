package opengd77

type (
	Contact struct {
		Name             PaddedName
		TalkGroupNumber  uint32
		CallType         byte
		CallRxTone       byte
		RingStyle        byte
		TimeslotOverride PackedTimeslotOverride
		IndexNumber      int
	}

	PackedTimeslotOverride byte

	TimeslotOverride struct {
		HasOverride bool
		Timeslot    int
	}
)

func (v PackedTimeslotOverride) Unpack() *TimeslotOverride {
	return &TimeslotOverride{
		HasOverride: v&0x01 == 0,
		Timeslot:    int(v & 0x02 >> 1),
	}
}

func (v *TimeslotOverride) Pack() (res PackedTimeslotOverride) {
	if v.HasOverride {
		res &= 0xfe
	} else {
		res |= 0x01
	}

	res |= PackedTimeslotOverride(v.Timeslot << 1)

	return
}
