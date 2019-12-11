package openapi

import (
	"fmt"
	"github.com/easyopsapis/openapi-go/signature/v1"
	"net/http"
	"time"
)

type TransportOption func(*transport)

type transport struct {
	rt        http.RoundTripper
	accessKey string
	secretKey string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	now := time.Now()
	sign, err := signature.SignRequest(t.accessKey, t.secretKey, now, req)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("accesskey", t.accessKey)
	query.Add("expires", fmt.Sprintf("%d", now.Unix()))
	query.Add("signature", sign)
	req.URL.RawQuery = query.Encode()
	return t.rt.RoundTrip(req)
}

func NewTransport(accessKey, secretKey string, options ...TransportOption) (http.RoundTripper, error) {
	t := &transport{
		rt:        http.DefaultTransport,
		accessKey: accessKey,
		secretKey: secretKey,
	}
	for _, option := range options {
		option(t)
	}
	return t, nil
}

func RoundTripper(rt http.RoundTripper) TransportOption {
	return func(t *transport) {
		if rt != nil {
			t.rt = rt
		}
	}
}
