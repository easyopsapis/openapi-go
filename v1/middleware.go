package openapi

import (
	"encoding/json"
	"errors"
	"github.com/easyopsapis/openapi-go/gerr"
	"io/ioutil"
	"net/http"

	"github.com/easyops-cn/giraffe-micro"
	"github.com/easyops-cn/giraffe-micro/plugins/restv2"
)

type Middleware struct {
	restv2.Middleware
}

func (m *Middleware) NewRequest(rule giraffe.HttpRule, in interface{}) (*http.Request, error) {
	req, err := m.Middleware.NewRequest(rule, in)
	if err != nil {
		return nil, err
	}

	// TODO 需要继续支持 body json 的处理
	if mask := fieldMask(in); mask != nil {
		withFieldMask(req, mask)
	}

	return req, nil
}

func (m *Middleware) ParseResponse(rule giraffe.HttpRule, resp *http.Response, out interface{}) error {
	if resp.StatusCode < 400 {
		return m.Middleware.ParseResponse(rule, resp, out)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	s := new(gerr.Message)
	if err := json.Unmarshal(body, s); err != nil {
		return err
	}

	if err := gerr.ErrorProto(s); err != nil {
		return err
	}

	return errors.New(resp.Status)
}
