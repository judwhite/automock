package chessdb

import (
	"context"
	"encoding/json"
	"net/url"

	"golang.org/x/xerrors"

	"stockhuman/httpcache"
)

const endpointURL = "https://www.chessdb.cn/cdb.php"

type QueryAllResponse struct {
	Status string `json:"status"`
	Moves  []Move `json:"moves"`
	Ply    int    `json:"ply"`
}

type Move struct {
	UCI string `json:"uci"`
	SAN string `json:"san"`
	// Score contains a large value e.g. 29999 if mate, -29998 is opponent can mate next turn.
	Score   int    `json:"score"`
	Rank    int    `json:"rank"`
	Note    string `json:"note"`
	WinRate string `json:"winrate"`
}

func QueryAll(ctx context.Context, fen string) (QueryAllResponse, error) {
	// https://www.chessdb.cn/cdb.php?action=queryall&json=1&board=rnbqkbnr/5ppp/p3p3/1p6/2BP4/5N2/PP3PPP/RNBQ1RK1%20w%20kq%20-%200%208
	// https://www.chessdb.cn/cdb.php?action=queryall&json=1&board=rnbqkbnr/pppp1ppp/8/4p3/6P1/5P2/PPPPP2P/RNBQKBNR%20b%20KQkq%20-%200%202
	// https://www.chessdb.cn/cdb.php?action=queryall&json=1&board=rnbqkbnr/pppp1ppp/8/4p3/8/5P2/PPPPP1PP/RNBQKBNR+w+KQkq+-+0+2
	params := make(url.Values)
	params.Set("action", "queryall")
	params.Set("json", "1")
	params.Set("board", fen)

	b, err := httpcache.Get(ctx, endpointURL, params, nil)
	if err != nil {
		return QueryAllResponse{}, xerrors.Errorf("%w", err)
	}

	var response QueryAllResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return QueryAllResponse{}, xerrors.Errorf("%w", err)
	}

	return response, nil
}

type QueryPVResponse struct {
	Status string   `json:"status"`
	Score  int      `json:"score"`
	Depth  int      `json:"depth"`
	PVUCI  []string `json:"pv"`
	PVSAN  []string `json:"pvSAN"`
}

func QueryPV(ctx context.Context, fen string) (QueryPVResponse, error) {
	// https://www.chessdb.cn/cdb.php??action=querypv&json=1&board=rnbqkbnr/5ppp/p3p3/1p6/2BP4/5N2/PP3PPP/RNBQ1RK1%20w%20kq%20-%200%208
	params := make(url.Values)
	params.Set("action", "querypv")
	params.Set("json", "1")
	params.Set("board", fen)

	b, err := httpcache.Get(ctx, endpointURL, params, nil)
	if err != nil {
		return QueryPVResponse{}, xerrors.Errorf("%w", err)
	}

	var response QueryPVResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return QueryPVResponse{}, xerrors.Errorf("%w", err)
	}

	return response, nil
}
