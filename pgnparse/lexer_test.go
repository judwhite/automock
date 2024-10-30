package pgnparse

import (
	"fmt"
	"testing"
)

// TODO: Many good examples here
// TODO: https://github.com/kevinludwig/pgn-parser/blob/master/test/test_grammar.js

func TestLexerRun_SingleGames(t *testing.T) {
	// arrange
	cases := []struct {
		input string
	}{
		{input: "1. e4 e5 2. Nf3 Nc6"},
		{input: "1. e4 e5 (1...c5 2. Nf3 Nc6) 2. f4 exf4 *"},
		{input: "1. f3 $2 e5 $1 2. g4 $4 $15 Qh4#\n0-1"},
		{input: "1. d4 {some commentary then} 1...d5 2. c4 dxc4 *"},
		{input: "1. e4 (1. d4 d5 ) e5 2. d4 (2. Nf3 Nc6 ) exd4 *"},
		{input: "1. e4 ({c1} {c2} 1. d4 d5) e5 2. d4 exd4 *"},
		{input: "1. e4 (1. d4 d5) (1. c4 e5) e5 *"},
		{input: "1. e4 (1. d4 d5)(1. f4 e5)e5 2. d4 exd4 *"},
		{input: "1. e4 (1. d4 d5)e5 2. d4 exd4 *"},
		{input: "1. e4 (1. d4 ) e5 2. d4 exd4 *"},
		{input: "1. e4? e5! 2. Nf3?? Nc6?! 3. Bc4!? Nf6!!\n0-1"},
		{input: "1. e4 f5 2. exf5 g6 3. fxg6 h6 4. g7 Rh7 5. gxf8Q+ *"},
		{input: "{start of game} 1. f3 e5 2. g4 Qh4#\n0-1"},
		{input: "1. e4 {[%clk 0:15:09.9]} 1... Nc6 {[%clk 0:15:06.6]} 2. Nf3 {[%clk 0:15:15.1]} 2... e5 {[%clk 0:15:02.9]} 3. Bc4 {[%clk 0:15:16.9]} 3... h6 {[%clk 0:14:39.1]} 4. d4 {[%clk 0:15:23.1]} 4... exd4 {[%clk 0:14:30.7]} 5. Nxd4 {[%clk 0:15:29.1]} 5... Bc5 {[%clk 0:14:19]} 6. c3 {[%clk 0:15:06.4]} 6... Qe7 {[%clk 0:14:17.3]} 7. Qf3 {[%clk 0:14:48]} 7... Nf6 {[%clk 0:14:22.9]} 8. O-O {[%clk 0:14:49.9]} 8... Qxe4 {[%clk 0:14:00.9]} 9. Qxe4+ {[%clk 0:14:51]} 9... Nxe4 {[%clk 0:14:07.7]} 10. Re1 {[%clk 0:15:00.2]} 10... f5 {[%clk 0:13:48.3]} 11. f3 {[%clk 0:14:47.9]} 11... Bxd4+ {[%clk 0:12:38.1]} 12. cxd4 {[%clk 0:14:57.8]} 12... Nxd4 {[%clk 0:12:47]} 1-0\n"},
		{input: `[Event "F/S Return Match"]
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
40. Rd6 Kc5 41. Ra6 Nf2 42. g4 Bd3 43. Re6 1/2-1/2`},
		{input: `[Event "Ch World (match)"]
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
28. Nc3 e2 29. Raxd1 Qxc3 0-1`},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			// act
			l, err := lex([]byte(c.input))

			// assert
			if err != nil {
				// debug output
				for _, item := range l.items {
					fmt.Printf("%20s pos: %4d val: %s\n", item.typ, item.pos, item.val)
				}

				t.Error(err)
			}
		})
	}
}

func TestLexerRun_MultipleGames(t *testing.T) {
	// arrange
	cases := []struct {
		input string
	}{
		{
			input: "1. d4 d5 2. c4 c6 *\n1. e4 e5 2. d4 exd4 3. c3 dxc3 4. Bc4 cxb2 5. Bxb2 d6 1/2-1/2",
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			// act
			l, err := lex([]byte(c.input))

			// assert
			if err != nil {
				// debug output
				for _, item := range l.items {
					fmt.Printf("%20s pos: %4d val: %s\n", item.typ, item.pos, item.val)
				}

				t.Error(err)
			}
		})
	}
}
