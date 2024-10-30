package bitboard

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/xerrors"
)

const (
	StartPos    = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	StartPosKey = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -"

	White = 0
	Black = 1

	Pawn   = 0
	Knight = 1
	Bishop = 2
	Rook   = 3
	Queen  = 4
	King   = 5
)

var (
	startPosBoard Board
)

func init() {
	var err error
	startPosBoard, err = ParseFEN(StartPos)
	if err != nil {
		panic(err)
	}
}

func StartPosBoard() Board {
	return startPosBoard
}

func PlyToFullMove(ply int) int {
	return (ply + 1) / 2
}

func PlyToColor(ply int) Color {
	if ply%2 == 1 {
		return White
	}
	return Black
}

type Board struct {
	Pieces [2][6]Bits
	Units  [2]Bits
	All    Bits

	ActiveColor    Color
	Castle         uint8
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

func (b Board) pseudoLegalPawnMoves() []uint64 {
	s := b.ActiveColor

	pawns := b.Pieces[s][Pawn]
	if pawns == 0 {
		return nil
	}

	moves := make([]uint64, 0, 32)

	xs := 1 - s
	enemyUnits := b.Units[xs]

	pawnForwardMoves := PawnMoves[s]
	pawnCaptures := PawnCaptures[s]

	for pawns != 0 {
		sq1 := pawns.NextBit()
		pawns &= pawns - 1

		baseMove := Pawn<<14 | sq1<<7

		// forward moves
		forwardMoves := pawnForwardMoves[sq1]
		for forwardMoves != 0 {
			sq2 := forwardMoves.NextBit()
			forwardMoves &= forwardMoves - 1

			sq2Pos := Bits(1 << sq2)

			// check if blocked, including target square
			if b.All&(BitBetween[sq1][sq2]|sq2Pos) != 0 {
				continue
			}

			moves = append(moves, uint64(baseMove|sq2))
		}

		// captures
		captures := pawnCaptures[sq1]
		ep := b.EPTargetSquare
		for captures != 0 {
			sq2 := captures.NextBit()
			captures &= captures - 1

			sq2Pos := Bits(1 << sq2)
			if enemyUnits&sq2Pos == sq2Pos {
				moves = append(moves, uint64(baseMove|sq2))
			}

			// en-passant
			if ep > 0 && ep == sq2 {
				moves = append(moves, uint64(baseMove|sq2))
			}
		}
	}

	// add distinct promotion moves for each promotable piece type

	var promotions []uint64

	for i := 0; i < len(moves); i++ {
		move := moves[i]
		toIdx := move & 0x7F
		if toIdx >= H8 || toIdx <= A1 {
			promotions = append(promotions,
				Queen<<17|move,
				Rook<<17|move,
				Knight<<17|move,
				Bishop<<17|move,
			)
			moves = append(moves[:i], moves[i+1:]...)
		}
	}

	moves = append(moves, promotions...)

	return moves
}

func (b Board) pseudoLegalKnightMoves() []uint64 {
	s := b.ActiveColor
	knights := b.Pieces[s][Knight]
	if knights == 0 {
		return nil
	}

	moves := make([]uint64, 0, 16)

	friendlyUnits := b.Units[s]

	for knights != 0 {
		sq1 := knights.NextBit()
		knights &= knights - 1

		baseMove := Knight<<14 | sq1<<7

		knightMoves := PieceMoves[Knight][sq1]
		for knightMoves != 0 {
			sq2 := knightMoves.NextBit()
			knightMoves &= knightMoves - 1

			sq2Pos := Bits(1 << sq2)
			if friendlyUnits&sq2Pos == sq2Pos {
				continue
			}

			moves = append(moves, uint64(baseMove|sq2))
		}
	}

	return moves
}

func (b Board) pseudoLegalSliderMoves(pieceType int) []uint64 {
	s := b.ActiveColor
	sliders := b.Pieces[s][pieceType]
	if sliders == 0 {
		return nil
	}

	moves := make([]uint64, 0, 14)

	friendlyUnits := b.Units[s]

	sliderMoves := PieceMoves[pieceType]

	for sliders != 0 {
		sq1 := sliders.NextBit()
		sliders &= sliders - 1

		baseMove := pieceType<<14 | sq1<<7

		pieceMoves := sliderMoves[sq1]
		for pieceMoves != 0 {
			sq2 := pieceMoves.NextBit()
			pieceMoves &= pieceMoves - 1

			if b.IsBlocked(sq1, sq2) {
				continue
			}

			sq2Pos := Bits(1 << sq2)
			if friendlyUnits&sq2Pos == sq2Pos {
				continue
			}

			moves = append(moves, uint64(baseMove|sq2))
		}
	}

	return moves
}

func (b Board) pseudoLegalKingMoves() []uint64 {
	moves := make([]uint64, 0, 8)

	s := b.ActiveColor
	xs := 1 - s
	king := b.Pieces[s][King]
	kingSquare := king.NextBit()

	friendlyUnits := b.Units[s]

	baseMove := King<<14 | kingSquare<<7

	kingMoves := PieceMoves[King][kingSquare]
	for kingMoves != 0 {
		sq2 := kingMoves.NextBit()
		kingMoves &= kingMoves - 1

		sq2Pos := Bits(1 << sq2)
		if friendlyUnits&sq2Pos == sq2Pos {
			continue
		}

		moves = append(moves, uint64(baseMove|sq2))
	}

	// castling rights
	castle := (b.Castle >> (2 * s)) & 0b11

	// if they have castling rights and aren't in check
	if castle != 0 && !b.Attack(xs, kingSquare) {
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
				moves = append(moves, uint64(baseMove|kingDestSq))
			}
		}
	}

	return moves
}

