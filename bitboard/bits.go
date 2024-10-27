package bitboard

import "strings"

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
