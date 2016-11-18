package groupvarint

func Encode1(dst []byte, n uint32) []byte {
	for n >= 0x80 {
		b := byte(n) | 0x80
		dst = append(dst, b)
		n >>= 7
	}
	return append(dst, byte(n))
}

func Decode1(dst *uint32, src []byte) int {

	var n uint32

	var shift uint
	for i, s := range src[:4] {
		if s < 0x80 {
			n |= uint32(s) << shift
			*dst = n
			return i + 1
		}

		n |= uint32(s&^0x80) << shift
		shift += 7
	}

	*dst = n | uint32(src[4])<<shift
	return 5
}
