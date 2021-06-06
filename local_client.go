package natureremo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LocalClient struct {
	addr       string
	HTTPClient *http.Client
}

func NewLocalClient(addr string) *LocalClient {
	return &LocalClient{addr: addr}
}

func (cli *LocalClient) httpClient() *http.Client {
	if cli.HTTPClient != nil {
		return cli.HTTPClient
	}
	return http.DefaultClient
}

func (cli *LocalClient) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Set("X-Requested-With", "tenntenn/natureremo")
	req.Header.Set("Expect", "")
	return cli.httpClient().Do(req)
}

func (cli *LocalClient) get(ctx context.Context, path string, params url.Values, v interface{}) error {
	reqURL := "http://" + cli.addr + "/" + path
	if params != nil {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("cannot create HTTP request: %w", err)
	}

	resp, err := cli.do(ctx, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return cli.error(resp.StatusCode, resp.Body)
	}

	if v == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("cannot parse HTTP body: %w", err)
	}

	return nil
}

func (cli *LocalClient) post(ctx context.Context, path string, body io.Reader, v interface{}) error {
	reqURL := "http://" + cli.addr + "/" + path
	req, err := http.NewRequest(http.MethodPost, reqURL, body)
	if err != nil {
		return fmt.Errorf("cannot create HTTP request: %w", err)
	}

	resp, err := cli.do(ctx, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return cli.error(resp.StatusCode, resp.Body)
	}

	if v == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("cannot parse HTTP body: %w", err)
	}

	return nil
}

func (cli *LocalClient) error(statusCode int, body io.Reader) error {
	buf, err := ioutil.ReadAll(body)
	if err != nil || len(buf) == 0 {
		return fmt.Errorf("request failed with status code %d", statusCode)
	}
	return fmt.Errorf("StatusCode: %d, Error: %s", statusCode, string(buf))
}

func (cli *LocalClient) Fetch(ctx context.Context) (*IRSignal, error) {
	var ir IRSignal
	if err := cli.get(ctx, "messages", nil, &ir); err != nil {
		return nil, fmt.Errorf("GET messages failed: %w", err)
	}
	return &ir, nil
}

func (cli *LocalClient) Emit(ctx context.Context, ir *IRSignal) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ir); err != nil {
		return fmt.Errorf("cannot encode IRSignal %v: %w", ir, err)
	}

	if err := cli.post(ctx, "messages", &buf, nil); err != nil {
		return fmt.Errorf("POST messages failed with %s: %w", buf.String(), err)
	}

	return nil
}
