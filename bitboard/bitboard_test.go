package bitboard

import (
	"fmt"
	"strings"
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

func TestBoard_Apply(t *testing.T) {
	cases := []struct {
		fen   string
		moves []string
		want  string
	}{
		{
			fen:   StartPos,
			moves: []string{},
			want:  StartPos,
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4"},
			want:  "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6"},
			want:  "rnbqkb1r/pppppppp/5n2/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 1 2",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5"},
			want:  "rnbqkb1r/pppppppp/5n2/4P3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5"},
			want:  "rnbqkb1r/ppp1pppp/5n2/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6"},
			want:  "rnbqkb1r/ppp1pppp/3P1n2/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6"},
			want:  "rnbqkb1r/pp2pppp/2pP1n2/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 4",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7"},
			want:  "rnbqkb1r/pp2Pppp/2p2n2/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 4",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5"},
			want:  "rnb1kb1r/pp2Pppp/2p2n2/q7/8/8/PPPP1PPP/RNBQKBNR w KQkq - 1 5",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5", "e7f8q"},
			want:  "rnb1kQ1r/pp3ppp/2p2n2/q7/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 5",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5", "e7f8q", "h8f8"},
			want:  "rnb1kr2/pp3ppp/2p2n2/q7/8/8/PPPP1PPP/RNBQKBNR w KQq - 0 6",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5", "e7f8q", "h8f8", "g1f3"},
			want:  "rnb1kr2/pp3ppp/2p2n2/q7/8/5N2/PPPP1PPP/RNBQKB1R b KQq - 1 6",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5", "e7f8q", "h8f8", "g1f3", "c8g4"},
			want:  "rn2kr2/pp3ppp/2p2n2/q7/6b1/5N2/PPPP1PPP/RNBQKB1R w KQq - 2 7",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5", "e7f8q", "h8f8", "g1f3", "c8g4", "f1e2"},
			want:  "rn2kr2/pp3ppp/2p2n2/q7/6b1/5N2/PPPPBPPP/RNBQK2R b KQq - 3 7",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5", "e7f8q", "h8f8", "g1f3", "c8g4", "f1e2", "b8a6"},
			want:  "r3kr2/pp3ppp/n1p2n2/q7/6b1/5N2/PPPPBPPP/RNBQK2R w KQq - 4 8",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5", "e7f8q", "h8f8", "g1f3", "c8g4", "f1e2", "b8a6", "e1g1"},
			want:  "r3kr2/pp3ppp/n1p2n2/q7/6b1/5N2/PPPPBPPP/RNBQ1RK1 b q - 5 8",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5", "e7f8q", "h8f8", "g1f3", "c8g4", "f1e2", "b8a6", "e1g1", "e8c8"},
			want:  "2kr1r2/pp3ppp/n1p2n2/q7/6b1/5N2/PPPPBPPP/RNBQ1RK1 w - - 6 9",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5", "e7f8q", "h8f8", "g1f3", "c8g4", "f1e2", "b8a6", "e1g1", "e8c8", "e2a6"},
			want:  "2kr1r2/pp3ppp/B1p2n2/q7/6b1/5N2/PPPP1PPP/RNBQ1RK1 b - - 0 9",
		},
		{
			fen:   StartPos,
			moves: []string{"e2e4", "g8f6", "e4e5", "d7d5", "e5d6", "c7c6", "d6e7", "d8a5", "e7f8q", "h8f8", "g1f3", "c8g4", "f1e2", "b8a6", "e1g1", "e8c8", "e2a6", "b7a6"},
			want:  "2kr1r2/p4ppp/p1p2n2/q7/6b1/5N2/PPPP1PPP/RNBQ1RK1 w - - 0 10",
		},
		{
			fen:   "rnbqkb1r/ppp1pppp/5n2/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
			moves: []string{"e5d6"},
			want:  "rnbqkb1r/ppp1pppp/3P1n2/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
		},

		// position fen startpos moves e2e4 d7d5 d2d4
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s %s", c.fen, strings.Join(c.moves, " "))
		t.Run(name, func(t *testing.T) {
			b, err := ParseFEN(c.fen)
			if err != nil {
				t.Error(err)
				return
			}

			b2, err := b.Apply(c.moves)
			if err != nil {
				t.Error(err)
				return
			}

			got := b2.FEN()

			if c.want != got {
				t.Errorf("\nwant: %v\ngot:  %v", c.want, got)
			}
		})
	}
}
