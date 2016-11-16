// +build !amd64 noasm

package groupvarint

func nlz(x uint32) uint32 {

	var n uint32

	if x == 0 {
		return (32)
	}

	if x <= 0x0000FFFF {
		n = n + 16
		x = x << 16
	}

	if x <= 0x00FFFFFF {
		n = n + 8
		x = x << 8
	}

	if x <= 0x0FFFFFFF {
		n = n + 4
		x = x << 4
	}

	if x <= 0x3FFFFFFF {
		n = n + 2
		x = x << 2
	}

	if x <= 0x7FFFFFFF {
		n = n + 1
	}

	return n
}
