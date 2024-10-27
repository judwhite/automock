package bitboard

import (
	"fmt"
)

func init() {
	genSquareNames()
	genBitBetween()
	genBitAfter()

	genKingMoves()
	genKnightMoves()
	genPawnCaptures()
	genPawnMoves()

	genPawnWeaknesses()
	genSliderMoves()

	//b, err := ParseFEN(StartPos)
	//if err != nil {
	//	panic(err)
	//}
	//
	//bb := b.All
	//fmt.Printf("%s\n\n---\n\n", bb.String())
	//nextBit := NextBit(bb)
	//fmt.Printf("next bit: %d\n", nextBit)
}

var (
	BitBetween [64][64]Bits
	BitAfter   [64][64]Bits

	PawnCaptures [2][64]Bits
	PawnDefends  [2][64]Bits
	PawnMoves    [2][64]Bits

	PieceMoves [6][64]Bits

	PassedPawns   [2][64]Bits
	IsolatedPawns [64]Bits
	PawnPath      [2][64]Bits
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

var ranks = []Bits{
	0xFF000000_00000000,
	0x00FF0000_00000000,
	0x0000FF00_00000000,
	0x000000FF_00000000,
	0x00000000_FF000000,
	0x00000000_00FF0000,
	0x00000000_0000FF00,
	0x00000000_000000FF,
}

var files = []Bits{
	0x80808080_80808080,
	0x40404040_40404040,
	0x20202020_20202020,
	0x10101010_10101010,
	0x08080808_08080808,
	0x04040404_04040404,
	0x02020202_02020202,
	0x01010101_01010101,
}

var diagonals = []Bits{
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

func genSquareNames() {
	for file := 'a'; file <= 'h'; file++ {
		for rank := 1; rank <= 8; rank++ {
			pos := (rank-1)*8 + 7 - int(file-'a')
			squareNames[pos] = fmt.Sprintf("%c%d", file, rank)
		}
	}
}

func genBitBetween() {
	lines := make([]Bits, 0, len(ranks)+len(files)+len(diagonals))
	lines = append(lines, ranks...)
	lines = append(lines, files...)
	lines = append(lines, diagonals...)

	for i := 0; i < 64; i++ {
		for j := i + 1; j < 64; j++ {
			b1, b2 := Bits(1<<i), Bits(1<<j)

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
	lines := make([]Bits, 0, len(ranks)+len(files)+len(diagonals))
	lines = append(lines, ranks...)
	lines = append(lines, files...)
	lines = append(lines, diagonals...)

	for i := 0; i < 64; i++ {
		for j := i + 1; j < 64; j++ {
			sq1, sq2 := Bits(1<<i), Bits(1<<j)

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
		PieceMoves[King][i] = b
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

		PieceMoves[Knight][i] = b
	}
}

func genPawnCaptures() {
	for rank := 1; rank <= 6; rank++ {
		for file := 0; file <= 7; file++ {
			pos := rank*8 + (7 - file)
			posBits := Bits(1 << pos)

			if file > 0 {
				newFile := file - 1

				whiteSq := (rank+1)*8 + (7 - newFile)
				blackSq := (rank-1)*8 + (7 - newFile)

				PawnCaptures[White][pos] |= 1 << whiteSq
				PawnCaptures[Black][pos] |= 1 << blackSq

				PawnDefends[White][whiteSq] |= posBits
				PawnDefends[Black][blackSq] |= posBits
			}
			if file < 7 {
				newFile := file + 1

				whiteSq := (rank+1)*8 + (7 - newFile)
				blackSq := (rank-1)*8 + (7 - newFile)

				PawnCaptures[White][pos] |= 1 << whiteSq
				PawnCaptures[Black][pos] |= 1 << blackSq

				PawnDefends[White][whiteSq] |= posBits
				PawnDefends[Black][blackSq] |= posBits
			}
		}
	}
}

func genPawnMoves() {
	for rank := 1; rank <= 6; rank++ {
		for file := 0; file <= 7; file++ {
			pos := rank*8 + (7 - file)

			var white, black Bits

			white |= Bits(1 << (pos + 8))
			black |= Bits(1 << (pos - 8))

			if rank == 1 {
				white |= Bits(1 << (pos + 16))
			} else if rank == 6 {
				black |= Bits(1 << (pos - 16))
			}

			PawnMoves[White][pos] = white
			PawnMoves[Black][pos] = black
		}
	}
}

func genPawnWeaknesses() {
	for x := 0; x < 64; x++ {
		colX := x % 8
		rowX := x / 8
		for y := 0; y < 64; y++ {
			colY := y % 8
			rowY := y / 8

			if abs(colX-colY) < 2 {
				if rowX < rowY && rowY < 7 {
					PassedPawns[White][x] |= 1 << y
				}
				if rowX > rowY && rowY > 0 {
					PassedPawns[Black][x] |= 1 << y
				}
			}

			if abs(colX-colY) == 1 {
				IsolatedPawns[x] |= 1 << y
			}

			if colX == colY {
				if rowX < rowY {
					PawnPath[White][x] |= 1 << y
				}
				if rowX > rowY {
					PawnPath[Black][x] |= 1 << y
				}
			}
		}
	}

	//fmt.Fprintf(os.Stderr, "%s\n\n", PawnPath[White][E4].String())
	//fmt.Fprintf(os.Stderr, "%s\n\n", PawnPath[Black][D7].String())

	// PassedPawns[White][E2]:
	// 0 0 0 0 0 0 0 0
	// 0 0 0 1 1 1 0 0
	// 0 0 0 1 1 1 0 0
	// 0 0 0 1 1 1 0 0
	// 0 0 0 1 1 1 0 0
	// 0 0 0 1 1 1 0 0
	// 0 0 0 0 0 0 0 0
	// 0 0 0 0 0 0 0 0

	// IsolatedPawns[D3]
	// 0 0 1 0 1 0 0 0
	// 0 0 1 0 1 0 0 0
	// 0 0 1 0 1 0 0 0
	// 0 0 1 0 1 0 0 0
	// 0 0 1 0 1 0 0 0
	// 0 0 1 0 1 0 0 0
	// 0 0 1 0 1 0 0 0
	// 0 0 1 0 1 0 0 0

	// PawnPath[White][E4]
	// 0 0 0 0 1 0 0 0
	// 0 0 0 0 1 0 0 0
	// 0 0 0 0 1 0 0 0
	// 0 0 0 0 1 0 0 0
	// 0 0 0 0 0 0 0 0
	// 0 0 0 0 0 0 0 0
	// 0 0 0 0 0 0 0 0
	// 0 0 0 0 0 0 0 0
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

// NextBit computes the index of the least significant bit (LSB) in the bitboard `bb`.
func NextBit(b Bits) int {
	b ^= b - 1
	folded := uint32(b ^ (b >> 32)) // Fold the upper 32 bits into the lower 32 bits
	return lsb64Table[(folded*0x78291ACF)>>26]
}

func genSliderMoves() {
	var rookLines []Bits
	rookLines = append(rookLines, files...)
	rookLines = append(rookLines, ranks...)

	for i := 0; i < 64; i++ {
		a := Bits(1 << i)
		for j := 0; j < 64; j++ {
			if i == j {
				continue
			}
			b := Bits(1 << j)

			for _, diagonal := range diagonals {
				if diagonal&a == a && diagonal&b == b {
					PieceMoves[Bishop][i] |= b
					PieceMoves[Queen][i] |= b
				}
			}

			for _, rookLine := range rookLines {
				if rookLine&a == a && rookLine&b == b {
					PieceMoves[Rook][i] |= b
					PieceMoves[Queen][i] |= b
				}
			}
		}
	}
}
