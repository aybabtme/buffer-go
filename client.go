package buffer

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	DefaultBaseURL  = "https://api.bufferapp.com"
	DefaultAuthURL  = DefaultBaseURL + "/oauth2/auth"
	DefaultTokenURL = DefaultBaseURL + "/oauth2/token"
)

type ClientConfig struct {
	// Required
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AccessToken  string

	// Optional

	BaseURL  *string
	AuthURL  *string
	TokenURL *string
}

func stringp(str string) *string { return &str }

func (opts *ClientConfig) setDefault() error {
	if opts.BaseURL == nil {
		opts.BaseURL = stringp(DefaultBaseURL)
	}
	if opts.AuthURL == nil {
		opts.AuthURL = stringp(DefaultAuthURL)
	}
	if opts.TokenURL == nil {
		opts.TokenURL = stringp(DefaultTokenURL)
	}
	if opts.ClientID == "" {
		return errors.New("missing ClientID")
	}
	if opts.ClientSecret == "" {
		return errors.New("missing ClientSecret")
	}
	if opts.RedirectURL == "" {
		return errors.New("missing RedirectURL")
	}
	if opts.AccessToken == "" {
		return errors.New("missing AccessToken")
	}
	return nil
}

type Client struct {
	baseURL     *url.URL
	oauthCfg    *oauth2.Config
	accessToken string
}

func NewClient(opts *ClientConfig) (*Client, error) {
	if opts == nil {
		opts = new(ClientConfig)
	}
	if err := opts.setDefault(); err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	baseURL, err := url.Parse(*opts.BaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "parsing base URL")
	}

	return &Client{
		baseURL: baseURL,
		oauthCfg: &oauth2.Config{
			ClientID:     opts.ClientID,
			ClientSecret: opts.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  *opts.AuthURL,
				TokenURL: *opts.TokenURL,
			},
		},
		accessToken: opts.AccessToken,
	}, nil
}

func (c *Client) User() *UserService {
	return &UserService{client: c}
}

func (c *Client) Profiles() *ProfilesService {
	return &ProfilesService{client: c}
}

func (c *Client) get(ctx context.Context, ref *url.URL, res interface{}) error {
	r := mustRequest(http.NewRequest(
		"GET",
		c.baseURL.ResolveReference(ref).String(),
		nil,
	))
	return errors.Wrap(c.do(ctx, r, res), "performing GET request")
}

func (c *Client) post(ctx context.Context, ref *url.URL, req, res interface{}) error {
	body := bytes.NewBuffer(nil)
	if req != nil {
		if err := json.NewEncoder(body).Encode(req); err != nil {
			return errors.Wrap(err, "encoding post request body")
		}
	}
	r := mustRequest(http.NewRequest(
		"POST",
		c.baseURL.ResolveReference(ref).String(),
		body,
	))
	return errors.Wrap(c.do(ctx, r, res), "perfoming POST request")
}

func (c *Client) do(ctx context.Context, req *http.Request, res interface{}) error {
	req = req.WithContext(ctx)
	client := c.prepare(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "executing request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "reading response")
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	if d {
		log.Println(string(body))
		log.Println("===")
	}

	return errors.Wrap(json.Unmarshal(body, res), "decoding JSON response")
}

func (c *Client) prepare(ctx context.Context) *http.Client {
	return c.oauthCfg.Client(ctx, &oauth2.Token{
		AccessToken: c.accessToken,
	})
}

func mustRequest(req *http.Request, err error) *http.Request {
	if err != nil {
		panic(err)
	}
	return req
}

func mustURL(str string) *url.URL {
	u, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	return u
}
