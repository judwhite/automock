package bitboard

const (
	StartPos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	White = 0
	Black = 1

	Pawn   = 0
	Knight = 1
	Bishop = 2
	Rook   = 3
	Queen  = 4
	King   = 5
)

type Board struct {
	Pieces [2][6]Bits
	Units  [2]Bits
	All    Bits
}

func (b Board) IsBlocked(sq1, sq2 int) bool {
	return b.All&BitBetween[sq1][sq2] != 0
}
