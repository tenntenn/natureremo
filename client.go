package natureremo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const version = "0.0.1"

const (
	baseURL    = "https://api.nature.global/"
	apiVersion = "1"
)

var defaultUserAgent string

func init() {
	defaultUserAgent = "tenntenn-natureremo/" + version + " (+https://github.com/tenntenn/natureremo)"
}

// Client is an API client for Nature Remo Cloud API.
type Client struct {
	UserService      UserService
	DeviceService    DeviceService
	ApplianceService ApplianceService
	SignalService    SignalService

	HTTPClient    *http.Client
	AccessToken   string
	BaseURL       string
	UserAgent     string
	LastRateLimit *RateLimit
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
	cli.BaseURL = baseURL + apiVersion
	cli.UserAgent = defaultUserAgent
	return &cli
}

func (cli *Client) getUA() string {
	if cli.UserAgent != "" {
		return cli.UserAgent
	}
	return defaultUserAgent
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
	req.Header.Set("User-Agent", cli.getUA())
	return cli.httpClient().Do(req)
}

func (cli *Client) get(ctx context.Context, path string, params url.Values, v interface{}) error {
	reqURL := cli.BaseURL + "/" + path
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

	rl, err := RateLimitFromHeader(resp.Header)
	if err != nil {
		return err
	}
	cli.LastRateLimit = rl

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("cannot parse HTTP body: %w", err)
	}

	return nil
}

func (cli *Client) postForm(ctx context.Context, path string, data url.Values, v interface{}) error {
	reqURL := cli.BaseURL + "/" + path
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest(http.MethodPost, reqURL, body)
	if err != nil {
		return fmt.Errorf("cannot create HTTP request: %w", err)
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
	//	return fmt.Errorf("cannot parse HTTP body: %w")
	//}
	//fmt.Println(buf.String())
	//respBody = &buf

	if err := json.NewDecoder(respBody).Decode(v); err != nil {
		return fmt.Errorf("cannot parse HTTP body: %w", err)
	}

	return nil
}

func (cli *Client) post(ctx context.Context, path string, v interface{}) error {
	reqURL := cli.BaseURL + "/" + path
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
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

func (cli *Client) error(statusCode int, body io.Reader) error {
	var aerr APIError
	if err := json.NewDecoder(body).Decode(&aerr); err != nil {
		return &APIError{HTTPStatus: statusCode}
	}
	aerr.HTTPStatus = statusCode
	return &aerr
}

// RateLimit has values of X-Rate-Limit-* in the response header.
type RateLimit struct {
	// Limit is a limit of request.
	Limit int64
	// Remaining is remaining request count.
	Remaining int64
	// Reset is time which a limit of request would be reseted.
	Reset time.Time
}

func RateLimitFromHeader(h http.Header) (*RateLimit, error) {
	ls := h.Get("X-Rate-Limit-Limit")
	if ls == "" {
		return nil, errors.New("cannot get X-Rate-Limit-Limit from header")
	}

	l, err := strconv.ParseInt(ls, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("X-Rate-Limit-Limit is invalid value: %w", err)
	}

	rs := h.Get("X-Rate-Limit-Remaining")
	if rs == "" {
		return nil, errors.New("cannot get X-Rate-Limit-Remaining from header")
	}

	r, err := strconv.ParseInt(rs, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("X-Rate-Limit-Remaining is invalid value: %w", err)
	}

	ts := h.Get("X-Rate-Limit-Reset")
	if ts == "" {
		return nil, errors.New("cannot get X-Rate-Limit-Reset from header")
	}

	t, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("X-Rate-Limit-Reset is invalid value: %w", err)
	}

	return &RateLimit{
		Limit:     l,
		Remaining: r,
		Reset:     time.Unix(t, 0),
	}, nil
}
