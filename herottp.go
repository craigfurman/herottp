package herottp

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	*http.Client
}

type Config struct {
	NoFollowRedirect                  bool
	DisableTLSCertificateVerification bool
}

func New(config Config) *Client {
	c := &http.Client{}

	if config.NoFollowRedirect {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return NoFollowRedirect{}
		}
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.DisableTLSCertificateVerification,
		},
	}

	c.Transport = transport

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

func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)

	return c.Do(req)
}

type NoFollowRedirect struct{}

func (NoFollowRedirect) Error() string {
	return "This error should not ever be returned!"
}
