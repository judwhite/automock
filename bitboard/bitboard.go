package bitboard

import (
	"golang.org/x/xerrors"
)

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

	ActiveColor    Color
	Castle         int
	EPTargetSquare int
	HalfMoveClock  int
	FullMoveNumber int
}

type Color int

func (c Color) String() string {
	if c == White {
		return "w"
	}
	return "b"
}

func (b Board) IsBlocked(sq1, sq2 int) bool {
	return b.All&BitBetween[sq1][sq2] != 0
}

func (b Board) Attack(color int, sq int) bool {
	p := b.Pieces[color]

	if p[Pawn]&PawnDefends[color][sq] != 0 {
		return true
	}
	if p[Knight]&PieceMoves[Knight][sq] != 0 {
		return true
	}
	if p[King]&PieceMoves[King][sq] != 0 {
		return true
	}

	b1 := PieceMoves[Rook][sq] & (p[Rook] | p[Queen])
	b1 |= PieceMoves[Bishop][sq] & (p[Bishop] | p[Queen])

	for b1 != 0 {
		sq2 := b1.NextBit()
		if BitBetween[sq2][sq]&b.All == 0 {
			return true
		}
		b1 &= b1 - 1
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

func (b Board) Apply(moves []string) (Board, error) {
	bb := b

	for _, uci := range moves {
		// parse move

		if len(uci) < 4 || len(uci) > 5 {
			return Board{}, xerrors.Errorf("invalid uci move '%s'", uci)
		}

		from := uci[0:2]
		to := uci[2:4]

		fromPos, ok := squareNameToBits[from]
		if !ok {
			return Board{}, xerrors.Errorf("invalid uci move '%s'", uci)
		}

		toPos, ok := squareNameToBits[to]
		if !ok {
			return Board{}, xerrors.Errorf("invalid uci move '%s'", uci)
		}
		toIdx := squareNameToIndex[to]

		// set side (s) and opponent side (xs)
		s := bb.ActiveColor
		xs := 1 - s
		epIdx := 0
		castle := (bb.Castle >> (2 * s)) & 0b11

		// get piece type
		pieceType := bb.PieceType(fromPos, s)
		if pieceType == -1 {
			return Board{}, xerrors.Errorf("invalid uci move '%s', %s does not have a piece on %s", uci, s.String(), from)
		}

		// capture?
		resetHalfMoveClock := pieceType == Pawn
		if capturedPieceType := bb.PieceType(toPos, xs); capturedPieceType != -1 {
			// remove captured piece
			bb.Pieces[xs][capturedPieceType] &= ^toPos
			resetHalfMoveClock = true
		} else if pieceType == Pawn {
			if bb.EPTargetSquare != 0 && bb.EPTargetSquare == toIdx {
				// remove pawn captured by en passant
				var capturePos Bits
				if toIdx >= H6 {
					capturePos = toPos >> 8
				} else {
					capturePos = toPos << 8
				}
				epCapture := ^capturePos
				bb.Pieces[xs][Pawn] &= epCapture
				resetHalfMoveClock = true
			} else if epFrom&fromPos == fromPos && epTo&toPos == toPos && epMask[toIdx]&bb.Pieces[xs][Pawn] > 0 {
				epIdx = epTargetIndex[toIdx]
			}
		}

		if pieceType == King && castle != 0 {
			if PieceMoves[King][toIdx]&fromPos != fromPos {
				for sideIdx := ksIdx; sideIdx <= qsIdx; sideIdx++ {
					sideMask := sideIdx + 1
					if castle&sideMask == 0 {
						continue
					}
					if castleKingTo[s][sideIdx]&toPos == 0 {
						continue
					}
					rookFrom := castleRookFrom[s][sideIdx]
					rookTo := castleRookTo[s][sideIdx]

					bb.Pieces[s][Rook] &= ^rookFrom
					bb.Pieces[s][Rook] |= rookTo
				}
			}
			castle = 0
		} else if pieceType == Rook && castle != 0 {
			if castleRookFrom[s][ksIdx]&fromPos == fromPos {
				castle &= qs
			} else if castleRookFrom[s][qsIdx]&fromPos == fromPos {
				castle &= ks
			}
		}

		// remove piece from original square
		bb.Pieces[s][pieceType] &= ^fromPos

		if pieceType == Pawn && len(uci) == 5 {
			promotion := uci[4]
			switch promotion {
			case 'n':
				pieceType = Knight
			case 'b':
				pieceType = Bishop
			case 'r':
				pieceType = Rook
			case 'q':
				pieceType = Queen
			default:
				return Board{}, xerrors.Errorf("invalid uci move '%s', promotion piece '%c' is invalid", uci, promotion)
			}
		}

		bb.Pieces[s][pieceType] |= toPos

		bb.ActiveColor = xs

		castle = castle<<(s*2) | 0b11<<(xs*2)
		bb.Castle &= castle

		bb.EPTargetSquare = epIdx

		if resetHalfMoveClock {
			bb.HalfMoveClock = 0
		} else {
			bb.HalfMoveClock += 1
		}

		if s == Black {
			bb.FullMoveNumber += 1
		}
	}

	for i := 0; i < 2; i++ {
		bb.Units[i] = bb.Pieces[i][Pawn] |
			bb.Pieces[i][Knight] |
			bb.Pieces[i][Bishop] |
			bb.Pieces[i][Rook] |
			bb.Pieces[i][Queen] |
			bb.Pieces[i][King]
	}

	bb.All = bb.Units[Black] | bb.Units[White]

	return bb, nil
}

func (b Board) PieceType(pos Bits, s Color) int {
	pieces := b.Pieces[s]

	if pieces[Pawn]&pos == pos {
		return Pawn
	} else if pieces[Knight]&pos == pos {
		return Knight
	} else if pieces[Bishop]&pos == pos {
		return Bishop
	} else if pieces[Rook]&pos == pos {
		return Rook
	} else if pieces[Queen]&pos == pos {
		return Queen
	} else if pieces[King]&pos == pos {
		return King
	}

	return -1
}
