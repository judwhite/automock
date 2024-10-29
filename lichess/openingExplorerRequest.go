package lichess

import (
	"net/url"
	"strconv"

	"automock/bitboard"
)

type OpeningExplorerRequest struct {
	Variant     string
	FEN         string
	Play        string
	Speeds      Speeds
	Ratings     Ratings
	Since       Date
	Until       Date
	Moves       int
	TopGames    int
	RecentGames int
	History     bool
}

func (r OpeningExplorerRequest) QueryString() url.Values {
	// set defaults

	if r.Variant == "" {
		r.Variant = "standard"
	}
	if r.FEN == "" || r.FEN == "startpos" {
		r.FEN = bitboard.StartPos
	}
	if r.Moves == 0 {
		r.Moves = 20
	}
	if len(r.Speeds) == 0 {
		r.Speeds = ValidSpeeds
	}
	if len(r.Ratings) == 0 {
		r.Ratings = ValidRatings
	}

	// set query param values

	values := make(url.Values)

	values.Set("variant", r.Variant)
	values.Set("fen", r.FEN)

	if r.Play != "" {
		values.Set("play", r.Play)
	}

	values.Set("speeds", r.Speeds.String())
	values.Set("ratings", r.Ratings.String())

	if !r.Since.IsZero() {
		values.Set("since", r.Since.String())
	}
	if !r.Until.IsZero() {
		values.Set("until", r.Until.String())
	}

	values.Set("moves", strconv.Itoa(r.Moves))
	values.Set("topGames", strconv.Itoa(r.TopGames))
	values.Set("recentGames", strconv.Itoa(r.RecentGames))

	if r.History {
		values.Set("history", strconv.FormatBool(r.History))
	}

	return values
}
