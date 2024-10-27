package bitboard

import (
	"fmt"
	"testing"
)

func TestBoard_IsBlocked(t *testing.T) {
	cases := []struct {
		fen  string
		sq1  int
		sq2  int
		want bool
	}{
		{
			fen:  "4n3/6k1/8/8/8/8/6K1/4R3 w - -",
			sq1:  E1,
			sq2:  E8,
			want: false,
		},
		{
			fen:  "4n3/6k1/8/8/4P3/8/6K1/4R3 w - -",
			sq1:  E1,
			sq2:  E8,
			want: true,
		},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s %s %s",
			c.fen,
			squareNames[c.sq1],
			squareNames[c.sq2])

		t.Run(name, func(t *testing.T) {
			b, err := ParseFEN(c.fen)
			if err != nil {
				t.Error(err)
				return
			}

			got := b.IsBlocked(c.sq1, c.sq2)
			if c.want != got {
				t.Errorf("want: %v, got: %v", c.want, got)
			}
		})
	}
}

func TestBoard_Attack(t *testing.T) {
	cases := []struct {
		fen  string
		want [2][64]int
	}{
		{
			fen: "rn1q1bnr/ppp1kBp1/3p3p/4N3/4P3/2N5/PPPP1PPP/R1BbK2R w KQ - 0 1",
			want: [2][64]int{
				{
					0, 0, 0, 0, 1, 0, 1, 0,
					0, 0, 0, 1, 0, 1, 0, 0,
					0, 0, 1, 0, 1, 0, 1, 0,
					0, 1, 0, 1, 0, 1, 0, 1,
					1, 0, 1, 0, 1, 0, 1, 0,
					1, 1, 1, 1, 1, 1, 1, 1,
					1, 1, 0, 1, 1, 1, 0, 1,
					0, 1, 1, 1, 1, 1, 1, 0,
				},
				{
					0, 1, 1, 1, 1, 1, 1, 0,
					1, 0, 1, 1, 1, 1, 1, 1,
					1, 1, 1, 1, 1, 1, 0, 1,
					0, 0, 1, 0, 1, 0, 1, 1,
					0, 0, 0, 0, 0, 0, 1, 0,
					0, 0, 0, 0, 0, 1, 0, 0,
					0, 0, 1, 0, 1, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
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

			for sq := 0; sq < 64; sq++ {
				wantWhite := c.want[White][63-sq]
				gotWhite := b.Attack(White, sq)
				if (wantWhite == 1) != gotWhite {
					t.Errorf("sq: %d side: white, want: %v, got: %v", sq, wantWhite, gotWhite)
				}

				wantBlack := c.want[Black][63-sq]
				gotBlack := b.Attack(Black, sq)
				if (wantBlack == 1) != gotBlack {
					t.Errorf("sq: %d side: black, want: %v, got: %v", sq, wantBlack, gotBlack)
				}
			}
		})

	}
}

// TODO

/*
   c3d5
   f7g8
   d2d4
   h1f1
   f7g6
   f7b3
   f7c4
   f7d5
   d2d3
   c3b5
   b2b4
   e5g6
   f7e6
   c3a4
   f2f3
   b2b3
   a2a4
   a1b1
   g2g3
   e1g1
   g2g4
   h2h3
   f2f4
   h2h4
   h1g1
   c3e2
   a2a3
   e5c6
   f7e8
   e5f3
   c3b1
   e1d1
   e5d7
   e5c4
   e5d3
   c3d1
   e1f1
   f7h5
   e5g4

*/
