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

func (b Board) Attack(s Color, sq int) bool {
	p := b.Pieces[s]

	if p[Pawn]&PawnDefends[s][sq] != 0 {
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
		b1 &= b1 - 1

		if BitBetween[sq][sq2]&b.All == 0 {
			return true
		}

	}

	return false
}

func (b Board) LegalMoves() []string {
	legalMoves := make([]string, 0, 30)

	s := b.ActiveColor
	xs := 1 - s

	sUnits := b.Units[s]
	xsUnits := b.Units[xs]

	pieces := b.Pieces[s]
	pawns := pieces[Pawn]
	pawnCaptures := PawnCaptures[s]

	for pawns != 0 {
		sq1 := pawns.NextBit()
		pawns &= pawns - 1

		// forward moves
		forwardMoves := PawnMoves[s][sq1]
		for forwardMoves != 0 {
			sq2 := forwardMoves.NextBit()
			forwardMoves &= forwardMoves - 1

			if b.IsBlocked(sq1, sq2) {
				continue
			}
			sq2Pos := Bits(1 << sq2)
			if b.All&sq2Pos == sq2Pos {
				continue
			}

			name := squareNames[sq1] + squareNames[sq2]

			// promotion squares
			if sq2 >= H8 || sq2 <= A1 {
				legalMoves = append(legalMoves,
					name+"q",
					name+"r",
					name+"n",
					name+"b",
				)
			} else {
				legalMoves = append(legalMoves, name)
			}
		}

		// captures
		captures := pawnCaptures[sq1]
		for captures != 0 {
			sq2 := captures.NextBit()
			captures &= captures - 1

			sq2Pos := Bits(1 << sq2)
			if xsUnits&sq2Pos == sq2Pos {
				name := squareNames[sq1] + squareNames[sq2]

				// promotion squares
				if sq2 >= H8 || sq2 <= A1 {
					legalMoves = append(legalMoves,
						name+"q",
						name+"r",
						name+"n",
						name+"b",
					)
				} else {
					legalMoves = append(legalMoves, name)
				}
			}

			// en-passant
			if b.EPTargetSquare != 0 && b.EPTargetSquare == sq2 {
				legalMoves = append(legalMoves, squareNames[sq1]+squareNames[sq2])
			}
		}
	}

	knights := pieces[Knight]
	for knights != 0 {
		sq1 := knights.NextBit()
		knights &= knights - 1

		knightMoves := PieceMoves[Knight][sq1]
		for knightMoves != 0 {
			sq2 := knightMoves.NextBit()
			knightMoves &= knightMoves - 1

			sq2Pos := Bits(1 << sq2)
			if sUnits&sq2Pos == sq2Pos {
				continue
			}

			legalMoves = append(legalMoves, squareNames[sq1]+squareNames[sq2])
		}
	}

	// bishop, rook, queen
	for _, pieceType := range []int{Bishop, Rook, Queen} {
		pcs := pieces[pieceType]
		for pcs != 0 {
			sq1 := pcs.NextBit()
			pcs &= pcs - 1

			pieceMoves := PieceMoves[pieceType][sq1]
			for pieceMoves != 0 {
				sq2 := pieceMoves.NextBit()
				pieceMoves &= pieceMoves - 1

				if b.IsBlocked(sq1, sq2) {
					continue
				}

				sq2Pos := Bits(1 << sq2)
				if sUnits&sq2Pos == sq2Pos {
					continue
				}

				legalMoves = append(legalMoves, squareNames[sq1]+squareNames[sq2])
			}
		}
	}

	// king
	king := pieces[King]
	kingSquare := king.NextBit()

	kingMoves := PieceMoves[King][kingSquare]
	for kingMoves != 0 {
		sq2 := kingMoves.NextBit()
		kingMoves &= kingMoves - 1

		sq2Pos := Bits(1 << sq2)
		if sUnits&sq2Pos == sq2Pos {
			continue
		}

		legalMoves = append(legalMoves, squareNames[kingSquare]+squareNames[sq2])
	}

	// castling

	castle := (b.Castle >> (2 * s)) & 0b11
	if castle != 0 {
		inCheck := b.Attack(xs, kingSquare)
		if !inCheck {
			castleSides := make([]int, 2)

			if castle&ks == ks {
				castleSides = append(castleSides, ksIdx)
			}
			if castle&qs == qs {
				castleSides = append(castleSides, qsIdx)
			}

			for _, castleSideIdx := range castleSides {
				kingDestPos := castleKingTo[s][castleSideIdx]
				kingDestSq := kingDestPos.NextBit()

				canCastle := true
				squares := BitBetween[kingSquare][kingDestSq] | kingDestPos
				if squares&b.All != 0 {
					continue
				}

				for squares != 0 {
					sq2 := squares.NextBit()
					squares &= squares - 1

					if b.Attack(xs, sq2) {
						canCastle = false
						break
					}
				}
				if canCastle {
					legalMoves = append(legalMoves, squareNames[kingSquare]+squareNames[kingDestSq])
				}
			}
		}
	}

	// remove moves that leave the king in check
	for i := 0; i < len(legalMoves); i++ {
		move := legalMoves[i]
		b2, err := b.Apply([]string{move})
		if err != nil {
			panic(err)
		}

		newKingSquare := b2.Pieces[s][King].NextBit()

		if b2.Attack(xs, newKingSquare) {
			legalMoves = append(legalMoves[:i], legalMoves[i+1:]...)
			i--
			continue
		}
	}

	return legalMoves
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
