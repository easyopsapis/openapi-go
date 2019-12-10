package openapi

import (
	"fmt"
	"github.com/easyops-cn/giraffe-micro"
	"github.com/easyops-cn/giraffe-micro/plugins/restv2"
	"net/http"
)

type wrapperMiddleware struct {
	name string
	restv2.Middleware
}

func (w *wrapperMiddleware) NewRequest(rule giraffe.HttpRule, in interface{}) (*http.Request, error) {
	req, err := w.Middleware.NewRequest(rule, in)
	if err == nil {
		req.URL.Path = fmt.Sprintf("/%s%s", w.name, req.URL.Path)
	}
	return req, nil
}

type wrapper struct {
	*restv2.Client
}

func WrapClient(name string, client *Client) giraffe.Client {
	return &wrapper{
		Client: &restv2.Client{
			Client:      client.Client.Client,
			Middleware:  &wrapperMiddleware{name: name, Middleware: client.Client.Middleware},
			NameService: client.Client.NameService,
		},
	}
}
