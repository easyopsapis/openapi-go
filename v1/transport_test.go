package openapi

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type testRoundTripper struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (t *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.roundTrip(req)
}

type testReader struct {
	read func(p []byte) (n int, err error)
}

func (t testReader) Read(p []byte) (n int, err error) {
	return t.read(p)
}

func Test_transport_RoundTrip(t1 *testing.T) {
	type fields struct {
		rt        http.RoundTripper
		accessKey string
		secretKey string
		expires   func() string
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			fields: fields{
				rt: &testRoundTripper{
					roundTrip: func(request *http.Request) (response *http.Response, err error) {
						query := request.URL.Query()
						want := "e01d1fee0425994caa85a9ff46e6ba1630cea4b7"
						got := query.Get("signature")
						if !reflect.DeepEqual(got, want) {
							t1.Errorf("want %s, got %s", want, got)
						}
						return nil, nil
					},
				},
				accessKey: "3fc93fed595063856df3ee1a",
				secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
				expires: func() string {
					return "1460314842"
				},
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1/cmdb/object/list?page=1&pageSize=30", nil)
					return req
				}(),
			},
		},
		{
			fields: fields{
				rt: &testRoundTripper{
					roundTrip: func(request *http.Request) (response *http.Response, err error) {
						query := request.URL.Query()
						want := "0c24f5b2d88b056fca35e650047baab2047d7989"
						got := query.Get("signature")
						if !reflect.DeepEqual(got, want) {
							t1.Errorf("want %s, got %s", want, got)
						}
						return nil, nil
					},
				},
				accessKey: "3fc93fed595063856df3ee1a",
				secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
				expires: func() string {
					return "1460314842"
				},
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1/cmdb/object/list", bytes.NewReader([]byte("{\"id\":123}")))
					req.Header.Add("Content-Type", "application/json")
					return req
				}(),
			},
		},
		{
			fields: fields{
				rt: &testRoundTripper{
					roundTrip: func(request *http.Request) (response *http.Response, err error) {
						t1.Errorf("should not call")
						return nil, nil
					},
				},
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1/cmdb/object/list", bytes.NewReader([]byte("{\"id\":123}")))
					req.Header.Add("Content-Type", "application/json")
					req.GetBody = func() (closer io.ReadCloser, err error) {
						return nil, errors.New("unknown error")
					}
					return req
				}(),
			},
			wantErr: true,
		},
		{
			fields: fields{
				rt: &testRoundTripper{
					roundTrip: func(request *http.Request) (response *http.Response, err error) {
						t1.Errorf("should not call")
						return nil, nil
					},
				},
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1/cmdb/object/list", bytes.NewReader([]byte("{\"id\":123}")))
					req.Header.Add("Content-Type", "application/json")
					req.GetBody = func() (closer io.ReadCloser, err error) {
						return ioutil.NopCloser(&testReader{read: func(p []byte) (n int, err error) {
							return 0, errors.New("unknown error")
						}}), nil
					}
					return req
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &transport{
				rt:        tt.fields.rt,
				accessKey: tt.fields.accessKey,
				secretKey: tt.fields.secretKey,
				expires:   tt.fields.expires,
			}
			got, err := t.RoundTrip(tt.args.req)
			if (err != nil) != tt.wantErr {
				t1.Errorf("RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("RoundTrip() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTransport(t *testing.T) {
	type args struct {
		accessKey string
		secretKey string
		options   []TransportOption
	}
	tests := []struct {
		name string
		args args
		want http.RoundTripper
	}{
		{
			args: args{
				accessKey: "3fc93fed595063856df3ee1a",
				secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
				options:   []TransportOption{RoundTripper(http.DefaultTransport)},
			},
			want: &transport{
				rt:        http.DefaultTransport,
				accessKey: "3fc93fed595063856df3ee1a",
				secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := NewTransport(tt.args.accessKey, tt.args.secretKey, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransport() = %v, want %v", got, tt.want)
			}
		})
	}
}
