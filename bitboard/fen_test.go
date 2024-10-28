package bitboard

import (
	"testing"
)

func TestParseFEN(t *testing.T) {
	type boardTest struct {
		name string
		bb   func(Board) Bits
		want string
	}

	cases := []struct {
		fen                string
		wantActiveColor    Color
		wantCastle         int
		wantEPSquare       int
		wantHalfMoveClock  int
		wantFullMoveNumber int

		boardTests []boardTest
	}{
		{
			fen:                StartPos,
			wantActiveColor:    White,
			wantCastle:         0b1111,
			wantEPSquare:       0,
			wantHalfMoveClock:  0,
			wantFullMoveNumber: 1,

			boardTests: []boardTest{
				{
					name: "black pawns",
					bb:   func(b Board) Bits { return b.Pieces[Black][Pawn] },
					want: "" +
						"0 0 0 0 0 0 0 0\n" +
						"1 1 1 1 1 1 1 1\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n",
				},
				{
					name: "black knights",
					bb:   func(b Board) Bits { return b.Pieces[Black][Knight] },
					want: "" +
						"0 1 0 0 0 0 1 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n",
				},
				{
					name: "black bishops",
					bb:   func(b Board) Bits { return b.Pieces[Black][Bishop] },
					want: "" +
						"0 0 1 0 0 1 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n",
				},
				{
					name: "black rooks",
					bb:   func(b Board) Bits { return b.Pieces[Black][Rook] },
					want: "" +
						"1 0 0 0 0 0 0 1\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n",
				},
				{
					name: "black queens",
					bb:   func(b Board) Bits { return b.Pieces[Black][Queen] },
					want: "" +
						"0 0 0 1 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n",
				},
				{
					name: "black king",
					bb:   func(b Board) Bits { return b.Pieces[Black][King] },
					want: "" +
						"0 0 0 0 1 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n",
				},
				{
					name: "white pawns",
					bb:   func(b Board) Bits { return b.Pieces[White][Pawn] },
					want: "" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"1 1 1 1 1 1 1 1\n" +
						"0 0 0 0 0 0 0 0\n",
				},
				{
					name: "white knights",
					bb:   func(b Board) Bits { return b.Pieces[White][Knight] },
					want: "" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 1 0 0 0 0 1 0\n",
				},
				{
					name: "white bishops",
					bb:   func(b Board) Bits { return b.Pieces[White][Bishop] },
					want: "" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 1 0 0 1 0 0\n",
				},
				{
					name: "white rooks",
					bb:   func(b Board) Bits { return b.Pieces[White][Rook] },
					want: "" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"1 0 0 0 0 0 0 1\n",
				},
				{
					name: "white queens",
					bb:   func(b Board) Bits { return b.Pieces[White][Queen] },
					want: "" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 1 0 0 0 0\n",
				},
				{
					name: "white king",
					bb:   func(b Board) Bits { return b.Pieces[White][King] },
					want: "" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 1 0 0 0\n",
				},
				{
					name: "black units",
					bb:   func(b Board) Bits { return b.Units[Black] },
					want: "" +
						"1 1 1 1 1 1 1 1\n" +
						"1 1 1 1 1 1 1 1\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n",
				},
				{
					name: "white units",
					bb:   func(b Board) Bits { return b.Units[White] },
					want: "" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"1 1 1 1 1 1 1 1\n" +
						"1 1 1 1 1 1 1 1\n",
				},
				{
					name: "all",
					bb:   func(b Board) Bits { return b.All },
					want: "" +
						"1 1 1 1 1 1 1 1\n" +
						"1 1 1 1 1 1 1 1\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"0 0 0 0 0 0 0 0\n" +
						"1 1 1 1 1 1 1 1\n" +
						"1 1 1 1 1 1 1 1\n",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.fen, func(t *testing.T) {
			b, err := ParseFEN(c.fen)
			if err != nil {
				t.Error(err)
				return
			}

			if c.wantActiveColor != b.ActiveColor {
				t.Errorf("ActiveColor, want: %d, got: %d", c.wantActiveColor, b.ActiveColor)
			}

			if c.wantCastle != b.Castle {
				t.Errorf("Castle, want: %04b, got: %04b", c.wantCastle, b.Castle)
			}

			for _, bt := range c.boardTests {
				t.Run(bt.name, func(t *testing.T) {
					got := bt.bb(b).String()
					if bt.want != got {
						t.Errorf("want:\n%s\ngot:\n%s", got, bt.want)
					}
				})
			}
		})
	}
}

func TestBoard_FEN(t *testing.T) {
	cases := []struct {
		fen string
	}{
		{StartPos},
		{"r1bqkb1r/ppp2ppp/2n2n2/1B1pp3/4P3/P1N2N2/1PPP1PPP/R1BQK2R b KQkq - 1 5"},
		{"r1bqkb1r/ppp2ppp/2n2n2/1B2N3/4p3/P1N5/1PPP1PPP/R1BQK2R b KQkq - 0 6"},
		{"rnbqkb1r/1p2pppp/p2p1n2/8/3NP3/2N5/PPP2PPP/R1BQKB1R w KQkq - 0 6"},
		{"rnbqkb1r/ppp1pppp/5n2/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3"},
	}

	for _, c := range cases {
		t.Run(c.fen, func(t *testing.T) {
			b, err := ParseFEN(c.fen)
			if err != nil {
				t.Error(err)
				return
			}

			got := b.FEN()

			if c.fen != got {
				t.Errorf("\nwant: %v\ngot:  %v", c.fen, got)
			}
		})
	}
}
