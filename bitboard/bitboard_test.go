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

func TestBoard_LegalMoves(t *testing.T) {
	cases := []struct {
		fen  string
		want []string
	}{
		{
			fen: "1r3b2/1bR3rp/pNp2k2/2p1pP2/2P1P3/1P3N2/P5PP/3R2K1 w - - 8 26",
			want: []string{
				"b6d7",
				"d1d7",
				"c7d7",
				"c7g7",
				"c7b7",
				"f3e5",
				"b6a8",
				"g2g4",
				"g1f2",
				"g1f1",
				"f3g5",
				"h2h3",
				"h2h4",
				"d1d2",
				"g2g3",
				"d1d3",
				"a2a4",
				"g1h1",
				"b3b4",
				"b6d5",
				"c7f7",
				"c7c8",
				"d1d5",
				"c7e7",
				"b6a4",
				"a2a3",
				"f3e1",
				"d1d6",
				"c7c6",
				"b6c8",
				"f3h4",
				"d1d4",
				"d1d8",
				"d1f1",
				"d1b1",
				"f3d2",
				"d1a1",
				"f3d4",
				"d1c1",
				"d1e1",
			},
		},
		{
			fen: "1r3br1/1bR4p/pNp2k2/2p1pP2/2P1P3/1P3N2/P5PP/3R2K1 b - - 7 25",
			want: []string{
				"g8g7",
				"f8e7",
				"f8h6",
				"g8g4",
				"b7a8",
				"a6a5",
				"b8e8",
				"h7h5",
				"h7h6",
				"g8h8",
				"f8g7",
				"b8a8",
				"b8d8",
				"b8c8",
				"b7c8",
				"g8g2",
				"g8g6",
				"g8g5",
				"g8g3",
				"f8d6",
			},
		},
		{
			fen: "rnbqkb1r/ppp1pppp/5n2/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
			want: []string{
				"e5f6",
				"d2d4",
				"g1f3",
				"b1c3",
				"c2c3",
				"f2f4",
				"f1e2",
				"c2c4",
				"h2h3",
				"d2d3",
				"b1a3",
				"a2a4",
				"a2a3",
				"d1e2",
				"g1e2",
				"b2b3",
				"g2g3",
				"h2h4",
				"f1d3",
				"f1b5",
				"g1h3",
				"f2f3",
				"d1f3",
				"b2b4",
				"e5d6",
				"e5e6",
				"g2g4",
				"e1e2",
				"f1c4",
				"f1a6",
				"d1h5",
				"d1g4",
			},
		},
		{
			fen: "rnbqkb1r/ppp1pppp/3P1n2/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
			want: []string{
				"c7d6",
				"e7d6",
				"d8d6",
				"e7e5",
				"c7c5",
				"b8c6",
				"c8f5",
				"c7c6",
				"c8g4",
				"a7a6",
				"h7h5",
				"g7g6",
				"e7e6",
				"h7h6",
				"c8d7",
				"c8e6",
				"d8d7",
				"a7a5",
				"b8d7",
				"f6g4",
				"b7b6",
				"g7g5",
				"b8a6",
				"f6d5",
				"h8g8",
				"f6e4",
				"f6d7",
				"f6g8",
				"e8d7",
				"b7b5",
				"f6h5",
				"c8h3",
			},
		},
		{
			fen: "rnbqkb1r/pp2Pppp/2p2n2/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 4",
			want: []string{
				"f8e7",
				"d8e7",
				"e8e7",
				"d8d4",
				"d8d5",
				"d8d6",
				"d8c7",
				"d8b6",
				"d8d7",
				"d8a5",
				"f6d5",
				"c8g4",
				"b8d7",
				"b8a6",
				"c8f5",
				"f6g4",
				"c8e6",
				"f6e4",
				"c8d7",
				"d8d2",
				"b7b5",
				"b7b6",
				"h7h6",
				"e8d7",
				"h7h5",
				"a7a5",
				"c6c5",
				"g7g6",
				"a7a6",
				"h8g8",
				"f6h5",
				"g7g5",
				"f6d7",
				"d8d3",
				"f6g8",
				"c8h3",
			},
		},
		{
			fen: "rnb1kb1r/pp2Pppp/2p2n2/q7/8/8/PPPP1PPP/RNBQKBNR w KQkq - 1 5",
			want: []string{
				"e7f8q",
				"e7f8r",
				"e7f8n",
				"e7f8b",
				"b2b4",
				"f1a6",
				"g1e2",
				"g1f3",
				"d1f3",
				"f1e2",
				"f1c4",
				"b1c3",
				"f1b5",
				"a2a3",
				"f1d3",
				"g1h3",
				"c2c3",
				"h2h3",
				"a2a4",
				"h2h4",
				"d1e2",
				"b1a3",
				"c2c4",
				"g2g3",
				"b2b3",
				"f2f4",
				"f2f3",
				"g2g4",
				"e1e2",
				"d1h5",
				"d1g4",
			},
		},
		{
			fen: "rnb1kQ1r/pp3ppp/2p2n2/q7/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 5",
			want: []string{
				"e8f8",
				"h8f8",
				"e8d7",
			},
		},
		{
			fen: "r3kr2/pp3ppp/n1p2n2/q7/6b1/5N2/PPPPBPPP/RNBQK2R w KQq - 4 8",
			want: []string{
				"e1g1",
				"b1c3",
				"c2c3",
				"h2h3",
				"b2b4",
				"e1f1",
				"a2a3",
				"a2a4",
				"h2h4",
				"b2b3",
				"b1a3",
				"f3d4",
				"f3g1",
				"c2c4",
				"h1g1",
				"h1f1",
				"e2c4",
				"g2g3",
				"e2a6",
				"e2f1",
				"e2d3",
				"f3h4",
				"f3g5",
				"f3e5",
				"e2b5",
			},
		},
		{
			fen: "2kr1r2/pp3ppp/B1p2n2/q7/6b1/5N2/PPPP1PPP/RNBQ1RK1 b - - 0 9",
			want: []string{
				"a5a6",
				"a5h5",
				"f8e8",
				"a5f5",
				"b7a6",
				"g4f3",
				"f8h8",
				"f6e4",
				"g7g5",
				"a5c7",
				"f8g8",
				"d8d7",
				"g4h5",
				"d8d5",
				"a5c5",
				"d8d4",
				"d8e8",
				"a5b6",
				"h7h6",
				"a5d5",
				"g4h3",
				"h7h5",
				"d8d6",
				"g7g6",
				"c6c5",
				"f6g8",
				"f6h5",
				"a5b4",
				"g4f5",
				"a5a4",
				"f6d5",
				"g4e6",
				"f6d7",
				"f6e8",
				"c8b8",
				"g4d7",
				"c8c7",
				"a5b5",
				"a5g5",
				"a5a2",
				"a5e5",
				"a5a3",
				"d8d2",
				"c8d7",
				"a5d2",
				"a5c3",
				"d8d3",
			},
		},
		{
			// mated position
			fen:  "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3",
			want: nil,
		},
		{
			fen: "rn1qkbnr/p1pp1ppp/bp6/4p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4",
			want: []string{
				"f3e5",
				"d2d3",
				"d2d4",
				"b1c3",
				"c2c4",
				"b2b4",
				"a2a3",
				"a2a4",
				"h2h4",
				"b2b3",
				"h2h3",
				"b1a3",
				"g2g3",
				"f3d4",
				"c2c3",
				"f3g5",
				"h1g1",
				"f3g1",
				"g2g4",
				"h1f1",
				"f3h4",
				"d1e2",
			},
		},
		{
			fen: "rn1qkbnr/p1pp1ppp/bp6/4p3/2P1P3/5N2/PP1P1PPP/RNBQK2R b KQkq - 0 4",
			want: []string{
				"b8c6",
				"f7f6",
				"f8d6",
				"f8c5",
				"f8b4",
				"a6b7",
				"d8f6",
				"d8e7",
				"h7h5",
				"g8f6",
				"g8e7",
				"g8h6",
				"f8e7",
				"g7g6",
				"c7c5",
				"h7h6",
				"d7d6",
				"a6c4",
				"d8c8",
				"c7c6",
				"b6b5",
				"a6c8",
				"g7g5",
				"f7f5",
				"d7d5",
				"e8e7",
				"f8a3",
				"a6b5",
				"d8h4",
				"d8g5",
			},
		},
		{
			fen: "r2qk2r/p1ppbppp/1pn2n2/1b2p3/2P1P3/2NP1N2/PP1BQPPP/R3K2R b KQkq - 4 8",
			want: []string{
				"b5a6",
				"b5c4",
				"a7a6",
				"b5a4",
				"e7c5",
				"e8g8",
				"h7h6",
				"d7d6",
				"a8c8",
				"a7a5",
				"d8c8",
				"e7b4",
				"a8b8",
				"c6d4",
				"d8b8",
				"h7h5",
				"d7d5",
				"f6h5",
				"g7g6",
				"e7f8",
				"c6b4",
				"f6e4",
				"e7d6",
				"c6a5",
				"c6b8",
				"e7a3",
				"f6g8",
				"f6g4",
				"e8f8",
				"g7g5",
				"h8g8",
				"h8f8",
				"f6d5",
			},
		},
		{
			fen: "r2qk2r/p1p1bppp/1pnp1n2/4p3/b1P1P3/2NP1N2/1P1BQPPP/R3K2R w KQkq - 0 10",
			want: []string{
				"a1a4",
				"c3a4",
				"c3d5",
				"d2g5",
				"b2b4",
				"d2e3",
				"f3g5",
				"c3b5",
				"d2c1",
				"d3d4",
				"h2h3",
				"e1g1",
				"c4c5",
				"h2h4",
				"g2g3",
				"g2g4",
				"f3h4",
				"d2h6",
				"e2e3",
				"e2f1",
				"e1f1",
				"c3b1",
				"d2f4",
				"b2b3",
				"f3g1",
				"a1a3",
				"c3d1",
				"h1g1",
				"h1f1",
				"c3a2",
				"a1c1",
				"f3e5",
				"a1a2",
				"a1b1",
				"a1d1",
				"f3d4",
				"e2d1",
			},
		},
		{
			fen: "r2qk2r/p1p2ppp/1pnp4/4p1b1/b1P1P1n1/1PNP1N2/4QPPP/R3K2R w KQkq - 0 12",
			want: []string{
				"f3g5",
				"a1a4",
				"b3a4",
				"c3a4",
				"c3d5",
				"e1g1",
				"d3d4",
				"c4c5",
				"e2a2",
				"c3b5",
				"h2h4",
				"f3h4",
				"b3b4",
				"h2h3",
				"e2c2",
				"g2g3",
				"e2d1",
				"e2b2",
				"c3d1",
				"e2f1",
				"c3b1",
				"f3g1",
				"f3d2",
				"e1f1",
				"a1d1",
				"a1a3",
				"a1a2",
				"h1f1",
				"e1d1",
				"h1g1",
				"a1b1",
				"f3e5",
				"c3a2",
				"a1c1",
				"f3d4",
				"e2e3",
				"e2d2",
			},
		},
		{
			fen: "r2qk2r/p1p2ppp/1pnp4/4p1N1/b1P1P1n1/1PNP4/4QPPP/R3K2R b KQkq - 0 12",
			want: []string{
				"d8g5",
				"g4f6",
				"g4h6",
				"a4b3",
				"h7h5",
				"c6d4",
				"h7h6",
				"a7a5",
				"d8c8",
				"d8d7",
				"e8g8",
				"f7f5",
				"c6b4",
				"g4f2",
				"a7a6",
				"f7f6",
				"g7g6",
				"c6e7",
				"a8b8",
				"g4h2",
				"g4e3",
				"d8e7",
				"a4b5",
				"e8f8",
				"d8b8",
				"d8f6",
				"c6b8",
				"c6a5",
				"h8g8",
				"a8c8",
				"d6d5",
				"h8f8",
				"b6b5",
				"e8e7",
				"e8d7",
			},
		},
		{
			fen: "r2qk2r/p1p2ppp/1pnp4/4p1N1/b1P1P3/1PNP4/4QPPn/R3K2R w KQkq - 0 13",
			want: []string{
				"g5f7",
				"e2h5",
				"g5h7",
				"f2f4",
				"a1a4",
				"g5e6",
				"e2e3",
				"b3a4",
				"e2a2",
				"c3d5",
				"c4c5",
				"e2d2",
				"f2f3",
				"g5h3",
				"g2g3",
				"c3b5",
				"c3a4",
				"h1h2",
				"d3d4",
				"e2b2",
				"b3b4",
				"a1a2",
				"a1a3",
				"a1b1",
				"g5f3",
				"e2d1",
				"e2c2",
				"c3b1",
				"g2g4",
				"a1d1",
				"e1d1",
				"c3d1",
				"c3a2",
				"h1g1",
				"a1c1",
				"h1f1",
				"e1d2",
				"e2f1",
				"e1c1",
				"e2g4",
				"e2f3",
			},
		},
		{
			fen: "r2qk2r/p1p2ppp/1pnp4/3Np1N1/b1P1P3/1P1P1n2/4QPP1/R3K2R w KQkq - 2 14",
			want: []string{
				"e2f3",
				"g5f3",
				"g2f3",
				"e1f1",
				"e1d1",
			},
		},
		{
			fen: "r3k2r/p1pq1ppp/1pnp4/3Np1N1/R1P1P3/1P1P1Q2/5PP1/4K2R b Kkq - 0 15",
			want: []string{
				"c6d4",
				"h7h6",
				"f7f6",
				"a7a5",
				"c6d8",
				"h7h5",
				"h8f8",
				"e8c8",
				"a7a6",
				"a8b8",
				"a8c8",
				"a8d8",
				"c6a5",
				"f7f5",
				"h8g8",
				"c6b8",
				"e8g8",
				"c6b4",
				"e8f8",
				"b6b5",
				"d7e7",
				"e8d8",
				"g7g6",
				"d7c8",
				"c6e7",
				"d7f5",
				"d7e6",
				"d7g4",
				"d7h3",
				"d7d8",
			},
		},
		{
			fen: "r1q1k2r/p1p2ppp/1pnp4/3Np1N1/2P1P3/1P1P1Q2/5PP1/R3K2R b Kkq - 2 16",
			want: []string{
				"c8d7",
				"c6d8",
				"h8f8",
				"f7f6",
				"e8g8",
				"f7f5",
				"h7h6",
				"c8b7",
				"h7h5",
				"c6d4",
				"e8d8",
				"a7a5",
				"c6b4",
				"e8d7",
				"a7a6",
				"c6e7",
				"g7g6",
				"b6b5",
				"c8f5",
				"c8b8",
				"a8b8",
				"c6a5",
				"c8e6",
				"c8g4",
				"h8g8",
				"c8h3",
				"c8a6",
				"c6b8",
				"c8d8",
				"e8f8",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.fen, func(t *testing.T) {
			b, err := ParseFEN(c.fen)
			if err != nil {
				t.Fatal(err)
			}

			wantMap := make(map[string]struct{})
			for _, wantMove := range c.want {
				wantMap[wantMove] = struct{}{}
			}

			got := b.LegalMoves()

			gotMap := make(map[string]struct{})
			for _, gotMove := range got {
				gotMap[gotMove] = struct{}{}
			}

			for k := range gotMap {
				if _, ok := wantMap[k]; !ok {
					t.Errorf("'%s' is not a legal move.", k)
				}
			}

			for k := range wantMap {
				if _, ok := gotMap[k]; !ok {
					t.Errorf("'%s' is missing.", k)
				}
			}
		})
	}
}
