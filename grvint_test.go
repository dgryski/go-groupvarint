package groupvarint

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestRoundtrip(t *testing.T) {

	for i := 0; i < 10000; i++ {

		var u32s [4]uint32

		b := rand.Intn(256)

		for j := 0; j < 4; j++ {
			size := int(b & 3)
			switch size {
			case 0:
				u32s[j] = uint32(rand.Intn(1 << 8))
			case 1:
				u32s[j] = 1<<8 + uint32(rand.Intn((1<<16)-(1<<8)))
			case 2:
				u32s[j] = 1<<16 + uint32(rand.Intn((1<<24)-(1<<16)))
			case 3:
				u32s[j] = 1<<24 + uint32(rand.Intn((1<<32)-(1<<24)))
			}

			b >>= 2
		}

		var dst [17]byte

		d := Encode4(dst[:], u32s[:])

		if bytesUsed[d[0]] != len(d) {
			t.Errorf("bytesUsed[%d]=%d, want %d\n", d[0], bytesUsed[d[0]], len(d))
		}

		var got [4]uint32
		Decode4(got[:], dst[:])

		if !reflect.DeepEqual(u32s, got) {
			t.Fatalf("failed roundtrip: got=%x, want %x (src[0]=%08b)\n", got, u32s, dst[0])
		}
	}
}

func makeInput(n int) []byte {

	rand.Seed(0)

	var input []byte

	var padding int

	for i := 0; i < n; i++ {

		var u32s [4]uint32

		b := rand.Intn(256)

		for j := 0; j < 4; j++ {
			size := int(b & 3)
			switch size {
			case 0:
				u32s[j] = uint32(rand.Intn(1 << 8))
			case 1:
				u32s[j] = 1<<8 + uint32(rand.Intn((1<<16)-(1<<8)))
			case 2:
				u32s[j] = 1<<16 + uint32(rand.Intn((1<<24)-(1<<16)))
			case 3:
				u32s[j] = 1<<24 + uint32(rand.Intn((1<<32)-(1<<24)))
			}

			b >>= 2
		}

		var dst [17]byte

		d := Encode4(dst[:], u32s[:])

		padding = 17 - len(d)

		input = append(input, d...)
	}

	// must be able to load 17 bytes from start of final block
	for i := 0; i < padding; i++ {
		input = append(input, 0)
	}

	return input
}

var sink uint32

func BenchmarkDecode(b *testing.B) {

	input := makeInput(4096)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		src := input
		for len(src) > 17 {
			var dst [4]uint32
			Decode4(dst[:], src)
			sink += dst[0] + dst[1] + dst[2] + dst[3]
			src = src[bytesUsed[src[0]]:]
		}
	}
}
