package bitboard

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

var fenCastlingAvailability = [16]string{
	" -",    // 0000
	" K",    // 0001
	" Q",    // 0010
	" KQ",   // 0011
	" k",    // 0100
	" Kk",   // 0101
	" Qk",   // 0110
	" KQk",  // 0111
	" q",    // 1000
	" Kq",   // 1001
	" Qq",   // 1010
	" KQq",  // 1011
	" kq",   // 1100
	" Kkq",  // 1101
	" Qkq",  // 1110
	" KQkq", // 1111
}

var benchCastlingAvailabilityMethod int

func ParseFEN(fen string) (Board, error) {
	// [0]: piece placement
	// [1]: active color
	// [2]: castling availability
	// [3]: ep target square
	// [4]: halfmove clock
	// [5]: fullmove number

	if fen == "startpos" {
		return ParseFEN(StartPos)
	}

	fenParts := strings.Split(fen, " ")
	if len(fenParts) != 4 && len(fenParts) != 6 {
		return Board{}, xerrors.Errorf("invalid FEN '%s', execpted 4 or 6 parts, got %d", fen, len(fenParts))
	}

	// [0]: piece placement

	ranks := strings.Split(fenParts[0], "/")
	if len(ranks) != 8 {
		return Board{}, xerrors.Errorf("invalid FEN '%s', expected 8 ranks, got %d", fen, len(ranks))
	}

	var b Board

	for i := 0; i < len(ranks); i++ {
		rank := (7 - i) * 8
		file := 7
		for j := 0; j < len(ranks[i]); j++ {
			if file < 0 {
				return Board{}, xerrors.Errorf("invalid FEN '%s' in rank index %d '%s'", fen, i, ranks[i])
			}

			c := ranks[i][j]
			if '1' <= c && c <= '8' {
				file -= int(c - '0')
				continue
			}

			bitPos := rank + file
			switch c {
			case 'p':
				b.Pieces[Black][Pawn] |= 1 << bitPos
			case 'n':
				b.Pieces[Black][Knight] |= 1 << bitPos
			case 'b':
				b.Pieces[Black][Bishop] |= 1 << bitPos
			case 'r':
				b.Pieces[Black][Rook] |= 1 << bitPos
			case 'q':
				b.Pieces[Black][Queen] |= 1 << bitPos
			case 'k':
				b.Pieces[Black][King] |= 1 << bitPos
			case 'P':
				b.Pieces[White][Pawn] |= 1 << bitPos
			case 'N':
				b.Pieces[White][Knight] |= 1 << bitPos
			case 'B':
				b.Pieces[White][Bishop] |= 1 << bitPos
			case 'R':
				b.Pieces[White][Rook] |= 1 << bitPos
			case 'Q':
				b.Pieces[White][Queen] |= 1 << bitPos
			case 'K':
				b.Pieces[White][King] |= 1 << bitPos
			default:
				return Board{}, xerrors.Errorf("invalid FEN '%s' in rank index %d '%s' character '%c'", fen, i, ranks[i], c)
			}

			file--
		}

		if file != -1 {
			return Board{}, xerrors.Errorf("invalid FEN '%s' in rank index %d '%s'", fen, i, ranks[i])
		}
	}

	for i := 0; i < 2; i++ {
		b.Units[i] = b.Pieces[i][Pawn] |
			b.Pieces[i][Knight] |
			b.Pieces[i][Bishop] |
			b.Pieces[i][Rook] |
			b.Pieces[i][Queen] |
			b.Pieces[i][King]
	}

	b.All = b.Units[Black] | b.Units[White]

	// [1]: active color

	activeColor := fenParts[1]
	switch activeColor {
	case "w":
		b.ActiveColor = White
	case "b":
		b.ActiveColor = Black
	default:
		return Board{}, xerrors.Errorf("invalid FEN '%s', active color expected to be 'w' or 'b', got '%s'", fen, activeColor)
	}

	// [2]: castling availability

	castling := fenParts[2]
	for _, c := range castling {
		switch c {
		case 'K':
			b.Castle |= 1
		case 'Q':
			b.Castle |= 2
		case 'k':
			b.Castle |= 4
		case 'q':
			b.Castle |= 8
		case '-':
			// no-op
		default:
			return Board{}, xerrors.Errorf("invalid FEN '%s', castling availability '%s' has unexpected character '%c'", fen, castling, c)
		}
	}

	// [3]: ep target square
	epTargetSquare := fenParts[3]
	if epTargetSquare != "-" {
		idx, ok := squareNameToIndex[epTargetSquare]
		if !ok {
			return Board{}, xerrors.Errorf("invalid FEN '%s', ep target square '%s' name invalid", fen, epTargetSquare)
		}
		b.EPTargetSquare = idx
	}

	// check for short 'fen key' version; early exit
	if len(fenParts) == 4 {
		b.HalfMoveClock = 0
		b.FullMoveNumber = 1
		return b, nil
	}

	// [4]: halfmove clock
	hmc := fenParts[4]
	halfMoveClock, err := strconv.Atoi(hmc)
	if err != nil {
		return Board{}, xerrors.Errorf("invalid fen '%s', halfmove clock '%s' is not an int", fen, hmc)
	}
	b.HalfMoveClock = halfMoveClock

	// [5]: fullmove number
	fmn := fenParts[5]
	fullMoveNumber, err := strconv.Atoi(fmn)
	if err != nil {
		return Board{}, xerrors.Errorf("invalid fen '%s', fullmove number '%s' is not an int", fen, fmn)
	}
	b.FullMoveNumber = fullMoveNumber

	return b, nil
}

