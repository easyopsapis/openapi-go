package openapi

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	expiresQueryName   = "expires"
	accessKeyQueryName = "accessKey"
	signatureQueryName = "signature"
)

type Request interface {
	SetSignature(accessKey, signature string, expires time.Time)
	StringToSign() (string, error)
}

type request struct {
	*http.Request
}

func (r request) SetSignature(accessKey, signature string, expires time.Time) {
	query := r.URL.Query()
	query.Set(accessKeyQueryName, accessKey)
	query.Set(expiresQueryName, fmt.Sprintf("%d", expires.Unix()))
	query.Set(signatureQueryName, signature)
	r.URL.RawQuery = query.Encode()
}

func (r request) StringToSign() (string, error) {
	verb := strings.ToUpper(r.Method)
	path := r.URL.Path
	query := r.URL.Query()
	query.Del(expiresQueryName)
	query.Del(accessKeyQueryName)
	query.Del(signatureQueryName)
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parameters := ""
	for _, k := range keys {
		parameters = fmt.Sprintf("%s%s%s", parameters, k, query.Get(k))
	}
	contentType := r.Header.Get("Content-type")
	contentMD5 := ""
	if r.GetBody != nil {
		reader, err := r.GetBody()
		if err != nil {
			return "", err
		}
		b, err := ioutil.ReadAll(reader)
		if err != nil {
			return "", err
		}
		contentMD5 = fmt.Sprintf("%x", md5.Sum(b))
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s", verb, path, parameters, contentType, contentMD5), nil
}
