package herottp

import "net/http"

type Client struct {
	*http.Client
}

type Config struct {
	FollowRedirects bool
}

func New(config Config) *Client {
	return &Client{
		Client: http.DefaultClient,
	}
}
