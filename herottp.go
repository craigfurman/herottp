package herottp

import (
	"crypto/tls"
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

type NoFollowRedirect struct{}

func (NoFollowRedirect) Error() string {
	return "This error should not ever be returned!"
}
