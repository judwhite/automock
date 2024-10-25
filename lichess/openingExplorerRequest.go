package lichess

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"stockhuman/utils"
)

type OpeningExplorerRequest struct {
	Speeds          Speeds
	Ratings         Ratings
	Since           Date
	Modes           Modes
	FEN             string
	History         bool
	OmitTopGames    bool
	OmitRecentGames bool
}

func (r OpeningExplorerRequest) QueryString() string {
	values := url.Values{}
	if len(r.Speeds) > 0 {
		values.Set("speeds", r.Speeds.String())
	}
	if len(r.Ratings) > 0 {
		values.Set("ratings", r.Ratings.String())
	}
	if !r.Since.IsZero() {
		values.Set("since", r.Since.String())
	}
	if len(r.Modes) > 0 {
		values.Set("modes", r.Modes.String())
	}
	if r.FEN != "" {
		values.Set("fen", r.FEN)
	}
	if r.History {
		values.Set("history", "true")
	}
	if r.OmitTopGames {
		values.Set("topGames", "0")
	}
	if r.OmitRecentGames {
		values.Set("recentGames", "0")
	}
	return values.Encode()
}

// date

type Date struct {
	Year  int
	Month int
}

func (d Date) IsZero() bool {
	return d.Year == 0 && d.Month == 0
}

func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return fmt.Sprintf("%04d-%02d", d.Year, d.Month)
}

func (d *Date) UnmarshalText(text []byte) error {
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

const (
	UltraBullet    Speed = "ultrabullet"
	Bullet         Speed = "bullet"
	Blitz          Speed = "blitz"
	Rapid          Speed = "rapid"
	Classical      Speed = "classical"
	Correspondence Speed = "correspondence"
)

type Speeds []Speed

func (s Speeds) String() string {
	return utils.SliceToString(s)
}

// ratings

type Rating int

const (
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

func (r Rating) String() string {
	return strconv.Itoa(int(r))
}

type Ratings []Rating

func (r Ratings) String() string {
	return utils.SliceToString(r)
}

// modes

type Mode string

func (m Mode) String() string {
	return string(m)
}

const (
	Casual Mode = "casual"
	Rated  Mode = "rated"
)

type Modes []Mode

func (m Modes) String() string {
	return utils.SliceToString(m)
}
