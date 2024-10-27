package bitboard

import (
	"math/bits"
	"strings"
)

type Bits uint64

func (b Bits) String() string {
	s := []byte(strings.Repeat("0 0 0 0 0 0 0 0\n", 8))
	mask := uint64(1 << 63)
	idx := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			val := uint64(b)&mask == mask
			if val {
				s[idx] = '1'
			}
			mask >>= 1
			idx += 2
		}
	}
	return string(s)
}

// lsb_64_table is the lookup table for LSB index calculation.
var lsb64Table = [64]int{
	63, 30, 3, 32, 59, 14, 11, 33,
	60, 24, 50, 9, 55, 19, 21, 34,
	61, 29, 2, 53, 51, 23, 41, 18,
	56, 28, 1, 43, 46, 27, 0, 35,
	62, 31, 58, 4, 5, 49, 54, 6,
	15, 52, 12, 40, 7, 42, 45, 16,
	25, 57, 48, 13, 10, 39, 8, 44,
	20, 47, 38, 22, 17, 37, 36, 26,
}

func (b Bits) NextBit() int {
	return bits.TrailingZeros64(uint64(b))
}

func (b Bits) NextBitOld() int {
	b ^= b - 1
	folded := uint32(b ^ (b >> 32)) // Fold the upper 32 bits into the lower 32 bits
	return lsb64Table[(folded*0x78291ACF)>>26]
}
