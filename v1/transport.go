package openapi

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

type TransportOption func(*transport)

type transport struct {
	rt        http.RoundTripper
	accessKey string
	secretKey string
	expires   func() string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	verb := strings.ToUpper(req.Method)
	url := req.URL.Path
	query := req.URL.Query()
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parameters := ""
	for _, k := range keys {
		parameters = fmt.Sprintf("%s%s%s", parameters, k, query.Get(k))
	}
	contentType := req.Header.Get("Content-type")
	contentMD5 := ""
	if req.GetBody != nil {
		reader, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		contentMD5 = fmt.Sprintf("%x", md5.Sum(b))
	}
	expires := fmt.Sprintf("%d", time.Now().Unix())
	if t.expires != nil {
		expires = t.expires()
	}

	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s", verb, url, parameters, contentType, contentMD5, expires, t.accessKey)
	mac := hmac.New(sha1.New, []byte(t.secretKey))
	mac.Write([]byte(stringToSign))
	signature := mac.Sum(nil)
	query.Add("accesskey", t.accessKey)
	query.Add("expires", expires)
	query.Add("signature", fmt.Sprintf("%x", signature))
	req.URL.RawQuery = query.Encode()
	return t.rt.RoundTrip(req)
}

func NewTransport(accessKey, secretKey string, options ...TransportOption) http.RoundTripper {
	t := &transport{
		rt:        http.DefaultTransport,
		accessKey: accessKey,
		secretKey: secretKey,
	}
	for _, option := range options {
		option(t)
	}
	return t
}

func RoundTripper(rt http.RoundTripper) TransportOption {
	return func(t *transport) {
		if rt != nil {
			t.rt = rt
		}
	}
}
