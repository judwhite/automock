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
