package opengd77

func FromBCD(val, length int) int {
	var acc int
	for i := length; i > 0; i-- {
		digit := (val >> (4 * (i - 1)) & 0xf)
		acc = (acc * 10) + digit
	}

	return acc
}
