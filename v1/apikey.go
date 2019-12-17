package openapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"time"
)

type Signer interface {
	Sign(expires time.Time, request Request) error
}

type ApiKey struct {
	AccessKey, SecretKey string
}

func (a ApiKey) Sign(expires time.Time, r Request) error {
	s, err := r.StringToSign()
	if err != nil {
		return err
	}
	t := fmt.Sprintf("%d", expires.Unix())
	stringToSign := fmt.Sprintf("%s\n%s\n%s", s, t, a.AccessKey)
	mac := hmac.New(sha1.New, []byte(a.SecretKey))
	mac.Write([]byte(stringToSign))
	signature := mac.Sum(nil)
	r.SetSignature(a.AccessKey, fmt.Sprintf("%x", signature), expires)
	return nil
}
