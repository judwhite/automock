package bitboard

import "fmt"

var (
	BitBetween [64][64]Bits
	BitAfter   [64][64]Bits

	KingMoves   [64]Bits
	KnightMoves [64]Bits
	// TODO: PawnMoves
	// TODO: PawnCaptures
	// TODO: KnightMoves
)

const (
	FileA = 7
	FileB = 6
	FileC = 5
	FileD = 4
	FileE = 3
	FileF = 2
	FileG = 1
	FileH = 0

	Rank1 = 0
	Rank2 = 8 * 1
	Rank3 = 8 * 2
	Rank4 = 8 * 3
	Rank5 = 8 * 4
	Rank6 = 8 * 5
	Rank7 = 8 * 6
	Rank8 = 8 * 7

	A1 = FileA + Rank1
	A2 = FileA + Rank2
	A3 = FileA + Rank3
	A4 = FileA + Rank4
	A5 = FileA + Rank5
	A6 = FileA + Rank6
	A7 = FileA + Rank7
	A8 = FileA + Rank8

	B1 = FileB + Rank1
	B2 = FileB + Rank2
	B3 = FileB + Rank3
	B4 = FileB + Rank4
	B5 = FileB + Rank5
	B6 = FileB + Rank6
	B7 = FileB + Rank7
	B8 = FileB + Rank8

	C1 = FileC + Rank1
	C2 = FileC + Rank2
	C3 = FileC + Rank3
	C4 = FileC + Rank4
	C5 = FileC + Rank5
	C6 = FileC + Rank6
	C7 = FileC + Rank7
	C8 = FileC + Rank8

	D1 = FileD + Rank1
	D2 = FileD + Rank2
	D3 = FileD + Rank3
	D4 = FileD + Rank4
	D5 = FileD + Rank5
	D6 = FileD + Rank6
	D7 = FileD + Rank7
	D8 = FileD + Rank8

	E1 = FileE + Rank1
	E2 = FileE + Rank2
	E3 = FileE + Rank3
	E4 = FileE + Rank4
	E5 = FileE + Rank5
	E6 = FileE + Rank6
	E7 = FileE + Rank7
	E8 = FileE + Rank8

	F1 = FileF + Rank1
	F2 = FileF + Rank2
	F3 = FileF + Rank3
	F4 = FileF + Rank4
	F5 = FileF + Rank5
	F6 = FileF + Rank6
	F7 = FileF + Rank7
	F8 = FileF + Rank8

	G1 = FileG + Rank1
	G2 = FileG + Rank2
	G3 = FileG + Rank3
	G4 = FileG + Rank4
	G5 = FileG + Rank5
	G6 = FileG + Rank6
	G7 = FileG + Rank7
	G8 = FileG + Rank8

	H1 = FileH + Rank1
	H2 = FileH + Rank2
	H3 = FileH + Rank3
	H4 = FileH + Rank4
	H5 = FileH + Rank5
	H6 = FileH + Rank6
	H7 = FileH + Rank7
	H8 = FileH + Rank8
)

var squareNames [64]string

var ranks = []uint64{
	0xFF000000_00000000,
	0x00FF0000_00000000,
	0x0000FF00_00000000,
	0x000000FF_00000000,
	0x00000000_FF000000,
	0x00000000_00FF0000,
	0x00000000_0000FF00,
	0x00000000_000000FF,
}

var files = []uint64{
	0x80808080_80808080,
	0x40404040_40404040,
	0x20202020_20202020,
	0x10101010_10101010,
	0x08080808_08080808,
	0x04040404_04040404,
	0x02020202_02020202,
	0x01010101_01010101,
}

var diagonals = []uint64{
	0x4080000000000000,
	0x2040800000000000,
	0x1020408000000000,
	0x0810204080000000,
	0x0408102040800000,
	0x0204081020408000,
	0x0102040810204080,
	0x0001020408102040,
	0x0000010204081020,
	0x0000000102040810,
	0x0000000001020408,
	0x0000000000010204,
	0x0000000000000102,
	0x0201000000000000,
	0x0402010000000000,
	0x0804020100000000,
	0x1008040201000000,
	0x2010080402010000,
	0x4020100804020100,
	0x8040201008040201,
	0x0080402010080402,
	0x0000804020100804,
	0x0000008040201008,
	0x0000000080402010,
	0x0000000000804020,
	0x0000000000008040,
}

