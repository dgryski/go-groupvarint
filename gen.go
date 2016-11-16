// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
)

func generateBitsUsed() {

	fmt.Println("package groupvarint")
	fmt.Println()
	fmt.Println("var bytesUsed = []int{")

	for i := 0; i < 256; i++ {
		b := byte(i)

		used := 1

		for j := 0; j < 4; j++ {
			used += int(b&3) + 1
			b >>= 2
		}

		if i == 0 || (i > 1 && i%16 == 1) {
			fmt.Printf("\t")
		}
		fmt.Printf("%d, ", used)
		if i > 0 && i%16 == 0 {
			fmt.Printf("\n")
		}
	}

	fmt.Println()
	fmt.Println("}")
}

func generateSSEMasks() {

	fmt.Println("// +build amd64 !noasm")
	fmt.Println()
	fmt.Println("package groupvarint")
	fmt.Println()
	fmt.Println("var sseMasks = []uint64{")

	for i := uint(0); i < 256; i++ {

		var offs uint32
		var vals [4]uint32

		for j := uint(0); j < 4; j++ {
			d := 1 + ((i >> (2 * j)) & 3)

			for k := uint(0); k < d; k++ {
				vals[j] |= offs << (8 * k)
				offs++
			}

			for k := d; k < 4; k++ {
				vals[j] |= 0xff << (8 * k)
			}
		}

		fmt.Printf("\t0x%08x%08x, 0x%08x%08x,\n", vals[1], vals[0], vals[3], vals[2])
	}

	fmt.Println("}")
}

func main() {

	table := flag.String("table", "", "which table to generate: bytesused,ssemasks")

	flag.Parse()

	switch *table {
	case "bytesused":
		generateBitsUsed()
	case "ssemasks":
		generateSSEMasks()
	default:
		log.Fatalf("unknown table: %q", *table)
	}
}
