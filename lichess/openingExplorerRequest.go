package lichess

import (
	"net/url"
	"strconv"
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

func NewOpeningExplorerRequest() OpeningExplorerRequest {
	return OpeningExplorerRequest{
		Variant:     "standard",
		FEN:         StartPos,
		Play:        "",
		Speeds:      ValidSpeeds,
		Ratings:     ValidRatings,
		Since:       Date{},
		Until:       Date{},
		Moves:       20,
		TopGames:    0,
		RecentGames: 0,
		History:     false,
	}
}

func (r OpeningExplorerRequest) QueryString() url.Values {
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
		values.Set("until", r.Since.String())
	}

	values.Set("moves", strconv.Itoa(r.Moves))
	values.Set("topGames", strconv.Itoa(r.TopGames))
	values.Set("recentGames", strconv.Itoa(r.RecentGames))

	if r.History {
		values.Set("history", strconv.FormatBool(r.History))
	}

	return values
}
