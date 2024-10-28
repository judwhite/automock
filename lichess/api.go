package lichess

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"golang.org/x/xerrors"

	"stockhuman/httpcache"
)

func GetLichessGames(ctx context.Context, req OpeningExplorerRequest) (OpeningExplorerResponse, error) {
	const endpointURL = "https://explorer.lichess.ovh/lichess"

	b, err := httpcache.Get(ctx, endpointURL, req.QueryString(), authHeader)
	if err != nil {
		return OpeningExplorerResponse{}, xerrors.Errorf("%w", err)
	}

	var response OpeningExplorerResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return OpeningExplorerResponse{}, xerrors.Errorf("%w", err)
	}

	return response, nil
}

type CloudEvalResponse struct {
	FEN    string        `json:"fen"`
	KNodes int           `json:"knodes"`
	Depth  int           `json:"depth"`
	PVs    []CloudEvalPV `json:"pvs"`
}

type CloudEvalPV struct {
	// MovesUCI contains a space separated list of UCI moves
	MovesUCI string `json:"moves"`
	CP       int    `json:"cp"`
	Mate     int    `json:"mate,omitempty"`
}

func GetCloudEval(ctx context.Context, fen string, multiPV int) (CloudEvalResponse, error) {
	// https://lichess.org/api/cloud-eval?fen=rnbqkbnr/pppp1ppp/8/4p3/6P1/5P2/PPPPP2P/RNBQKBNR%20b%20KQkq%20-%200%202&multiPv=3
	const endpointURL = "https://lichess.org/api/cloud-eval"

	params := make(url.Values)
	params.Set("fen", fen)
	params.Set("multiPv", strconv.Itoa(multiPV))

	b, err := httpcache.Get(ctx, endpointURL, params, authHeader)
	if err != nil {
		return CloudEvalResponse{}, xerrors.Errorf("%w", err)
	}

	var response CloudEvalResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return CloudEvalResponse{}, xerrors.Errorf("%w", err)
	}

	return response, nil
}
