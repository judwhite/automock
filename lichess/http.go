package lichess

import (
	"context"
	"io"
	"net/http"

	"golang.org/x/xerrors"
)

func httpGet(ctx context.Context, url string, query QueryStringer) ([]byte, error) {
	queryString := query.QueryString()
	if queryString != "" {
		url += "?" + queryString
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to create request: %w", err)
	}

	var c http.Client

	resp, err := c.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("GET %s, error reading body: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, xerrors.Errorf("unexpected status code: %s response body: %s", resp.Status, string(responseBody))
	}

	return responseBody, nil
}
