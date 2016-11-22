package groupvarint

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestRoundtrip(t *testing.T) {

	input := makeInput(4096)

	for len(input) > 0 {
		u32s := input[:4]
		var dst [17]byte
		d := Encode4(dst[:], u32s)

		if BytesUsed[d[0]] != len(d) {
			t.Errorf("BytesUsed[%d]=%d, want %d\n", d[0], BytesUsed[d[0]], len(d))
		}

		var got [4]uint32
		Decode4(got[:], dst[:])

		if !reflect.DeepEqual(u32s, got[:]) {
			t.Fatalf("failed roundtrip: got=%x, want %x (src[0]=%08b)\n", got, u32s, dst[0])
		}

		input = input[4:]
	}
}

func makeInput(n int) []uint32 {
	rand.Seed(0)

	var input []uint32

	for i := 0; i < n; i++ {

		for j := 0; j < 4; j++ {

			b := uint32(rand.Int31())

			size := nlz(b)

			var u32 uint32

			switch size {
			// case 0: none, because b > 0
			case 1:
				u32 = uint32(rand.Intn(1 << 8))
			case 2:
				u32 = 1<<8 + uint32(rand.Intn((1<<16)-(1<<8)))
			case 3:
				u32 = 1<<16 + uint32(rand.Intn((1<<24)-(1<<16)))
			default:
				u32 = 1<<24 + uint32(rand.Intn((1<<32)-(1<<24)))
			}

			input = append(input, u32)
		}

	}

	return input
}

func encodeGroupVarint(input []uint32) []byte {

	var r []byte

	var padding int
	for len(input) > 0 {
		var dst [17]byte

		d := Encode4(dst[:], input)

		padding = 17 - len(d)

		r = append(r, d...)

		input = input[4:]
	}

	// must be able to load 17 bytes from start of final block
	for i := 0; i < padding; i++ {
		r = append(r, 0)
	}

	return r
}

func encodeVarint(input []uint32) []byte {
	var r []byte
	for _, u32 := range input {
		var dst [5]byte
		d := Encode1(dst[:], u32)
		r = append(r, d...)
	}

	return r
}

var sink uint32

func BenchmarkDecode(b *testing.B) {

	input := encodeGroupVarint(makeInput(4096))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		src := input
		for len(src) > 17 {
			var dst [4]uint32
			Decode4(dst[:], src)
			sink += dst[0] + dst[1] + dst[2] + dst[3]
			src = src[BytesUsed[src[0]]:]
		}
	}
}

func BenchmarkVint(b *testing.B) {

	input := encodeVarint(makeInput(4096))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		src := input
		for len(src) > 0 {
			var dst uint32
			used := Decode1(&dst, src)
			sink += dst
			src = src[used:]
		}
	}
}

func BenchmarkBaseline(b *testing.B) {

	input := makeInput(4096)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		src := input
		for len(src) > 0 {
			sink += src[0] + src[1] + src[2] + src[3]
			src = src[4:]
		}
	}
}
