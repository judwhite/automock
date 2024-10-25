package lichess

import (
	"context"
	"encoding/json"

	"golang.org/x/xerrors"
)

const StartPos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type QueryStringer interface {
	QueryString() string
}

func GetLichessGames(ctx context.Context, req OpeningExplorerRequest) (OpeningExplorerResponse, error) {
	const endpointURL = "https://explorer.lichess.ovh/lichess"

	b, err := httpGet(ctx, endpointURL, req)
	if err != nil {
		return OpeningExplorerResponse{}, xerrors.Errorf("%w", err)
	}

	//fmt.Printf("ORIGINAL RESPONSE:\n\n%s\n\n%s\n\n", string(b), strings.Repeat("=", 40))

	var response OpeningExplorerResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return OpeningExplorerResponse{}, xerrors.Errorf("%w", err)
	}

	return response, nil
}
