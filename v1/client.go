package openapi

import (
	"github.com/easyops-cn/giraffe-micro"
	"github.com/easyops-cn/giraffe-micro/plugins/restv2"
	"net/http"
)

var _ giraffe.Client = &Client{}

type Client struct {
	*restv2.Client
	transportOptions []TransportOption
}

type ClientOption func(*Client)

func NewClient(address, accessKey, secretKey string, options ...ClientOption) (*Client, error) {
	c := &Client{Client: &restv2.Client{
		Client:      &http.Client{},
		Middleware:  restv2.DefaultMiddleware,
		NameService: restv2.StaticAddress(address),
	}}
	for _, option := range options {
		option(c)
	}

	c.transportOptions = append(
		c.transportOptions,
		// the following Client settings override provided transport settings.
		RoundTripper(c.Client.Transport),
	)
	transport, _ := NewTransport(accessKey, secretKey, c.transportOptions...)
	c.Client.Transport = transport

	return c, nil
}

func WithClient(client *http.Client) ClientOption {
	return func(c *Client) {
		if client == nil {
			client = &http.Client{}
		}
		c.Client.Client = client
	}
}
