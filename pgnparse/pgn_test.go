package pgnparse

import (
	"fmt"
	"testing"
	"time"
)

func TestGameDate(t *testing.T) {
	cases := []struct {
		date string
		time string
		want time.Time
	}{
		{
			date: "2022.10.03",
			time: "13:34:56",
			want: time.Date(2022, 10, 03, 13, 34, 56, 0, time.UTC),
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s %s", c.date, c.time), func(t *testing.T) {
			// arrange
			var game Game
			if len(c.date) > 0 {
				game.Tags = append(game.Tags, Tag{"UTCDate", c.date})
			}
			if len(c.time) > 0 {
				game.Tags = append(game.Tags, Tag{"UTCTime", c.time})
			}

			// act
			got := game.Date()

			// assert
			if c.want != got {
				t.Errorf("want: %v got: %v", c.want, got)
			}
		})
	}
}
