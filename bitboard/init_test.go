package bitboard

import "testing"

func TestInit_BitAfter(t *testing.T) {
	cases := []struct {
		name  string
		input Bits
		want  string
	}{
		{
			name:  "a4 e4",
			input: BitAfter[A4][E4],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 1 1 1 1\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "e4 a4",
			input: BitAfter[E4][A4],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "f4 e4",
			input: BitAfter[F4][E4],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 1 1 1 1 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "d1 d4",
			input: BitAfter[D1][D4],
			want: "" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "d8 d4",
			input: BitAfter[D8][D4],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n",
		},
		{
			name:  "b2 d4",
			input: BitAfter[B2][D4],
			want: "" +
				"0 0 0 0 0 0 0 1\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 1 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "f7 c4",
			input: BitAfter[F7][C4],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 1 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "b7 e4",
			input: BitAfter[B7][E4],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 0 1 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 0 1\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}

func TestInit_BitBetween(t *testing.T) {
	cases := []struct {
		name  string
		input Bits
		want  string
	}{
		{
			name:  "e1 e8",
			input: BitBetween[E1][E8],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "e8 e1",
			input: BitBetween[E8][E1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "a1 h8",
			input: BitBetween[A1][H8],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 1 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 1 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "a8 h1",
			input: BitBetween[A8][H1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 1 0 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 0 1 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "a1 h1",
			input: BitBetween[A1][H1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 1 1 1 1 1 1 0\n",
		},
		{
			name:  "h8 a8",
			input: BitBetween[H8][A8],
			want: "" +
				"0 1 1 1 1 1 1 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}

func TestInit_PawnCaptures(t *testing.T) {
	cases := []struct {
		name  string
		input Bits
		want  string
	}{
		{
			name:  "white a2",
			input: PawnCaptures[White][A2],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white b2",
			input: PawnCaptures[White][B2],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 0 1 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white g2",
			input: PawnCaptures[White][G2],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 1 0 1\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white h2",
			input: PawnCaptures[White][H2],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white a7",
			input: PawnCaptures[White][A7],
			want: "" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white a1",
			input: PawnCaptures[White][A1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "black a7",
			input: PawnCaptures[Black][A7],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "black b7",
			input: PawnCaptures[Black][B7],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 0 1 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "black h7",
			input: PawnCaptures[Black][H7],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}

func TestInit_PawnDefends(t *testing.T) {
	cases := []struct {
		name  string
		input Bits
		want  string
	}{
		{
			name:  "white b3",
			input: PawnDefends[White][B3],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 0 1 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white a3",
			input: PawnDefends[White][A3],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white g3",
			input: PawnDefends[White][G3],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 1 0 1\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white h3",
			input: PawnDefends[White][H3],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white a8",
			input: PawnDefends[White][A8],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white a1",
			input: PawnDefends[White][A1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "black a7",
			input: PawnDefends[Black][A7],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "black b6",
			input: PawnDefends[Black][B6],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"1 0 1 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "black h5",
			input: PawnDefends[Black][H5],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}

func TestInit_PawnMoves(t *testing.T) {
	cases := []struct {
		name  string
		input Bits
		want  string
	}{

		{
			name:  "white a2",
			input: PawnMoves[White][A2],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white h2",
			input: PawnMoves[White][H2],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 1\n" +
				"0 0 0 0 0 0 0 1\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white d5",
			input: PawnMoves[White][D5],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "white e7",
			input: PawnMoves[White][E7],
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
			name:  "black a7",
			input: PawnMoves[Black][A7],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "black e7",
			input: PawnMoves[Black][E7],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "black d2",
			input: PawnMoves[Black][D2],
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
			name:  "black e6",
			input: PawnMoves[Black][E6],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}

func TestInit_PieceMoves_Bishop(t *testing.T) {
	cases := []struct {
		name  string
		input Bits
		want  string
	}{
		{
			name:  "a1",
			input: PieceMoves[Bishop][A1],
			want: "" +
				"0 0 0 0 0 0 0 1\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 1 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 1 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "d4",
			input: PieceMoves[Bishop][D4],
			want: "" +
				"0 0 0 0 0 0 0 1\n" +
				"1 0 0 0 0 0 1 0\n" +
				"0 1 0 0 0 1 0 0\n" +
				"0 0 1 0 1 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 1 0 1 0 0 0\n" +
				"0 1 0 0 0 1 0 0\n" +
				"1 0 0 0 0 0 1 0\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}

func TestInit_PieceMoves_King(t *testing.T) {
	cases := []struct {
		name  string
		input Bits
		want  string
	}{
		{
			name:  "a1",
			input: PieceMoves[King][A1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 1 0 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n",
		},
		{
			name:  "h1",
			input: PieceMoves[King][H1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 1 1\n" +
				"0 0 0 0 0 0 1 0\n",
		},
		{
			name:  "a8",
			input: PieceMoves[King][A8],
			want: "" +
				"0 1 0 0 0 0 0 0\n" +
				"1 1 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "h8",
			input: PieceMoves[King][H8],
			want: "" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 1 1\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "e1",
			input: PieceMoves[King][E1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 1 1 1 0 0\n" +
				"0 0 0 1 0 1 0 0\n",
		},
		{
			name:  "e8",
			input: PieceMoves[King][E8],
			want: "" +
				"0 0 0 1 0 1 0 0\n" +
				"0 0 0 1 1 1 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "a4",
			input: PieceMoves[King][A4],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 1 0 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"1 1 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "h5",
			input: PieceMoves[King][H5],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 1 1\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 1 1\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "e4",
			input: PieceMoves[King][E4],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 1 1 1 0 0\n" +
				"0 0 0 1 0 1 0 0\n" +
				"0 0 0 1 1 1 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}

func TestInit_PieceMoves_Knight(t *testing.T) {
	cases := []struct {
		name  string
		input Bits
		want  string
	}{
		{
			name:  "b1",
			input: PieceMoves[Knight][B1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 0 1 0 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "g1",
			input: PieceMoves[Knight][G1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 1 0 1\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "b8",
			input: PieceMoves[Knight][B8],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"1 0 1 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "g8",
			input: PieceMoves[Knight][G8],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 0 1 0 1\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "h4",
			input: PieceMoves[Knight][H4],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 1 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 1 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "a1",
			input: PieceMoves[Knight][A1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 1 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "a8",
			input: PieceMoves[Knight][A8],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 1 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "h1",
			input: PieceMoves[Knight][H1],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 1 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "h8",
			input: PieceMoves[Knight][H8],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 1 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
		{
			name:  "f5",
			input: PieceMoves[Knight][F5],
			want: "" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 1 0 1 0\n" +
				"0 0 0 1 0 0 0 1\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 1 0 0 0 1\n" +
				"0 0 0 0 1 0 1 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}

func TestInit_PieceMoves_Rook(t *testing.T) {
	cases := []struct {
		name  string
		input Bits
		want  string
	}{
		{
			name:  "a1",
			input: PieceMoves[Rook][A1],
			want: "" +
				"1 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 0\n" +
				"0 1 1 1 1 1 1 1\n",
		},
		{
			name:  "d4",
			input: PieceMoves[Rook][D4],
			want: "" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"1 1 1 0 1 1 1 1\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}

func TestInit_PieceMoves_Queen(t *testing.T) {
	cases := []struct {
		name  string
		input Bits
		want  string
	}{
		{
			name:  "h1",
			input: PieceMoves[Queen][H1],
			want: "" +
				"1 0 0 0 0 0 0 1\n" +
				"0 1 0 0 0 0 0 1\n" +
				"0 0 1 0 0 0 0 1\n" +
				"0 0 0 1 0 0 0 1\n" +
				"0 0 0 0 1 0 0 1\n" +
				"0 0 0 0 0 1 0 1\n" +
				"0 0 0 0 0 0 1 1\n" +
				"1 1 1 1 1 1 1 0\n",
		},
		{
			name:  "e4",
			input: PieceMoves[Queen][E4],
			want: "" +
				"1 0 0 0 1 0 0 0\n" +
				"0 1 0 0 1 0 0 1\n" +
				"0 0 1 0 1 0 1 0\n" +
				"0 0 0 1 1 1 0 0\n" +
				"1 1 1 1 0 1 1 1\n" +
				"0 0 0 1 1 1 0 0\n" +
				"0 0 1 0 1 0 1 0\n" +
				"0 1 0 0 1 0 0 1\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}
