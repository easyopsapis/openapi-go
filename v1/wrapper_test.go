package openapi

import (
	"errors"
	"github.com/easyops-cn/giraffe-micro"
	"github.com/easyops-cn/giraffe-micro/plugins/restv2"
	"net/http"
	"reflect"
	"testing"
)

type testMiddleware struct {
	newRequest    func(rule giraffe.HttpRule, in interface{}) (*http.Request, error)
	parseResponse func(rule giraffe.HttpRule, resp *http.Response, out interface{}) error
}

func (t *testMiddleware) NewRequest(rule giraffe.HttpRule, in interface{}) (*http.Request, error) {
	return t.newRequest(rule, in)
}

func (t *testMiddleware) ParseResponse(rule giraffe.HttpRule, resp *http.Response, out interface{}) error {
	return t.parseResponse(rule, resp, out)
}

func TestWrapClient(t *testing.T) {
	type args struct {
		name   string
		client *Client
	}
	tests := []struct {
		name string
		args args
		want giraffe.Client
	}{
		{
			args: args{
				name: "cmdb",
				client: &Client{
					Client: &restv2.Client{
						Client: &http.Client{
							Transport: &transport{
								sig: ApiKey{
									AccessKey: "3fc93fed595063856df3ee1a",
									SecretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
								},
								rt: http.DefaultTransport,
							},
						},
						Middleware:  restv2.DefaultMiddleware,
						NameService: restv2.StaticAddress("192.168.100.162:8080"),
					},
					transportOptions: []TransportOption{RoundTripper(nil)},
				},
			},
			want: &wrapper{&restv2.Client{
				Client: &http.Client{
					Transport: &transport{
						sig: ApiKey{
							AccessKey: "3fc93fed595063856df3ee1a",
							SecretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
						},
						rt: http.DefaultTransport,
					},
				},
				Middleware:  &wrapperMiddleware{name: "cmdb", Middleware: restv2.DefaultMiddleware},
				NameService: restv2.StaticAddress("192.168.100.162:8080"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WrapClient(tt.args.name, tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WrapClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrapperMiddleware_NewRequest(t *testing.T) {
	type fields struct {
		name       string
		Middleware restv2.Middleware
	}
	type args struct {
		rule giraffe.HttpRule
		in   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			fields: fields{
				name: "cmdb",
				Middleware: &testMiddleware{
					newRequest: func(rule giraffe.HttpRule, in interface{}) (request *http.Request, err error) {
						return http.NewRequest("GET", "/api/instance", nil)
					},
				},
			},
			want: func() *http.Request {
				req, _ := http.NewRequest("GET", "/cmdb/api/instance", nil)
				return req
			}(),
		},
		{
			fields: fields{
				name: "cmdb",
				Middleware: &testMiddleware{
					newRequest: func(rule giraffe.HttpRule, in interface{}) (request *http.Request, err error) {
						return nil, errors.New("unknown error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &wrapperMiddleware{
				name:       tt.fields.name,
				Middleware: tt.fields.Middleware,
			}
			got, err := w.NewRequest(tt.args.rule, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapClientWithScheme(t *testing.T) {
	type args struct {
		name   string
		scheme string
		client *Client
	}
	tests := []struct {
		name string
		args args
		want giraffe.Client
	}{
		{
			name: "https",
			args: args{
				name:   "cmdb",
				scheme: "https",
				client: &Client{
					Client: &restv2.Client{
						Client: &http.Client{
							Transport: &transport{
								sig: ApiKey{
									AccessKey: "3fc93fed595063856df3ee1a",
									SecretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
								},
								rt: http.DefaultTransport,
							},
						},
						Middleware:  restv2.DefaultMiddleware,
						NameService: restv2.StaticAddress("192.168.100.162:8080"),
					},
					transportOptions: []TransportOption{RoundTripper(nil)},
				},
			},
			want: &wrapper{&restv2.Client{
				Client: &http.Client{
					Transport: &transport{
						sig: ApiKey{
							AccessKey: "3fc93fed595063856df3ee1a",
							SecretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
						},
						rt: http.DefaultTransport,
					},
				},
				Middleware:  &wrapperMiddleware{name: "cmdb", Middleware: restv2.DefaultMiddleware},
				NameService: restv2.StaticAddress("192.168.100.162:8080"),
				Scheme:      "https",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WrapClientWithScheme(tt.args.name, tt.args.scheme, tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WrapClientWithScheme() = %v, want %v", got, tt.want)
			}
		})
	}
}
