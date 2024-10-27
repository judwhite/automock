package lichess

import (
	"context"
	"encoding/json"

	"golang.org/x/xerrors"
)

func GetLichessGames(ctx context.Context, req OpeningExplorerRequest) (OpeningExplorerResponse, error) {
	const endpointURL = "https://explorer.lichess.ovh/lichess"

	b, err := httpGet(ctx, endpointURL, req.QueryString())
	if err != nil {
		return OpeningExplorerResponse{}, xerrors.Errorf("%w", err)
	}

	var response OpeningExplorerResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return OpeningExplorerResponse{}, xerrors.Errorf("%w", err)
	}

	return response, nil
}
