package bitboard

import "testing"

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
