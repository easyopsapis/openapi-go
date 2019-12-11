package signature

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

func SignRequest(accessKey, secretKey string, t time.Time, req *http.Request) (string, error) {
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
			return "", err
		}
		b, err := ioutil.ReadAll(reader)
		if err != nil {
			return "", err
		}
		contentMD5 = fmt.Sprintf("%x", md5.Sum(b))
	}
	expires := fmt.Sprintf("%d", t.Unix())

	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s", verb, url, parameters, contentType, contentMD5, expires, accessKey)
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(stringToSign))
	signature := mac.Sum(nil)
	return fmt.Sprintf("%x", signature), nil
}
