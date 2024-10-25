package lichess

type OpeningExplorerGame struct {
	UCI    string `json:"uci,omitempty"`
	ID     string `json:"id"`
	Winner string `json:"winner"`
	Speed  string `json:"speed"`
	Mode   string `json:"mode"`
	Black  struct {
		Name   string `json:"name"`
		Rating int    `json:"rating"`
	} `json:"black"`
	White struct {
		Name   string `json:"name"`
		Rating int    `json:"rating"`
	} `json:"white"`
	Year  int    `json:"year"`
	Month string `json:"month"`
}

type OpeningExplorerOpening struct {
	ECO  string `json:"eco"`
	Name string `json:"name"`
}

// OpeningExplorerMove represents a move played from this position.
type OpeningExplorerMove struct {
	UCI           string `json:"uci"`
	SAN           string `json:"san"`
	AverageRating int    `json:"averageRating"`
	White         int    `json:"white"`
	Draws         int    `json:"draws"`
	Black         int    `json:"black"`
	// Game is populated if only a single game was played from this position.
	// TODO: do we really care?
	// Game *OpeningExplorerGame `json:"game,omitempty"`
}

func (oem OpeningExplorerMove) Total() int {
	return oem.White + oem.Draws + oem.Black
}

// OpeningExplorerHistory contains historical data about the position.
type OpeningExplorerHistory struct {
	Month string `json:"month"`
	Black int    `json:"black"`
	Draws int    `json:"draws"`
	White int    `json:"white"`
}

type OpeningExplorerResponse struct {
	Opening     *OpeningExplorerOpening  `json:"opening,omitempty"`
	White       int                      `json:"white"`
	Draws       int                      `json:"draws"`
	Black       int                      `json:"black"`
	Moves       []OpeningExplorerMove    `json:"moves,omitempty"`
	TopGames    []OpeningExplorerGame    `json:"topGames,omitempty"`
	RecentGames []OpeningExplorerGame    `json:"recentGames,omitempty"`
	History     []OpeningExplorerHistory `json:"history,omitempty"`
}

func (oer OpeningExplorerResponse) Total() int {
	return oer.White + oer.Draws + oer.Black
}
