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

	ActiveColor int
}

func (b Board) IsBlocked(sq1, sq2 int) bool {
	return b.All&BitBetween[sq1][sq2] != 0
}

func (b Board) Attack(color int, sq int) bool {
	if b.Pieces[color][Pawn]&PawnDefends[color][sq] != 0 {
		return true
	}
	if b.Pieces[color][Knight]&PieceMoves[Knight][sq] != 0 {
		return true
	}
	if b.Pieces[color][King]&PieceMoves[King][sq] != 0 {
		return true
	}

	b1 := PieceMoves[Rook][sq] & (b.Pieces[color][Rook] | b.Pieces[color][Queen])
	b1 |= PieceMoves[Bishop][sq] & (b.Pieces[color][Bishop] | b.Pieces[color][Queen])

	for b1 != 0 {
		sq2 := NextBit(b1)
		if BitBetween[sq2][sq]&b.All == 0 {
			return true
		}
		b1 &= ^(1 << sq2) // TODO: test if pre-calculating a `not_mask` is faster
	}

	return false
}

func (b Board) LegalMoves() {
	pieces := []int{Knight, Bishop, Rook, Queen, King}

	for i := 0; i < 64; i++ {
		bit := Bits(1 << i)

		if b.Units[b.ActiveColor]&bit == 0 {
			continue
		}

		if b.Pieces[b.ActiveColor][Pawn]&bit == bit {
			//TODO: check en passant target square
			//PawnMoves[b.ActiveColor][i]
			//PawnCaptures[b.ActiveColor][i]
		}

		for _, piece := range pieces {
			if b.Pieces[b.ActiveColor][piece]&bit == bit {
				//TODO: knights will get special treatment
				//PieceMoves[piece][i]
			}
		}
	}
}

func (b Board) String() string {
	/*
	   +---+---+---+---+---+---+---+---+
	   | r | n |   | q |   | b | n | r | 8
	   +---+---+---+---+---+---+---+---+
	   | p | p | p |   | k | B | p |   | 7
	   +---+---+---+---+---+---+---+---+
	   |   |   |   | p |   |   |   | p | 6
	   +---+---+---+---+---+---+---+---+
	   |   |   |   |   | N |   |   |   | 5
	   +---+---+---+---+---+---+---+---+
	   |   |   |   |   | P |   |   |   | 4
	   +---+---+---+---+---+---+---+---+
	   |   |   | N |   |   |   |   |   | 3
	   +---+---+---+---+---+---+---+---+
	   | P | P | P | P |   | P | P | P | 2
	   +---+---+---+---+---+---+---+---+
	   | R |   | B | b | K |   |   | R | 1
	   +---+---+---+---+---+---+---+---+
	     a   b   c   d   e   f   g   h

	*/
	return "TODO"
}
