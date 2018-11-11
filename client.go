package natureremo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	baseURL = "https://api.nature.global/"
	version = "1"
)

var apiURL = baseURL + version

type Client struct {
	UserService      UserService
	DeviceService    DeviceService
	ApplianceService ApplianceService

	HTTPClient  *http.Client
	AccessToken string
}

func NewClient(accessToken string) *Client {
	var cli Client
	cli.AccessToken = accessToken
	cli.UserService = &userService{cli: &cli}
	cli.DeviceService = &deviceService{cli: &cli}
	cli.ApplianceService = &applianceService{cli: &cli}
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
	reqURL := apiURL + "/" + path
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
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("request failed with status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return errors.Wrap(err, "cannot parse HTTP body")
	}

	return nil
}

func (cli *Client) postForm(ctx context.Context, path string, data url.Values, v interface{}) error {
	reqURL := apiURL + "/" + path
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
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("request failed with status code %d", resp.StatusCode)
	}

	if v == nil {
		return nil
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return errors.Wrap(err, "cannot parse HTTP body")
	}

	return nil
}
