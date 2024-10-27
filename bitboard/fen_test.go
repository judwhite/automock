package bitboard

import "testing"

func TestParseFEN(t *testing.T) {
	type fenTest struct {
		name string
		bb   func(Board) Bits
		want string
	}

	cases := []struct {
		fen   string
		tests []fenTest
	}{
		{
			fen: StartPos,
			tests: []fenTest{
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

			for _, c := range c.tests {
				t.Run(c.name, func(t *testing.T) {
					got := c.bb(b).String()
					if c.want != got {
						t.Errorf("want:\n%s\ngot:\n%s", got, c.want)
					}
				})
			}
		})
	}
}