func init() {
	genSquareNames()
	genBitBetween()
	genBitAfter()
	genKingMoves()
	genKnightMoves()
}

func genSquareNames() {
	for file := 'a'; file <= 'h'; file++ {
		for rank := 1; rank <= 8; rank++ {
			pos := (rank-1)*8 + 7 - int(file-'a')
			squareNames[pos] = fmt.Sprintf("%c%d", file, rank)
		}
	}
}

func genBitBetween() {
	lines := make([]uint64, 0, len(ranks)+len(files)+len(diagonals))
	lines = append(lines, ranks...)
	lines = append(lines, files...)
	lines = append(lines, diagonals...)

	for i := 0; i < 64; i++ {
		for j := i + 1; j < 64; j++ {
			b1, b2 := uint64(1<<i), uint64(1<<j)

			var b Bits

			for _, diag := range lines {
				if b1&diag == b1 && b2&diag == b2 {
					// b1 is always less than b2

					// anti1 = all bits > b1
					anti1 := ^(b1 | (b1 - 1))
					// anti2 = all bits < b2
					anti2 := b2 - 1

					// diag & (all greater than b1) & (all less than b2)
					b = Bits(diag & anti1 & anti2)
					break
				}
			}

			BitBetween[i][j] = b
			BitBetween[j][i] = b
		}
	}
}

func genBitAfter() {
	lines := make([]uint64, 0, len(ranks)+len(files)+len(diagonals))
	lines = append(lines, ranks...)
	lines = append(lines, files...)
	lines = append(lines, diagonals...)

	for i := 0; i < 64; i++ {
		for j := i + 1; j < 64; j++ {
			sq1, sq2 := uint64(1<<i), uint64(1<<j)

			var b1, b2 Bits

			for _, diag := range lines {
				if sq1&diag == sq1 && sq2&diag == sq2 {
					// sq1 is always less than sq2

					anti1 := ^(sq2 - 1)      // all bits >= sq2
					anti2 := sq1 | (sq1 - 1) // all bits <= sq1

					b1 = Bits(diag & anti1)
					b2 = Bits(diag & anti2)
					break
				}
			}

			BitAfter[i][j] = b1
			BitAfter[j][i] = b2
		}
	}
}

func genKingMoves() {
	for i := 0; i < 64; i++ {
		rank, file := i/8, i%8

		var b Bits
		if rank < 7 {
			n := rank + 1

			b |= 1 << (n*8 + file) // n
			if file > 0 {
				b |= 1 << (n*8 + file - 1) // ne
			}
			if file < 7 {
				b |= 1 << (n*8 + file + 1) // nw
			}
		}

		if rank > 0 {
			s := rank - 1

			b |= 1 << (s*8 + file) // s
			if file > 0 {
				b |= 1 << (s*8 + file - 1) // se
			}
			if file < 7 {
				b |= 1 << (s*8 + file + 1) // sw
			}
		}

		if file > 0 {

			b |= 1 << (rank*8 + file - 1) // e
		}
		if file < 7 {

			b |= 1 << (rank*8 + file + 1) // w
		}
		KingMoves[i] = b
	}
}

func genKnightMoves() {
	moves := []struct {
		r int
		f int
	}{
		{2, 1},   // 1 o'clock
		{1, 2},   // 2 o'clock
		{-1, 2},  // 4 o'clock
		{-2, 1},  // 5 o'clock
		{-2, -1}, // 7 o'clock
		{-1, -2}, // 8 o'clock
		{1, -2},  // 10 o'clock
		{2, -1},  // 11 o'clock
	}

	for i := 0; i < 64; i++ {
		rank, file := i/8, i%8

		var b Bits
		for _, move := range moves {
			newRank := rank + move.r
			newFile := file + move.f

			if newRank < 0 || newRank > 7 ||
				newFile < 0 || newFile > 7 {
				continue
			}

			b |= 1 << (newRank*8 + newFile)
		}

		KnightMoves[i] = b
	}
}
