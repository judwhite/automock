package pgnparse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"automock/bitboard"
)

func TestParseSingleGames(t *testing.T) {
	// arrange
	cases := []struct {
		input string
		want  *Game
	}{
		{
			input: "1. e4 e5 2. Nf3 Nc6",
			want: &Game{
				Result: ResultUnknown,
				Moves: []*Move{
					{SAN: "e4", Ply: 1},
					{SAN: "e5", Ply: 2},
					{SAN: "Nf3", Ply: 3},
					{SAN: "Nc6", Ply: 4},
				},
			},
		},
		{
			input: "1. e4 e5 (1...c5 2. Nf3 Nc6) 2. f4 exf4 *",
			want: &Game{
				Result: ResultUnknown,
				Moves: []*Move{
					{SAN: "e4", Ply: 1},
					{SAN: "e5", Variations: []*Variation{
						{
							Moves: []*Move{
								{SAN: "c5", Ply: 2},
								{SAN: "Nf3", Ply: 3},
								{SAN: "Nc6", Ply: 4},
							},
						},
					}, Ply: 2},
					{SAN: "f4", Ply: 3},
					{SAN: "exf4", Ply: 4},
				},
			},
		},
		{
			input: "1. f3 $2 e5 $1 2. g4 $4 $15 Qh4#\n0-1",
			want: &Game{
				Result: ResultBlackWins,
				Moves: []*Move{
					{SAN: "f3", NAGs: []string{"$2"}, Ply: 1},
					{SAN: "e5", NAGs: []string{"$1"}, Ply: 2},
					{SAN: "g4", NAGs: []string{"$4", "$15"}, Ply: 3},
					{SAN: "Qh4#", Ply: 4},
				},
			},
		},
		{
			input: "1. d4 {some commentary then} 1. ...d5 2. c4 dxc4 *",
			want: &Game{
				Result: ResultUnknown,
				Moves: []*Move{
					{Ply: 1, SAN: "d4", Comment: "some commentary then"},
					{Ply: 2, SAN: "d5"},
					{Ply: 3, SAN: "c4"},
					{Ply: 4, SAN: "dxc4"},
				},
			},
		},
		{
			input: "1. e4 (1. d4 d5 ) e5 2. d4 (2. Nf3 Nc6 ) exd4 *",
			want: &Game{
				Result: ResultUnknown,
				Moves: []*Move{
					{Ply: 1, SAN: "e4", Variations: []*Variation{
						{Moves: []*Move{
							{Ply: 1, SAN: "d4"},
							{Ply: 2, SAN: "d5"},
						}},
					}},
					{Ply: 2, SAN: "e5"},
					{Ply: 3, SAN: "d4", Variations: []*Variation{
						{Moves: []*Move{
							{Ply: 3, SAN: "Nf3"},
							{Ply: 4, SAN: "Nc6"},
						}},
					}},
					{Ply: 4, SAN: "exd4"},
				},
			},
		},
		{
			input: "1. e4 ({c1} {c2} 1. d4 d5) e5 2. d4 exd4 *",
			want: &Game{
				Result: ResultUnknown,
				Moves: []*Move{
					{Ply: 1, SAN: "e4", Variations: []*Variation{
						{Comments: []string{"c1", "c2"}, Moves: []*Move{
							{Ply: 1, SAN: "d4"},
							{Ply: 2, SAN: "d5"},
						}},
					}},
					{Ply: 2, SAN: "e5"},
					{Ply: 3, SAN: "d4"},
					{Ply: 4, SAN: "exd4"},
				},
			},
		},
		{
			input: "1. e4 (1. d4 d5) (1. c4 e5) e5 *",
			want: &Game{
				Result: ResultUnknown,
				Moves: []*Move{
					{Ply: 1, SAN: "e4", Variations: []*Variation{
						{Moves: []*Move{
							{Ply: 1, SAN: "d4"},
							{Ply: 2, SAN: "d5"},
						}},
						{Moves: []*Move{
							{Ply: 1, SAN: "c4"},
							{Ply: 2, SAN: "e5"},
						}},
					}},
					{Ply: 2, SAN: "e5"},
				},
			},
		},
		{
			input: "1. e4 (1. d4 d5)(1. f4 e5)e5 2. d4 exd4 *",
			want: &Game{
				Result: ResultUnknown,
				Moves: []*Move{
					{Ply: 1, SAN: "e4", Variations: []*Variation{
						{Moves: []*Move{
							{Ply: 1, SAN: "d4"},
							{Ply: 2, SAN: "d5"},
						}},
						{Moves: []*Move{
							{Ply: 1, SAN: "f4"},
							{Ply: 2, SAN: "e5"},
						}},
					}},
					{Ply: 2, SAN: "e5"},
					{Ply: 3, SAN: "d4"},
					{Ply: 4, SAN: "exd4"},
				},
			},
		},
		{
			input: "1. e4 (1. d4 d5)e5 2. d4 exd4 *",
			want: &Game{
				Result: ResultUnknown,
				Moves: []*Move{
					{Ply: 1, SAN: "e4", Variations: []*Variation{
						{Moves: []*Move{
							{Ply: 1, SAN: "d4"},
							{Ply: 2, SAN: "d5"},
						}},
					}},
					{Ply: 2, SAN: "e5"},
					{Ply: 3, SAN: "d4"},
					{Ply: 4, SAN: "exd4"},
				},
			},
		},
		{
			input: "1. e4 (1. d4 ) e5 2. d4 exd4 *",
			want: &Game{
				Result: ResultUnknown,
				Moves: []*Move{
					{Ply: 1, SAN: "e4", Variations: []*Variation{
						{Moves: []*Move{
							{Ply: 1, SAN: "d4"},
						}},
					}},
					{Ply: 2, SAN: "e5"},
					{Ply: 3, SAN: "d4"},
					{Ply: 4, SAN: "exd4"},
				},
			},
		},
		{
			/*
				// 1    good move (traditional "!")
				// 2    poor move (traditional "?")
				// 3    very good move (traditional "!!")
				// 4    very poor move (traditional "??")
				// 5    speculative move (traditional "!?")
				// 6    questionable move (traditional "?!")
			*/
			input: "1. e4? e5! 2. Nf3?? Nc6?! 3. Bc4!? Nf6!!\n0-1",
			want: &Game{
				Result: ResultBlackWins,
				Moves: []*Move{
					{Ply: 1, SAN: "e4", NAGs: []string{"$2"}},
					{Ply: 2, SAN: "e5", NAGs: []string{"$1"}},
					{Ply: 3, SAN: "Nf3", NAGs: []string{"$4"}},
					{Ply: 4, SAN: "Nc6", NAGs: []string{"$6"}},
					{Ply: 5, SAN: "Bc4", NAGs: []string{"$5"}},
					{Ply: 6, SAN: "Nf6", NAGs: []string{"$3"}},
				},
			},
		},
		// TODO: this conversion check was slowing down parsing; do something else
		//{
		//	input: "1. e4 f5 2. exf5 g6 3. fxg6 h6 4. g7 Rh7 5. gxf8Q+ *",
		//	want: &Game{
		//		Result: ResultUnknown,
		//		Moves: []*Move{
		//			{Ply: 1, SAN: "e4"},
		//			{Ply: 2, SAN: "f5"},
		//			{Ply: 3, SAN: "exf5"},
		//			{Ply: 4, SAN: "g6"},
		//			{Ply: 5, SAN: "fxg6"},
		//			{Ply: 6, SAN: "h6"},
		//			{Ply: 7, SAN: "g7"},
		//			{Ply: 8, SAN: "Rh7"},
		//			{Ply: 9, SAN: "gxf8=Q+"},
		//		},
		//	},
		//},
		{
			input: "{start of game} 1. f3 e5 2. g4 Qh4#\n0-1",
			want: &Game{
				Comment: "start of game",
				Result:  ResultBlackWins,
				Moves: []*Move{
					{Ply: 1, SAN: "f3"},
					{Ply: 2, SAN: "e5"},
					{Ply: 3, SAN: "g4"},
					{Ply: 4, SAN: "Qh4#"},
				},
			},
		},
		{
			input: "1. e4 {[%clk 0:15:09.9]} 1... Nc6 {[%clk 0:15:06.6]} 2. Nf3 {[%clk 0:15:15.1]} 2... e5 {[%clk 0:15:02.9]} 3. Bc4 {[%clk 0:15:16.9]} 3... h6 {[%clk 0:14:39.1]} 4. d4 {[%clk 0:15:23.1]} 4... exd4 {[%clk 0:14:30.7]} 5. Nxd4 {[%clk 0:15:29.1]} 5... Bc5 {[%clk 0:14:19]} 6. c3 {[%clk 0:15:06.4]} 6... Qe7 {[%clk 0:14:17.3]} 7. Qf3 {[%clk 0:14:48]} 7... Nf6 {[%clk 0:14:22.9]} 8. O-O {[%clk 0:14:49.9]} 8... Qxe4 {[%clk 0:14:00.9]} 9. Qxe4+ {[%clk 0:14:51]} 9... Nxe4 {[%clk 0:14:07.7]} 10. Re1 {[%clk 0:15:00.2]} 10... f5 {[%clk 0:13:48.3]} 11. f3 {[%clk 0:14:47.9]} 11... Bxd4+ {[%clk 0:12:38.1]} 12. cxd4 {[%clk 0:14:57.8]} 12... Nxd4 {[%clk 0:12:47]} 1-0\n",
			// 1. e4 Nc6 2. Nf3 e5 3. Bc4 h6
			// 4. d4 exd4 5. Nxd4 Bc5 6. c3 Qe7 7. Qf3 Nf6 8. O-O Qxe4 9. Qxe4+ Nxe4 10. Re1 f5 11. f3 Bxd4+ 12. cxd4 Nxd4 *
			want: &Game{
				Result: ResultWhiteWins,
				Moves: []*Move{
					{Ply: 1, SAN: "e4", Comment: "[%clk 0:15:09.9]"},
					{Ply: 2, SAN: "Nc6", Comment: "[%clk 0:15:06.6]"},
					{Ply: 3, SAN: "Nf3", Comment: "[%clk 0:15:15.1]"},
					{Ply: 4, SAN: "e5", Comment: "[%clk 0:15:02.9]"},
					{Ply: 5, SAN: "Bc4", Comment: "[%clk 0:15:16.9]"},
					{Ply: 6, SAN: "h6", Comment: "[%clk 0:14:39.1]"},
					{Ply: 7, SAN: "d4", Comment: "[%clk 0:15:23.1]"},
					{Ply: 8, SAN: "exd4", Comment: "[%clk 0:14:30.7]"},
					{Ply: 9, SAN: "Nxd4", Comment: "[%clk 0:15:29.1]"},
					{Ply: 10, SAN: "Bc5", Comment: "[%clk 0:14:19]"},
					{Ply: 11, SAN: "c3", Comment: "[%clk 0:15:06.4]"},
					{Ply: 12, SAN: "Qe7", Comment: "[%clk 0:14:17.3]"},
					{Ply: 13, SAN: "Qf3", Comment: "[%clk 0:14:48]"},
					{Ply: 14, SAN: "Nf6", Comment: "[%clk 0:14:22.9]"},
					{Ply: 15, SAN: "O-O", Comment: "[%clk 0:14:49.9]"},
					{Ply: 16, SAN: "Qxe4", Comment: "[%clk 0:14:00.9]"},
					{Ply: 17, SAN: "Qxe4+", Comment: "[%clk 0:14:51]"},
					{Ply: 18, SAN: "Nxe4", Comment: "[%clk 0:14:07.7]"},
					{Ply: 19, SAN: "Re1", Comment: "[%clk 0:15:00.2]"},
					{Ply: 20, SAN: "f5", Comment: "[%clk 0:13:48.3]"},
					{Ply: 21, SAN: "f3", Comment: "[%clk 0:14:47.9]"},
					{Ply: 22, SAN: "Bxd4+", Comment: "[%clk 0:12:38.1]"},
					{Ply: 23, SAN: "cxd4", Comment: "[%clk 0:14:57.8]"},
					{Ply: 24, SAN: "Nxd4", Comment: "[%clk 0:12:47]"},
				},
			},
		},
		{
			input: `[Event "F/S Return Match"]
[Site "Belgrade, Serbia JUG"]
[Date "1992.11.04"]
[Round "29"]
[White "Fischer, Robert J."]
[Black "Spassky, Boris V."]
[Result "1/2-1/2"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3
O-O 9. h3 Nb8 10. d4 Nbd7 11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 15.
Nb1 h6 16. Bh4 c5 17. dxe5 Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 Nxd6 21.
Nc4 Nxc4 22. Bxc4 Nb6 23. Ne5 Rae8 24. Bxf7+ Rxf7 25. Nxf7 Rxe1+ 26. Qxe1 Kxf7
27. Qe3 Qg5 28. Qxg5 hxg5 29. b3 Ke6 30. a3 Kd6 31. axb4 cxb4 32. Ra5 Nd5 33.
f3 Bc8 34. Kf2 Bf5 35. Ra7 g6 36. Ra6+ Kc5 37. Ke1 Nf4 38. g3 Nxh3 39. Kd2 Kb5
40. Rd6 Kc5 41. Ra6 Nf2 42. g4 Bd3 43. Re6 1/2-1/2`,
			want: &Game{
				Result: ResultDraw,
				Tags: []Tag{
					{"Event", "F/S Return Match"},
					{"Site", "Belgrade, Serbia JUG"},
					{"Date", "1992.11.04"},
					{"Round", "29"},
					{"White", "Fischer, Robert J."},
					{"Black", "Spassky, Boris V."},
					{"Result", "1/2-1/2"},
				},
				Moves: []*Move{
					{Ply: 1, SAN: "e4"},
					{Ply: 2, SAN: "e5"},
					{Ply: 3, SAN: "Nf3"},
					{Ply: 4, SAN: "Nc6"},
					{Ply: 5, SAN: "Bb5"},
					{Ply: 6, SAN: "a6"},
					{Ply: 7, SAN: "Ba4"},
					{Ply: 8, SAN: "Nf6"},
					{Ply: 9, SAN: "O-O"},
					{Ply: 10, SAN: "Be7"},
					{Ply: 11, SAN: "Re1"},
					{Ply: 12, SAN: "b5"},
					{Ply: 13, SAN: "Bb3"},
					{Ply: 14, SAN: "d6"},
					{Ply: 15, SAN: "c3"},
					{Ply: 16, SAN: "O-O"},
					{Ply: 17, SAN: "h3"},
					{Ply: 18, SAN: "Nb8"},
					{Ply: 19, SAN: "d4"},
					{Ply: 20, SAN: "Nbd7"},
					{Ply: 21, SAN: "c4"},
					{Ply: 22, SAN: "c6"},
					{Ply: 23, SAN: "cxb5"},
					{Ply: 24, SAN: "axb5"},
					{Ply: 25, SAN: "Nc3"},
					{Ply: 26, SAN: "Bb7"},
					{Ply: 27, SAN: "Bg5"},
					{Ply: 28, SAN: "b4"},
					{Ply: 29, SAN: "Nb1"},
					{Ply: 30, SAN: "h6"},
					{Ply: 31, SAN: "Bh4"},
					{Ply: 32, SAN: "c5"},
					{Ply: 33, SAN: "dxe5"},
					{Ply: 34, SAN: "Nxe4"},
					{Ply: 35, SAN: "Bxe7"},
					{Ply: 36, SAN: "Qxe7"},
					{Ply: 37, SAN: "exd6"},
					{Ply: 38, SAN: "Qf6"},
					{Ply: 39, SAN: "Nbd2"},
					{Ply: 40, SAN: "Nxd6"},
					{Ply: 41, SAN: "Nc4"},
					{Ply: 42, SAN: "Nxc4"},
					{Ply: 43, SAN: "Bxc4"},
					{Ply: 44, SAN: "Nb6"},
					{Ply: 45, SAN: "Ne5"},
					{Ply: 46, SAN: "Rae8"},
					{Ply: 47, SAN: "Bxf7+"},
					{Ply: 48, SAN: "Rxf7"},
					{Ply: 49, SAN: "Nxf7"},
					{Ply: 50, SAN: "Rxe1+"},
					{Ply: 51, SAN: "Qxe1"},
					{Ply: 52, SAN: "Kxf7"},
					{Ply: 53, SAN: "Qe3"},
					{Ply: 54, SAN: "Qg5"},
					{Ply: 55, SAN: "Qxg5"},
					{Ply: 56, SAN: "hxg5"},
					{Ply: 57, SAN: "b3"},
					{Ply: 58, SAN: "Ke6"},
					{Ply: 59, SAN: "a3"},
					{Ply: 60, SAN: "Kd6"},
					{Ply: 61, SAN: "axb4"},
					{Ply: 62, SAN: "cxb4"},
					{Ply: 63, SAN: "Ra5"},
					{Ply: 64, SAN: "Nd5"},
					{Ply: 65, SAN: "f3"},
					{Ply: 66, SAN: "Bc8"},
					{Ply: 67, SAN: "Kf2"},
					{Ply: 68, SAN: "Bf5"},
					{Ply: 69, SAN: "Ra7"},
					{Ply: 70, SAN: "g6"},
					{Ply: 71, SAN: "Ra6+"},
					{Ply: 72, SAN: "Kc5"},
					{Ply: 73, SAN: "Ke1"},
					{Ply: 74, SAN: "Nf4"},
					{Ply: 75, SAN: "g3"},
					{Ply: 76, SAN: "Nxh3"},
					{Ply: 77, SAN: "Kd2"},
					{Ply: 78, SAN: "Kb5"},
					{Ply: 79, SAN: "Rd6"},
					{Ply: 80, SAN: "Kc5"},
					{Ply: 81, SAN: "Ra6"},
					{Ply: 82, SAN: "Nf2"},
					{Ply: 83, SAN: "g4"},
					{Ply: 84, SAN: "Bd3"},
					{Ply: 85, SAN: "Re6"},
				},
			},
		},
		{
			input: `[Event "Ch World (match)"]
[Site "New York (USA)"]
[Date "1886.03.24"]
[EventDate "?"]
[Round "19"]
[Result "0-1"]
[White "Johannes Zukertort"]
[Black "Wilhelm Steinitz"]
[ECO "D53"]
[WhiteElo "?"]
[BlackElo "?"]
[PlyCount "58"]

1. d4 {Notes by Robert James Fischer from a television
interview. } d5 2. c4 e6 3. Nc3 Nf6 4. Bg5 Be7 5. Nf3 O-O
6. c5 {White plays a mistake already; he should just play e3,
naturally.--Fischer} b6 7. b4 bxc5 8. dxc5 a5 9. a3 {Now he
plays this fantastic move; it's the winning move. -- Fischer}
d4 {He can't take with the knight, because of axb4.--Fischer}
10. Bxf6 gxf6 11. Na4 e5 {This kingside weakness is nothing;
the center is easily winning.--Fischer} 12. b5 Be6 13. g3 c6
14. bxc6 Nxc6 15. Bg2 Rb8 {Threatening Bb3.--Fischer} 16. Qc1
d3 17. e3 e4 18. Nd2 f5 19. O-O Re8 {A very modern move; a
quiet positional move. The rook is doing nothing now, but
later...--Fischer} 20. f3 {To break up the center, it's his
only chance.--Fischer} Nd4 21. exd4 Qxd4+ 22. Kh1 e3 23. Nc3
Bf6 24. Ndb1 d2 25. Qc2 Bb3 26. Qxf5 d1=Q 27. Nxd1 Bxd1
28. Nc3 e2 29. Raxd1 Qxc3 0-1`,
			want: &Game{
				Tags: []Tag{
					{"Event", "Ch World (match)"},
					{"Site", "New York (USA)"},
					{"Date", "1886.03.24"},
					{"EventDate", "?"},
					{"Round", "19"},
					{"Result", "0-1"},
					{"White", "Johannes Zukertort"},
					{"Black", "Wilhelm Steinitz"},
					{"ECO", "D53"},
					{"WhiteElo", "?"},
					{"BlackElo", "?"},
					{"PlyCount", "58"},
				},
				Result: ResultBlackWins,
				Moves: []*Move{
					{Ply: 1, SAN: "d4", Comment: "Notes by Robert James Fischer from a television interview."},
					{Ply: 2, SAN: "d5"},
					{Ply: 3, SAN: "c4"},
					{Ply: 4, SAN: "e6"},
					{Ply: 5, SAN: "Nc3"},
					{Ply: 6, SAN: "Nf6"},
					{Ply: 7, SAN: "Bg5"},
					{Ply: 8, SAN: "Be7"},
					{Ply: 9, SAN: "Nf3"},
					{Ply: 10, SAN: "O-O"},
					{Ply: 11, SAN: "c5", Comment: "White plays a mistake already; he should just play e3, naturally.--Fischer"},
					{Ply: 12, SAN: "b6"},
					{Ply: 13, SAN: "b4"},
					{Ply: 14, SAN: "bxc5"},
					{Ply: 15, SAN: "dxc5"},
					{Ply: 16, SAN: "a5"},
					{Ply: 17, SAN: "a3", Comment: "Now he plays this fantastic move; it's the winning move. -- Fischer"},
					{Ply: 18, SAN: "d4", Comment: "He can't take with the knight, because of axb4.--Fischer"},
					{Ply: 19, SAN: "Bxf6"},
					{Ply: 20, SAN: "gxf6"},
					{Ply: 21, SAN: "Na4"},
					{Ply: 22, SAN: "e5", Comment: "This kingside weakness is nothing; the center is easily winning.--Fischer"},
					{Ply: 23, SAN: "b5"},
					{Ply: 24, SAN: "Be6"},
					{Ply: 25, SAN: "g3"},
					{Ply: 26, SAN: "c6"},
					{Ply: 27, SAN: "bxc6"},
					{Ply: 28, SAN: "Nxc6"},
					{Ply: 29, SAN: "Bg2"},
					{Ply: 30, SAN: "Rb8", Comment: "Threatening Bb3.--Fischer"},
					{Ply: 31, SAN: "Qc1"},
					{Ply: 32, SAN: "d3"},
					{Ply: 33, SAN: "e3"},
					{Ply: 34, SAN: "e4"},
					{Ply: 35, SAN: "Nd2"},
					{Ply: 36, SAN: "f5"},
					{Ply: 37, SAN: "O-O"},
					{Ply: 38, SAN: "Re8", Comment: "A very modern move; a quiet positional move. The rook is doing nothing now, but later...--Fischer"},
					{Ply: 39, SAN: "f3", Comment: "To break up the center, it's his only chance.--Fischer"},
					{Ply: 40, SAN: "Nd4"},
					{Ply: 41, SAN: "exd4"},
					{Ply: 42, SAN: "Qxd4+"},
					{Ply: 43, SAN: "Kh1"},
					{Ply: 44, SAN: "e3"},
					{Ply: 45, SAN: "Nc3"},
					{Ply: 46, SAN: "Bf6"},
					{Ply: 47, SAN: "Ndb1"},
					{Ply: 48, SAN: "d2"},
					{Ply: 49, SAN: "Qc2"},
					{Ply: 50, SAN: "Bb3"},
					{Ply: 51, SAN: "Qxf5"},
					{Ply: 52, SAN: "d1=Q"},
					{Ply: 53, SAN: "Nxd1"},
					{Ply: 54, SAN: "Bxd1"},
					{Ply: 55, SAN: "Nc3"},
					{Ply: 56, SAN: "e2"},
					{Ply: 57, SAN: "Raxd1"},
					{Ply: 58, SAN: "Qxc3"},
				},
			},
		},
		/*
			want: &Game{
				Result: "",
				Moves: []*Move{
					{Ply: 0, SAN: ""},
					{Ply: 0, SAN: ""},
					{Ply: 0, SAN: "", Variations: []*Variation{
						{Moves: []*Move{
							{Ply: 0, SAN: ""},
						}},
					}},
				},
			},
		*/
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			// act
			pgn, err := Parse(c.input)
			if err != nil {
				t.Error(err)
				return
			}

			// assert
			if pgn == nil {
				t.Error("pgn == nil")
				return
			}

			if len(pgn.Games) != 1 {
				t.Errorf("len(pgn.Games) want: 1 got: %d", len(pgn.Games))
				return
			}

			game := pgn.Games[0]

			if !c.want.Equals(game) {
				wantString := c.want.String()
				gotString := game.String()
				t.Errorf("\nwant:\n%s\n\ngot:\n%s\n", wantString, gotString)
			}
		})
	}
}

