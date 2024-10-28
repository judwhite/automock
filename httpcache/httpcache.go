package httpcache

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"golang.org/x/xerrors"

	"stockhuman/commas"
	"stockhuman/utils"
)

const cacheDir = "./cache"

var (
	cache        = make(map[string][]byte)
	memCacheMtx  sync.RWMutex
	fileCacheMtx sync.RWMutex
)

func Get(ctx context.Context, skipCache bool, url string, query url.Values, header http.Header) ([]byte, bool, error) {
	queryKey := getQueryKey(url, query)

	if !skipCache {
		cachedResponse := getCachedResponse(queryKey)
		if cachedResponse != nil {
			return cachedResponse, true, nil
		}
	}

	queryString := query.Encode()
	if queryString != "" {
		url += "?" + queryString
	}

	utils.Log(fmt.Sprintf("... calling GET %s", url))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, false, xerrors.Errorf("failed to create request: %w", err)
	}

	for k, v := range header {
		req.Header.Set(k, v[0])
	}

	var c http.Client

	resp, err := c.Do(req)
	if err != nil {
		return nil, false, xerrors.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, xerrors.Errorf("GET %s, error reading body: %w", url, err)
	}

	utils.Log(fmt.Sprintf("... GET %s returned HTTP %s %s bytes", url, resp.Status, commas.Int(len(responseBody))))

	if resp.StatusCode != http.StatusOK {
		return nil, false, xerrors.Errorf("unexpected status code: %s response body: %s", resp.Status, string(responseBody))
	}

	storeCachedResponse(queryKey, url, query.Encode(), responseBody)

	return responseBody, false, nil
}

func getQueryKey(url string, query url.Values) string {
	names := make([]string, 0, len(query))
	for name := range query {
		names = append(names, name)
	}

	sort.Strings(names)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("url=%s", url))
	for _, key := range names {
		sb.WriteString(fmt.Sprintf("&%s=%s", key, query.Get(key)))
	}

	h := sha1.New()
	h.Write([]byte(sb.String()))
	return hex.EncodeToString(h.Sum(nil))
}

func getCachedResponse(queryKey string) []byte {
	memCacheMtx.RLock()
	cachedResponse, ok := cache[queryKey]
	memCacheMtx.RUnlock()

	if ok {
		return cachedResponse
	}

	fs := queryKeyToFS(queryKey)

	fileCacheMtx.RLock()
	b, err := os.ReadFile(fs.Filename)
	fileCacheMtx.RUnlock()
	if err != nil {
		return nil
	}

	memCacheMtx.Lock()
	cache[queryKey] = b
	memCacheMtx.Unlock()

	return b
}

type queryKeyFileSystem struct {
	Dir              string
	Filename         string
	MetadataFilename string
}

func queryKeyToFS(queryKey string) queryKeyFileSystem {
	prefix := queryKey[:2]
	dir := filepath.Join(cacheDir, prefix)

	return queryKeyFileSystem{
		Dir:              dir,
		Filename:         filepath.Join(dir, queryKey[:8]+".json"),
		MetadataFilename: filepath.Join(dir, queryKey[:8]+"-metadata.json"),
	}
}

func storeCachedResponse(queryKey, url, query string, body []byte) {
	memCacheMtx.Lock()
	cache[queryKey] = body
	memCacheMtx.Unlock()

	fs := queryKeyToFS(queryKey)
	metadataBody := []byte(fmt.Sprintf("url=%s\nquery=%s\nqueryKey=%s\n", url, query, queryKey))

	fileCacheMtx.Lock()
	defer fileCacheMtx.Unlock()

	if err := os.MkdirAll(fs.Dir, 0755); err != nil {
		return
	}
	if err := os.WriteFile(fs.Filename, body, 0644); err != nil {
		return
	}
	if err := os.WriteFile(fs.MetadataFilename, metadataBody, 0644); err != nil {
		return
	}
}
