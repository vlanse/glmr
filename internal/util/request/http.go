package request

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GET(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	return doRequest(ctx, http.MethodGet, url, headers)
}

func MustURL(hostWithSchemaAndPath string, queryKV ...string) string {
	u, err := url.Parse(hostWithSchemaAndPath)
	if err != nil {
		panic(fmt.Errorf("bad URL %q: %w", hostWithSchemaAndPath, err))
	}

	var k string
	q := u.Query()
	for i, v := range queryKV {
		if i%2 == 0 {
			k = v
		} else {
			q.Add(k, v)
		}

	}

	u.RawQuery = q.Encode()
	return u.String()
}

func doRequest(ctx context.Context, method string, url string, headers map[string]string) ([]byte, error) {
	req, _ := http.NewRequestWithContext(ctx, method, url, nil)
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("GET request to %s error: %w", url, err)
	}

	defer func() { _ = resp.Body.Close() }()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response from %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("%q request to %s failed with code %s", method, url, resp.Status)
	}

	return data, nil
}
