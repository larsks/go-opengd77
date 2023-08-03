package opengd77

func FromBCD(val int) int {
	acc := 0
	mult := 1

	for i := 0; val > 0; i++ {
		digit := val & 0xf
		val >>= 4
		acc += digit * mult
		mult *= 10
	}

	return acc
}

func ToBCD(val int) int {
	var acc int
	for i := 0; val > 0; i++ {
		digit := val - ((val / 10) * 10)
		val = val / 10
		acc |= digit << (4 * i)
	}

	return acc
}
