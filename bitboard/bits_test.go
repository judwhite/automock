package bitboard

import (
	"fmt"
	"testing"
)

func TestBits_String(t *testing.T) {
	cases := []struct {
		input Bits
		want  string
	}{
		{
			input: 0x8100000000000081,
			want: "" +
				"1 0 0 0 0 0 0 1\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"1 0 0 0 0 0 0 1\n",
		},
		{
			input: 0x8000000000000001,
			want: "" +
				"1 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 0\n" +
				"0 0 0 0 0 0 0 1\n",
		},
		{
			input: 0x8040201008040201,
			want: "" +
				"1 0 0 0 0 0 0 0\n" +
				"0 1 0 0 0 0 0 0\n" +
				"0 0 1 0 0 0 0 0\n" +
				"0 0 0 1 0 0 0 0\n" +
				"0 0 0 0 1 0 0 0\n" +
				"0 0 0 0 0 1 0 0\n" +
				"0 0 0 0 0 0 1 0\n" +
				"0 0 0 0 0 0 0 1\n",
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%016X", uint64(c.input)), func(t *testing.T) {
			if got := c.input.String(); got != c.want {
				t.Errorf("want:\n%s\ngot:\n%s\n", c.want, got)
			}
		})
	}
}
