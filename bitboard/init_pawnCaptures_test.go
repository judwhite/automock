package bitboard

import "testing"

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