func (b Board) filterPseudoLegalMoves(moves []uint64) []uint64 {
	// remove moves that leave the king in check
	s := b.ActiveColor
	xs := 1 - s

	for i := 0; i < len(moves); i++ {
		move := moves[i]
		b2 := b.apply(move)

		newKingSquare := b2.Pieces[s][King].NextBit()

		if b2.Attack(xs, newKingSquare) {
			moves = append(moves[:i], moves[i+1:]...)
			i--
			continue
		}
	}

	return moves
}

func (b Board) legalMoves() []uint64 {
	moves := make([]uint64, 0, 32)

	// pawns
	pawnMoves := b.pseudoLegalPawnMoves()
	moves = append(moves, pawnMoves...)

	// knights
	moves = append(moves, b.pseudoLegalKnightMoves()...)

	// bishops, rooks, queen
	for _, pieceType := range []int{Bishop, Rook, Queen} {
		moves = append(moves, b.pseudoLegalSliderMoves(pieceType)...)
	}

	// king
	moves = append(moves, b.pseudoLegalKingMoves()...)

	legalMoves := b.filterPseudoLegalMoves(moves)

	return legalMoves
}

func uciString(uci uint64) string {
	promo := (uci >> 17) & 0b111
	fromTo := uci & 0x3FFF

	if promo == 0 {
		return uciMoveStrings[fromTo]
	}

	return uciMoveStrings[fromTo] + uciMovePromo[promo]
}

