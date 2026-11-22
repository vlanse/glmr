package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func GET(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	return doRequest(ctx, http.MethodGet, url, headers, nil)
}

func POST(ctx context.Context, url string, headers map[string]string, body []byte) ([]byte, error) {
	return doRequest(ctx, http.MethodPost, url, headers, io.NopCloser(bytes.NewBuffer(body)))
}

func WithQuery(base string, queryKV ...string) string {
	var k string
	q := url.Values{}
	for i, v := range queryKV {
		if i%2 == 0 {
			k = v
		} else {
			q.Add(k, v)
		}
	}

	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}

	return base + sep + q.Encode()
}

func doRequest(ctx context.Context, method string, url string, headers map[string]string, body io.ReadCloser) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build %s request to %s: %w", method, url, err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Body = body

	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("%s request to %s error: %w", method, url, err)
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
