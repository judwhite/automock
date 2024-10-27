package bitboard

import (
	"fmt"
	"strings"
	"testing"
)

func TestBoard_IsBlocked(t *testing.T) {
	cases := []struct {
		fen  string
		sq1  int
		sq2  int
		want bool
	}{
		{
			fen:  "4n3/6k1/8/8/8/8/6K1/4R3 w - -",
			sq1:  E1,
			sq2:  E8,
			want: false,
		},
		{
			fen:  "4n3/6k1/8/8/4P3/8/6K1/4R3 w - -",
			sq1:  E1,
			sq2:  E8,
			want: true,
		},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s__%s_%s",
			strings.ReplaceAll(c.fen, "/", "_"),
			squareNames[c.sq1],
			squareNames[c.sq2])

		t.Run(name, func(t *testing.T) {
			b, err := ParseFEN(c.fen)
			if err != nil {
				t.Error(err)
				return
			}

			got := b.IsBlocked(c.sq1, c.sq2)
			if c.want != got {
				t.Errorf("want: %v, got: %v", c.want, got)
			}
		})
	}
}
