package bitboard

import (
	"testing"
)

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
