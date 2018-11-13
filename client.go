package natureremo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	baseURL = "https://api.nature.global/"
	version = "1"
)

// Client is an API client for Nature Remo Cloud API.
type Client struct {
	UserService      UserService
	DeviceService    DeviceService
	ApplianceService ApplianceService
	SignalService    SignalService

	HTTPClient  *http.Client
	AccessToken string
	BaseURL     string
}

// NewClient creates new client with access token of Nature Remo API.
// You can get access token from https://home.nature.global/.
func NewClient(accessToken string) *Client {
	var cli Client
	cli.AccessToken = accessToken
	cli.UserService = &userService{cli: &cli}
	cli.DeviceService = &deviceService{cli: &cli}
	cli.ApplianceService = &applianceService{cli: &cli}
	cli.SignalService = &signalService{cli: &cli}
	cli.BaseURL = baseURL + version
	return &cli
}

func (cli *Client) httpClient() *http.Client {
	if cli.HTTPClient != nil {
		return cli.HTTPClient
	}
	return http.DefaultClient
}

func (cli *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.AccessToken))
	return cli.httpClient().Do(req)
}

func (cli *Client) get(ctx context.Context, path string, params url.Values, v interface{}) error {
	reqURL := cli.BaseURL + "/" + path
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

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return errors.Wrap(err, "cannot parse HTTP body")
	}

	return nil
}

func (cli *Client) postForm(ctx context.Context, path string, data url.Values, v interface{}) error {
	reqURL := cli.BaseURL + "/" + path
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest(http.MethodPost, reqURL, body)
	if err != nil {
		return errors.Wrap(err, "cannot create HTTP request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	var respBody io.Reader = resp.Body

	// For Debug
	//var buf bytes.Buffer
	//if _, err := buf.ReadFrom(resp.Body); err != nil {
	//	return errors.Wrap(err, "cannot read HTTP body")
	//}
	//fmt.Println(buf.String())
	//respBody = &buf

	if err := json.NewDecoder(respBody).Decode(v); err != nil {
		return errors.Wrap(err, "cannot parse HTTP body")
	}

	return nil
}

func (cli *Client) post(ctx context.Context, path string, v interface{}) error {
	reqURL := cli.BaseURL + "/" + path
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
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

func (cli *Client) error(statusCode int, body io.Reader) error {
	buf, err := ioutil.ReadAll(body)
	if err != nil || len(buf) == 0 {
		return errors.Errorf("request failed with status code %d", statusCode)
	}
	return errors.Errorf("StatusCode: %d, Error: %s", statusCode, string(buf))
}
