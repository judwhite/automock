package lichess

import (
	"fmt"
	"strconv"
	"time"

	"automock/utils"
)

var (
	ValidSpeeds  = Speeds{UltraBullet, Bullet, Blitz, Rapid, Classical, Correspondence}
	ValidRatings = Ratings{R0, R1000, R1200, R1400, R1600, R1800, R2000, R2200, R2500}

	speedsOrder = map[Speed]int{
		UltraBullet:    1,
		Bullet:         2,
		Blitz:          3,
		Rapid:          4,
		Classical:      5,
		Correspondence: 6,
	}
)

const (
	UltraBullet    Speed = "ultraBullet"
	Bullet         Speed = "bullet"
	Blitz          Speed = "blitz"
	Rapid          Speed = "rapid"
	Classical      Speed = "classical"
	Correspondence Speed = "correspondence"

	R0    Rating = 0
	R1000 Rating = 1000
	R1200 Rating = 1200
	R1400 Rating = 1400
	R1600 Rating = 1600
	R1800 Rating = 1800
	R2000 Rating = 2000
	R2200 Rating = 2200
	R2500 Rating = 2500
)

// date

type Date struct {
	Year  int
	Month int
}

func (d *Date) IsZero() bool {
	return d.Year == 0 && d.Month == 0
}

func (d *Date) String() string {
	if d.IsZero() {
		return ""
	}
	return fmt.Sprintf("%04d-%02d", d.Year, d.Month)
}

func (d *Date) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	s := string(text)
	t, err := time.Parse("2006-01", s)
	if err != nil {
		return err
	}
	*d = Date{
		Year:  t.Year(),
		Month: int(t.Month()),
	}
	return nil
}

// speeds

type Speed string

func (s Speed) String() string {
	return string(s)
}

type Speeds []Speed

func (s Speeds) String() string {
	return utils.StringerSliceToString(s)
}

func (s Speeds) Contains(v string) (Speed, bool) {
	return utils.StringerSliceContains(s, v)
}

func (s Speeds) Len() int {
	return len(s)
}

func (s Speeds) Less(i, j int) bool {
	a, b := s[i], s[j]

	return speedsOrder[a] < speedsOrder[b]
}

func (s Speeds) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// ratings

type Rating int

func (r Rating) String() string {
	return strconv.Itoa(int(r))
}

type Ratings []Rating

func (r Ratings) String() string {
	return utils.StringerSliceToString(r)
}

func (r Ratings) Contains(v string) (Rating, bool) {
	return utils.StringerSliceContains(r, v)
}
