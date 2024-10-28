package lichess

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

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

	// NOTE: Response has been seen with duplicates.
	// URL: https://lichess.org/api/cloud-eval?fen=rnbqkb1r%2Fp1pppppp%2F1p3n2%2F8%2F8%2F1P3N2%2FPBPPPPPP%2FRN1QKB1R+b+KQkq+-+1+3&multiPv=3
	// Response:
	// {
	//     "fen": "rnbqkb1r/p1pppppp/1p3n2/8/8/1P3N2/PBPPPPPP/RN1QKB1R b KQkq - 1 3",
	//     "knodes": 50450,
	//     "depth": 31,
	//     "pvs": [
	//         {
	//             "moves": "d7d5 g2g3 g7g6 c2c4 e7e6 f1g2 f8g7 e1h1 e8h8 d2d4",
	//             "cp": 16
	//         },
	//         {
	//             "moves": "c7c5 c2c4 c8b7 g2g3 g7g6 f1g2 f8g7 e1h1 e8h8 d2d4",
	//             "cp": 17
	//         },
	//         {
	//             "moves": "c7c5 c2c4 c8b7 g2g3 g7g6 f1g2 f8g7 e1h1 e8h8 d2d4",
	//             "cp": 17
	//         },
	//         {
	//             "moves": "c8b7 g2g3 e7e6 c2c4 g7g6 f1g2 f8g7 e1h1 e8h8 d2d4",
	//             "cp": 19
	//         },
	//         {
	//             "moves": "c8b7 g2g3 e7e6 c2c4 g7g6 f1g2 f8g7 e1h1 e8h8 d2d4",
	//             "cp": 19
	//         },
	//         {
	//             "moves": "e7e6 g2g3 c8b7 c2c4 g7g6 f1g2 f8g7 e1h1 e8h8 d2d4",
	//             "cp": 21
	//         },
	//         {
	//             "moves": "e7e6 g2g3 c8b7 c2c4 g7g6 f1g2 f8g7 e1h1 e8h8 d2d4",
	//             "cp": 21
	//         },
	//         {
	//             "moves": "h7h6 c2c4 c8b7 e2e3 e7e6 b1c3 c7c5 g2g3 f8e7 f1g2",
	//             "cp": 30
	//         },
	//         {
	//             "moves": "h7h6 c2c4 c8b7 e2e3 e7e6 b1c3 c7c5 g2g3 f8e7 f1g2",
	//             "cp": 30
	//         }
	//     ]
	// }

	// filter duplicates
	seen := make(map[string]struct{})

	for i := 0; i < len(response.PVs); i++ {
		pv := response.PVs[i].MovesUCI

		uci := strings.Split(pv, " ")[0]
		if _, ok := seen[uci]; ok {
			response.PVs = append(response.PVs[:i], response.PVs[i+1:]...)
			i--
			continue
		}

		seen[uci] = struct{}{}
	}

	return response, nil
}
