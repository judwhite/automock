package bitboard

import "testing"

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
