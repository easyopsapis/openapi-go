package openapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"time"
)

type Signer interface {
	Sign(request Request) error
}

type ApiKey struct {
	AccessKey, SecretKey string
}

func (a ApiKey) Sign(r Request) error {
	s, err := r.StringToSign()
	if err != nil {
		return err
	}
	expires := time.Now()
	t := fmt.Sprintf("%d", expires.Unix())
	stringToSign := fmt.Sprintf("%ss\n%ss\n%ss", s, t, a.AccessKey)
	mac := hmac.New(sha1.New, []byte(a.SecretKey))
	mac.Write([]byte(stringToSign))
	signature := mac.Sum(nil)
	r.SetSignature(a.AccessKey, fmt.Sprintf("%x", signature), expires)
	return nil
}