func (b Board) LegalMoves() []string {
	moves := b.legalMoves()

	legalMoves := make([]string, len(moves))
	for i, uci := range moves {
		legalMoves[i] = uciString(uci)
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

func (b Board) apply(uciMove uint64) Board {
	// 000    000    0000000  0000000
	// 3bits  3bits  7bits    7bits
	// promo  type   from     to

	pieceType := (uciMove >> 14) & 0b111
	fromIdx := int((uciMove >> 7) & 0x7F)
	toIdx := int(uciMove & 0x7F)

	if pieceType > King || fromIdx > 63 || toIdx > 63 {
		panic(fmt.Errorf("invalid uciMove %017b. pieceType: %03b fromIdx: %d toIdx: %d", uciMove, pieceType, fromIdx, toIdx))
	}

	fromPos := Bits(1 << fromIdx)
	toPos := Bits(1 << toIdx)

	// set side (s) and opponent side (xs)
	s := b.ActiveColor
	xs := 1 - s
	epIdx := 0
	castle := (b.Castle >> (2 * s)) & 0b11

	// capture?
	resetHalfMoveClock := pieceType == Pawn
	if capturedPieceType := b.PieceType(toPos, xs); capturedPieceType != -1 {
		// remove captured piece
		b.Pieces[xs][capturedPieceType] &= ^toPos
		resetHalfMoveClock = true
	} else if pieceType == Pawn {
		if b.EPTargetSquare != 0 && b.EPTargetSquare == toIdx {
			// remove pawn captured by en passant
			var capturePos Bits
			if toIdx >= H6 {
				capturePos = toPos >> 8
			} else {
				capturePos = toPos << 8
			}
			epCapture := ^capturePos
			b.Pieces[xs][Pawn] &= epCapture
			resetHalfMoveClock = true
		} else if epFrom&fromPos == fromPos && epTo&toPos == toPos && epMask[toIdx]&b.Pieces[xs][Pawn] > 0 {
			epIdx = epTargetIndex[toIdx]
		}
	}

	if pieceType == King && castle != 0 {
		if PieceMoves[King][toIdx]&fromPos != fromPos {
			for sideIdx := ksIdx; sideIdx <= qsIdx; sideIdx++ {
				sideMask := uint8(sideIdx + 1)
				if castle&sideMask == 0 {
					continue
				}
				if castleKingTo[s][sideIdx]&toPos == 0 {
					continue
				}
				rookFrom := castleRookFrom[s][sideIdx]
				rookTo := castleRookTo[s][sideIdx]

				b.Pieces[s][Rook] &= ^rookFrom
				b.Pieces[s][Rook] |= rookTo
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
	b.Pieces[s][pieceType] &= ^fromPos

	newPieceType := pieceType
	if pieceType == Pawn {
		newPieceType = (uciMove >> 17) & 0b111
	}

	b.Pieces[s][newPieceType] |= toPos

	b.ActiveColor = xs

	castle = castle<<(s*2) | 0b11<<(xs*2)
	b.Castle &= castle

	b.EPTargetSquare = epIdx

	if resetHalfMoveClock {
		b.HalfMoveClock = 0
	} else {
		b.HalfMoveClock += 1
	}

	// increment full move number if it was black's turn
	b.FullMoveNumber += int(s)

	// update b.Units and b.All
	for i := 0; i < 2; i++ {
		b.Units[i] = b.Pieces[i][Pawn] |
			b.Pieces[i][Knight] |
			b.Pieces[i][Bishop] |
			b.Pieces[i][Rook] |
			b.Pieces[i][Queen] |
			b.Pieces[i][King]
	}

	b.All = b.Units[Black] | b.Units[White]

	return b
}

func (b Board) uciStringToMove(uci string) (uint64, error) {
	if len(uci) < 4 {
		return 0, xerrors.Errorf("invalid uci move '%s'", uci)
	}

	fromName := uci[0:2]
	toName := uci[2:4]

	var (
		fromIdx, toIdx int
		ok             bool
	)

	if fromIdx, ok = squareNameToIndex[fromName]; !ok {
		return 0, xerrors.Errorf("invalid uci move '%s'", uci)
	}
	if toIdx, ok = squareNameToIndex[toName]; !ok {
		return 0, xerrors.Errorf("invalid uci move '%s'", uci)
	}

	pieceType := b.PieceType(1<<fromIdx, b.ActiveColor)
	if pieceType == -1 {
		pieceType2 := b.PieceType(1<<fromIdx, b.ActiveColor)
		return 0, xerrors.Errorf("invalid uci move '%s', color '%s' does not have a piece on %s. other color piece type on same square: %d", uci, b.ActiveColor, fromName, pieceType2)
	}

	uciMove := uint64(pieceType<<14 | fromIdx<<7 | toIdx)

	// pawn promotion
	if pieceType == Pawn && (toIdx >= H8 || toIdx <= A1) {
		if len(uci) != 5 {
			return 0, xerrors.Errorf("invalid uci move '%s', pawn promotion type not specified", uci)
		}
		promotionType := uci[4]
		switch promotionType {
		case 'q':
			promotionType = Queen
		case 'r':
			promotionType = Rook
		case 'n':
			promotionType = Knight
		case 'b':
			promotionType = Bishop
		default:
			return 0, xerrors.Errorf("invalid uci move '%s', pawn promotion type '%c' not recognized. expected 'q', 'r', 'n', or 'b'", uci)
		}

		uciMove |= uint64(promotionType) << 17
	}

	return uciMove, nil
}

func (b Board) Apply(moves ...string) (Board, error) {
	bb := b

	for _, uci := range moves {
		uciMove, err := bb.uciStringToMove(uci)
		if err != nil {
			return bb, xerrors.Errorf("%w", err)
		}

		bb = bb.apply(uciMove)
	}

	return bb, nil
}

// PieceType returns which piece type is it a mask. -1 is returned if there's no piece or the piece is not of the color specified.
// TODO: this function would be more useful if it took an index (and used a lookup to find the mask), unless we're really testing a mask of piece bits.
// TODO: it looks like this function is minimally used, so it may be optimized away.
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

func (b Board) legalMovesByPieceType(pieceType int) ([]uint64, error) {
	var moves []uint64

	switch pieceType {
	case Queen, Rook, Bishop:
		moves = b.pseudoLegalSliderMoves(pieceType)
	case Knight:
		moves = b.pseudoLegalKnightMoves()
	case King:
		moves = b.pseudoLegalKingMoves()
	case Pawn:
		moves = b.pseudoLegalPawnMoves()
	default:
		return nil, xerrors.Errorf("invalid piece type '%d'", pieceType)
	}

	return b.filterPseudoLegalMoves(moves), nil
}

func (b Board) UCI(san string) (string, error) {
	originalSAN := san
	san = strings.TrimRight(originalSAN, "+#!?")

	if len(san) < 2 {
		return "", xerrors.Errorf("invalid SAN move '%s'", originalSAN)
	}

	var pieceType int

	switch pieceChar := san[0]; pieceChar {
	case 'Q':
		pieceType = Queen
	case 'R':
		pieceType = Rook
	case 'N':
		pieceType = Knight
	case 'B':
		pieceType = Bishop
	case 'K', 'O':
		pieceType = King
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
		pieceType = Pawn
	default:
		return "", xerrors.Errorf("invalid SAN move '%s'", originalSAN)
	}

	moves, err := b.legalMovesByPieceType(pieceType)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}

	if pieceType == Pawn {
		// check for promotion, save and trim it off
		var promotion rune
		if san[len(san)-2] == '=' {
			promotion = unicode.ToLower(rune(san[len(san)-1]))
			san = san[:len(san)-2]
		}

		toName := san[len(san)-2:]
		toIdx, ok := squareNameToIndex[toName]
		if !ok {
			return "", xerrors.Errorf("invalid SAN move '%s' (target square '%s' not found, san '%s', promotion '%c')", originalSAN, toName, san, promotion)
		}

		// when the SAN is e4, d5, etc. the target square is all we need
		if len(san) == 2 {
			for _, move := range moves {
				if toIdx != int(move&0x7F) {
					continue
				}

				fromIdx := int(move >> 7 & 0x7F)

				if promotion != 0 {
					return fmt.Sprintf("%s%s%c", squareNames[fromIdx], squareNames[toIdx], promotion), nil
				}
				return fmt.Sprintf("%s%s", squareNames[fromIdx], squareNames[toIdx]), nil
			}
			return "", xerrors.Errorf("FEN: '%s' SAN '%s' is not a legal move", b.FEN(), originalSAN)
		}

		// the other case is cxd4, exd4, etc. where we need the original column to identify the move
		fromColumn := int(7 - (san[0] - 'a'))
		for _, move := range moves {
			if toIdx != int(move&0x7F) {
				continue
			}

			fromIdx := int(move >> 7 & 0x7F)
			if fromIdx%8 != fromColumn {
				continue
			}

			if promotion != 0 {
				return fmt.Sprintf("%s%s%c", squareNames[fromIdx], squareNames[toIdx], promotion), nil
			}
			return fmt.Sprintf("%s%s", squareNames[fromIdx], squareNames[toIdx]), nil
		}
		return "", xerrors.Errorf("FEN: '%s' SAN '%s' is not a legal move", b.FEN(), originalSAN)
	}

	if len(san) < 3 {
		return "", xerrors.Errorf("invalid SAN move '%s'", originalSAN)
	}

	if pieceType == King && (san == "O-O" || san == "O-O-O") {
		expectedMove := castleSAN[b.ActiveColor][san]
		for _, move := range moves {
			if expectedMove != move {
				continue
			}

			fromIdx := int(move >> 7 & 0x7F)
			toIdx := int(move & 0x7F)

			return fmt.Sprintf("%s%s", squareNames[fromIdx], squareNames[toIdx]), nil
		}
		return "", xerrors.Errorf("FEN: '%s' SAN '%s' is not a legal move", b.FEN(), originalSAN)
	}

	toName := san[len(san)-2:]
	toIdx, ok := squareNameToIndex[toName]
	if !ok {
		return "", xerrors.Errorf("invalid SAN move '%s'", originalSAN)
	}

	// get the disambiguation string and trim off 'x'.
	disambiguation := strings.TrimRight(san[1:len(san)-2], "x")
	fromColumn, fromRow := -1, -1 // aka rank and file
	for _, c := range disambiguation {
		if c >= 'a' && c <= 'h' {
			if fromColumn != -1 {
				return "", xerrors.Errorf("invalid SAN move '%s'", originalSAN)
			}
			fromColumn = int(7 - (c - 'a'))
		} else if c >= '1' && c <= '8' {
			if fromRow != -1 {
				return "", xerrors.Errorf("invalid SAN move '%s'", originalSAN)
			}
			fromRow = int(c - '1')
		}
	}

	for _, move := range moves {
		if toIdx != int(move&0x7F) {
			continue
		}

		fromIdx := int(move >> 7 & 0x7F)

		if fromColumn != -1 && fromColumn != fromIdx%8 {
			continue
		}
		if fromRow != -1 && fromRow != fromIdx/8 {
			continue
		}

		// TODO: we could do more to validate the move, such as checking if other moves need a disambiguation, causing the SAN to be invalid.
		// TODO: any cases not handled currently are invalid SAN anyway.

		return fmt.Sprintf("%s%s", squareNames[fromIdx], squareNames[toIdx]), nil
	}

	return "", xerrors.Errorf("FEN: '%s' SAN '%s' is not a legal move", b.FEN(), originalSAN)
}

func (b Board) SAN(uci string) (string, error) {
	uciMove, err := b.uciStringToMove(uci)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}

	promo := int((uciMove >> 17) & 0b111)
	pieceType := int((uciMove >> 14) & 0b111)
	fromIdx := int((uciMove >> 7) & 0x7F)
	toIdx := int(uciMove & 0x7F)

	// get all legal moves for the piece
	allLegalPieceMoves, err := b.legalMovesByPieceType(pieceType)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}

	// filter by moves which are going to the destination square
	var moves []uint64
	for _, move := range allLegalPieceMoves {
		if toIdx != int(move&0x7F) {
			continue
		}
		moves = append(moves, move)
	}

	// find the move in the list of filtered moves
	var found bool
	for i := 0; i < len(moves); i++ {
		move := moves[i]
		if move == uciMove {
			moves = append(moves[:i], moves[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return "", xerrors.Errorf("UCI '%s' is not a legal move", uci)
	}

	// check if it's a capture
	var isCapture bool

	s := b.ActiveColor
	xs := 1 - s
	toPos := Bits(1 << toIdx)
	if b.Units[xs]&toPos == toPos {
		isCapture = true
	} else if pieceType == Pawn && b.EPTargetSquare == toIdx && b.EPTargetSquare != 0 {
		isCapture = true // en passant
	}

	var isCastling bool

	var sb strings.Builder
	switch pieceType {
	case King:
		castling := castleSAN[s]
		if castling["O-O"] == uciMove {
			isCastling = true
			sb.WriteString("O-O")
		} else if castling["O-O-O"] == uciMove {
			isCastling = true
			sb.WriteString("O-O-O")
		} else {
			sb.WriteByte('K')
		}
	case Queen:
		sb.WriteByte('Q')
	case Rook:
		sb.WriteByte('R')
	case Bishop:
		sb.WriteByte('B')
	case Knight:
		sb.WriteByte('N')
	}

	// check if move needs disambiguation
	if len(moves) != 0 && !isCastling {
		fromCol := fromIdx % 8
		fromRow := fromIdx / 8

		var sameCol, sameRow bool
		for _, move := range moves {
			moveFromIdx := int(move >> 7 & 0x7F)
			moveFromCol := moveFromIdx % 8
			moveFromRow := moveFromIdx / 8

			if moveFromCol == fromCol {
				sameCol = true
			}
			if moveFromRow == fromRow {
				sameRow = true
			}
		}

		if !sameCol {
			// if there's nothing on this column (file), use the file letter
			sb.WriteByte(byte('a' + (7 - fromCol)))
		} else if !sameRow {
			// if there's nothing on this row (rank), use the rank number
			sb.WriteByte(byte('1' + fromRow))
		} else {
			// otherwise, use the full square name
			sb.WriteString(squareNames[fromIdx])
		}
	}

	if isCapture {
		if pieceType == Pawn && sb.Len() == 0 {
			fromCol := fromIdx % 8
			sb.WriteByte(byte('a' + (7 - fromCol)))
		}
		sb.WriteByte('x')
	}

	if !isCastling {
		sb.WriteString(squareNames[toIdx])
	}

	if promo != 0 {
		sb.WriteByte('=')
		switch promo {
		case Queen:
			sb.WriteByte('Q')
		case Rook:
			sb.WriteByte('R')
		case Knight:
			sb.WriteByte('N')
		case Bishop:
			sb.WriteByte('B')
		}
	}

	// check and checkmate
	bb := b.apply(uciMove)
	if bb.IsCheck() {
		if len(bb.legalMoves()) == 0 {
			sb.WriteByte('#')
		} else {
			sb.WriteByte('+')
		}
	}

	return sb.String(), nil
}

func (b Board) IsCheck() bool {
	s := b.ActiveColor
	xs := 1 - s

	king := b.Pieces[s][King]
	kingSquare := king.NextBit()
	return b.Attack(xs, kingSquare)
}

func (b Board) IsCheckmate() bool {
	if !b.IsCheck() {
		return false
	}
	return len(b.legalMoves()) == 0
}

func (b Board) IsStalemate() bool {
	// TODO: draw by repetition, insufficient material
	if b.IsCheck() {
		return false
	}
	return len(b.legalMoves()) == 0
}
