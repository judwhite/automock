package bitboard

import "testing"

func TestParseFEN(t *testing.T) {
	type boardTest struct {
		name string
		bb   func(Board) Bits
		want string
	}

	cases := []struct {
		fen             string
		wantActiveColor int
		wantCastle      int

		boardTests []boardTest
	}{
		{
			fen:             StartPos,
			wantActiveColor: White,
			wantCastle:      0b1111,
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
