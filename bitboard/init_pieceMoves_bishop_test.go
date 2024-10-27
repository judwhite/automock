package bitboard

import "testing"

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
