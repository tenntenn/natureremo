package natureremo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
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
		return errors.Wrap(err, "cannot create HTTP request")
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
		return errors.Wrap(err, "cannot parse HTTP body")
	}

	return nil
}

func (cli *LocalClient) post(ctx context.Context, path string, body io.Reader, v interface{}) error {
	reqURL := "http://" + cli.addr + "/" + path
	req, err := http.NewRequest(http.MethodPost, reqURL, body)
	if err != nil {
		return errors.Wrap(err, "cannot create HTTP request")
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
		return errors.Wrap(err, "cannot parse HTTP body")
	}

	return nil
}

func (cli *LocalClient) error(statusCode int, body io.Reader) error {
	buf, err := ioutil.ReadAll(body)
	if err != nil || len(buf) == 0 {
		return errors.Errorf("request failed with status code %d", statusCode)
	}
	return errors.Errorf("StatusCode: %d, Error: %s", statusCode, string(buf))
}

func (cli *LocalClient) Fetch(ctx context.Context) (*IRSignal, error) {
	var ir IRSignal
	if err := cli.get(ctx, "messages", nil, &ir); err != nil {
		return nil, errors.Wrap(err, "GET messages failed")
	}
	return &ir, nil
}

func (cli *LocalClient) Emit(ctx context.Context, ir *IRSignal) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ir); err != nil {
		return errors.Wrapf(err, "cannot encode IRSignal %v", ir)
	}

	if err := cli.post(ctx, "messages", &buf, nil); err != nil {
		return errors.Wrapf(err, "POST messages failed with %s", buf.String())
	}

	return nil
}
