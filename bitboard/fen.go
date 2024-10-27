package bitboard

import (
	"strings"

	"golang.org/x/xerrors"
)

func ParseFEN(fen string) (Board, error) {
	// [0]: piece placement
	// [1]: active color
	// [2]: castling availability
	// [3]: ep target square
	// [4]: halfmove clock
	// [5]: fullmove number
	fenParts := strings.Split(fen, " ")
	if len(fenParts) != 4 && len(fenParts) != 6 {
		return Board{}, xerrors.Errorf("invalid FEN '%s', execpted 4 or 6 parts, got %d", fen, len(fenParts))
	}

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

	return b, nil
}