func TestParseMultipleGames(t *testing.T) {
	// arrange
	cases := []struct {
		input string
		want  []*Game
	}{
		{
			input: "1. d4 d5 2. c4 c6 *\n1. e4 e5 2. d4 exd4 3. c3 dxc3 4. Bc4 cxb2 5. Bxb2 d6 1/2-1/2",
			want: []*Game{
				{
					Result: ResultUnknown,
					Moves: []*Move{
						{Ply: 1, SAN: "d4"},
						{Ply: 2, SAN: "d5"},
						{Ply: 3, SAN: "c4"},
						{Ply: 4, SAN: "c6"},
					},
				},
				{
					Result: ResultDraw,
					Moves: []*Move{
						{Ply: 1, SAN: "e4"},
						{Ply: 2, SAN: "e5"},
						{Ply: 3, SAN: "d4"},
						{Ply: 4, SAN: "exd4"},
						{Ply: 5, SAN: "c3"},
						{Ply: 6, SAN: "dxc3"},
						{Ply: 7, SAN: "Bc4"},
						{Ply: 8, SAN: "cxb2"},
						{Ply: 9, SAN: "Bxb2"},
						{Ply: 10, SAN: "d6"},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			// act
			pgn, err := Parse(c.input)
			if err != nil {
				t.Error(err)
				return
			}

			// assert
			if pgn == nil {
				t.Error("pgn == nil")
				return
			}

			if len(pgn.Games) != len(c.want) {
				for _, game := range pgn.Games {
					t.Logf("%s\n", game)
				}
				t.Errorf("len(pgn.Games) want: %d got: %d", len(c.want), len(pgn.Games))
				return
			}

			for i := 0; i < len(pgn.Games); i++ {
				want := c.want[i]
				got := pgn.Games[i]

				if !want.Equals(got) {
					t.Errorf("\nwant:\n%s\n\ngot:\n%s\n", want, got)
				}
			}
		})
	}
}

func TestPGNtoMoves(t *testing.T) {
	// arrange
	cases := pgnMovesTestData(t)

	for i, c := range cases {
		t.Run(fmt.Sprintf("%04d", i+1), func(t *testing.T) {
			t.Parallel()

			// act
			pgnDB, err := Parse(c.PGN)
			if err != nil {
				t.Error(err)
				return
			}

			// assert
			if len(pgnDB.Games) != 1 {
				t.Errorf("len(pgnDB.Games): want: 1 got: %d", len(pgnDB.Games))
				return
			}

			pgn := pgnDB.Games[0]

			var uciMoves []string
			for _, m := range pgn.Moves {
				uciMoves = append(uciMoves, m.UCI)
			}

			if !reflect.DeepEqual(c.UCIMoves, uciMoves) {
				t.Errorf("\nwant:\n%v\ngot:\n%v", c.UCIMoves, uciMoves)
			}

			bb, err := bitboard.ParseFEN(bitboard.StartPos)
			if err != nil {
				t.Fatal(err)
			}

			var sanMoves []string
			for _, uciMove := range uciMoves {
				sanMove, err := bb.SAN(uciMove)
				if err != nil {
					t.Fatal(err)
				}

				sanMoves = append(sanMoves, sanMove)

				if bb, err = bb.Apply(uciMove); err != nil {
					t.Fatal(err)
				}
			}

			if !reflect.DeepEqual(c.SANMoves, sanMoves) {
				t.Errorf("\nwant:\n%v\ngot:\n%v", c.SANMoves, sanMoves)
			}

			// now test SAN to UCI. should break this out into a separate test.
			bb, err = bitboard.ParseFEN(bitboard.StartPos)
			if err != nil {
				t.Fatal(err)
			}

			for i, sanMove := range sanMoves {
				gotUCI, err := bb.UCI(sanMove)
				if err != nil {
					t.Fatal(err)
				}

				wantUCI := uciMoves[i]
				if wantUCI != gotUCI {
					t.Errorf("want: '%s' got: '%s'", wantUCI, gotUCI)
				}

				if bb, err = bb.Apply(wantUCI); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

type PGNMoves struct {
	PGN      string   `json:"pgn"`
	UCIMoves []string `json:"uciMoves"`
	SANMoves []string `json:"sanMoves"`
}

func pgnMovesTestData(tb testing.TB) []PGNMoves {
	fp, err := os.Open("testdata/pgn_uci_san.json")
	if err != nil {
		tb.Fatal(err)
	}
	defer fp.Close()

	var cases []PGNMoves

	dec := json.NewDecoder(fp)
	if err := dec.Decode(&cases); err != nil {
		tb.Fatal(err)
	}

	return cases
}

func BenchmarkParse(b *testing.B) {
	const pgnFilename = "testdata/TrollololFish.pgn"
	bytes, err := os.ReadFile(pgnFilename)
	if err != nil {
		b.Fatal(fmt.Errorf("'%s': %v", pgnFilename, err))
	}
	s := string(bytes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(s)
		if err != nil {
			b.Fatal(fmt.Errorf("'%s': %v", pgnFilename, err))
		}
	}
}

func BenchmarkParseReader(b *testing.B) {
	const pgnFilename = "testdata/TrollololFish.pgn"

	for i := 0; i < b.N; i++ {
		fp, err := os.Open(pgnFilename)
		if err != nil {
			b.Fatal(fmt.Errorf("'%s': %v", pgnFilename, err))
		}

		r := bufio.NewReaderSize(fp, 16384)

		if _, err = ParseReader(r); err != nil {
			_ = fp.Close()
			b.Fatal(fmt.Errorf("'%s': %v", pgnFilename, err))
		}

		_ = fp.Close()
	}
}

func BenchmarkIsCheck(b *testing.B) {
	const pgnFilename = "testdata/TrollololFish.pgn"
	bytes, err := os.ReadFile(pgnFilename)
	if err != nil {
		b.Fatal(fmt.Errorf("'%s': %v", pgnFilename, err))
	}
	s := string(bytes)

	pgn, err := Parse(s)
	if err != nil {
		b.Fatal(fmt.Errorf("'%s': %v", pgnFilename, err))
	}

	game := pgn.Games[0]

	var moves []string
	for _, move := range game.Moves {
		moves = append(moves, move.UCI)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb := bitboard.StartPosBoard()

		if _, err = bb.Apply(moves...); err != nil {
			b.Fatal(fmt.Errorf("'%s': moves: '%s': %v", pgnFilename, strings.Join(moves, ","), err))
		}
	}
}
