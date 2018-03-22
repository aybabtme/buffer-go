package buffer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	DefaultBaseURL  = "https://api.bufferapp.com"
	DefaultAuthURL  = DefaultBaseURL + "/oauth2/auth"
	DefaultTokenURL = DefaultBaseURL + "/oauth2/token"
)

type ClientOptions struct {
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

func (opts *ClientOptions) setDefault() error {
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

func NewClient(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		opts = new(ClientOptions)
	}
	if err := opts.setDefault(); err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(*opts.BaseURL)
	if err != nil {
		return nil, err
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

func (c *Client) get(ctx context.Context, ref *url.URL, v interface{}) error {

	req := mustRequest(http.NewRequest(
		"GET",
		c.baseURL.ResolveReference(ref).String(),
		nil,
	))
	req = req.WithContext(ctx)
	client := c.prepare(ctx)

	log.Print(req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) prepare(ctx context.Context) *http.Client {
	return c.oauthCfg.Client(ctx, &oauth2.Token{
		AccessToken: c.accessToken,
	})
}

func (c *Client) User() *UserService {
	return &UserService{client: c}
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
