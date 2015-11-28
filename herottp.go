package herottp

import (
	"net/http"
	"net/url"
)

type Client struct {
	*http.Client
}

type Config struct {
	NoFollowRedirect bool
}

func New(config Config) *Client {
	c := &http.Client{}
	if config.NoFollowRedirect {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return NoFollowRedirect{}
		}
	}

	return &Client{
		Client: c,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.Client.Do(req)
	if e, isURLErr := err.(*url.Error); isURLErr {
		if _, ok := e.Err.(NoFollowRedirect); ok {
			return resp, nil
		}
	}

	return resp, err
}

type NoFollowRedirect struct{}

func (NoFollowRedirect) Error() string {
	return "This error should not ever be returned!"
}