func (b Board) FEN() string {
	return b.makeFEN(false)
}

func (b Board) FENKey() string {
	return b.makeFEN(true)
}

func (b Board) makeFEN(keyOnly bool) string {
	var sb strings.Builder

	// [0]: piece placement
	pos := Bits(1 << 63)
	for row := 7; row >= 0; row-- {
		var blankCount byte

		for col := 7; col >= 0; col-- {
			if b.All&pos != pos {
				blankCount++
				pos >>= 1
				continue
			}

			if blankCount > 0 {
				sb.WriteByte('0' + blankCount)
				blankCount = 0
			}

			if b.Units[White]&pos == pos {
				pieces := b.Pieces[White]
				if pieces[Pawn]&pos == pos {
					sb.WriteByte('P')
				} else if pieces[Knight]&pos == pos {
					sb.WriteByte('N')
				} else if pieces[Bishop]&pos == pos {
					sb.WriteByte('B')
				} else if pieces[Rook]&pos == pos {
					sb.WriteByte('R')
				} else if pieces[Queen]&pos == pos {
					sb.WriteByte('Q')
				} else if pieces[King]&pos == pos {
					sb.WriteByte('K')
				} else {
					idx := row*8 + col
					sq := squareNames[idx]
					panic(fmt.Errorf("inconsistent internal board state. b.All&pos == pos && b.Units[White]&pos == pos, but couldn't find pos in piece types of b.Pieces[White]. idx: %d sq: %s pos: %016X", idx, sq, pos))
				}
			} else {
				pieces := b.Pieces[Black]
				if pieces[Pawn]&pos == pos {
					sb.WriteByte('p')
				} else if pieces[Knight]&pos == pos {
					sb.WriteByte('n')
				} else if pieces[Bishop]&pos == pos {
					sb.WriteByte('b')
				} else if pieces[Rook]&pos == pos {
					sb.WriteByte('r')
				} else if pieces[Queen]&pos == pos {
					sb.WriteByte('q')
				} else if pieces[King]&pos == pos {
					sb.WriteByte('k')
				} else {
					idx := row*8 + col
					sq := squareNames[idx]
					panic(fmt.Errorf("inconsistent internal board state. b.All&pos == pos && b.Units[White]&pos != pos, but couldn't find pos in piece types of b.Pieces[Black]. idx: %d sq: %s pos: %016X", idx, sq, pos))
				}
			}

			pos >>= 1
		}

		if blankCount > 0 {
			sb.WriteByte('0' + blankCount)
		}

		if row != 0 {
			sb.WriteByte('/')
		}
	}

	// [1]: active color
	sb.WriteByte(' ')
	if b.ActiveColor == White {
		sb.WriteByte('w')
	} else {
		sb.WriteByte('b')
	}

	// [2]: castling availability
	if benchCastlingAvailabilityMethod == 0 {
		sb.WriteString(fenCastlingAvailability[b.Castle])
	} else {
		sb.WriteByte(' ')
		if b.Castle == 0 {
			sb.WriteByte('-')
		} else {
			if b.Castle&1 == 1 {
				sb.WriteByte('K')
			}
			if b.Castle&2 == 2 {
				sb.WriteByte('Q')
			}
			if b.Castle&4 == 4 {
				sb.WriteByte('k')
			}
			if b.Castle&8 == 8 {
				sb.WriteByte('q')
			}
		}
	}

	// [3]: ep target square
	sb.WriteByte(' ')
	if b.EPTargetSquare == 0 {
		sb.WriteByte('-')
	} else {
		sq := squareNames[b.EPTargetSquare]
		sb.WriteString(sq)
	}

	if keyOnly {
		return sb.String()
	}

	// [4]: halfmove clock
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(b.HalfMoveClock))

	// [5]: fullmove number
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(b.FullMoveNumber))

	return sb.String()
}
