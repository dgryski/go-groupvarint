package groupvarint

import "testing"

func TestRoundtripVarint(t *testing.T) {
	input := makeInput(4096)

	for _, u32 := range input {
		var dst [5]byte

		d := Encode1(dst[:], u32)

		var n uint32
		l := Decode1(&n, d)

		if l != len(d) {
			t.Errorf("Decode1(): l=%d, want %d", l, len(d))
		}

		if n != u32 {
			t.Errorf("roundtrip = %x, want %x\n", n, u32)
		}
	}
}
