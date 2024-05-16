package openapi

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

type testReader struct {
	read func(p []byte) (n int, err error)
}

func (t testReader) Read(p []byte) (n int, err error) {
	return t.read(p)
}

func Test_request_SetSignature(t *testing.T) {
	type fields struct {
		Request *http.Request
	}
	type args struct {
		accessKey string
		signature string
		expires   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *http.Request
	}{
		{
			fields: fields{
				Request: &http.Request{
					URL: &url.URL{RawQuery: ""},
				},
			},
			args: args{
				accessKey: "3fc93fed595063856df3ee1a",
				signature: "e01d1fee0425994caa85a9ff46e6ba1630cea4b7",
				expires:   time.Unix(1460314842, 0),
			},
			want: &http.Request{
				URL: &url.URL{RawQuery: "accesskey=3fc93fed595063856df3ee1a&expires=1460314842&signature=e01d1fee0425994caa85a9ff46e6ba1630cea4b7"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request{
				Request: tt.fields.Request,
			}
			r.SetSignature(tt.args.accessKey, tt.args.signature, tt.args.expires)
			got := r.Request
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_request_StringToSign(t *testing.T) {
	type fields struct {
		Request *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "url部分存在中文、符号等需要url转码的内容",
			fields: fields{
				Request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1/cmdb/object/你好?page=1&pageSize=30", nil)
					return req
				}(),
			},
			want: "GET\n/cmdb/object/%E4%BD%A0%E5%A5%BD\npage1pageSize30\n\n",
		},
		{
			fields: fields{
				Request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1/cmdb/object/list?page=1&pageSize=30", nil)
					return req
				}(),
			},
			want: "GET\n/cmdb/object/list\npage1pageSize30\n\n",
		},
		{
			fields: fields{
				Request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1/cmdb/object/list", bytes.NewReader([]byte("{\"id\":123}")))
					return req
				}(),
			},
			want: "POST\n/cmdb/object/list\n\n\n07925d389335c0229b97393df477a438",
		},
		{
			fields: fields{
				Request: func() *http.Request {
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
				Request: func() *http.Request {
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
		t.Run(tt.name, func(t *testing.T) {
			r := request{
				Request: tt.fields.Request,
			}
			got, err := r.StringToSign()
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToSign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringToSign() got = %v, want %v", got, tt.want)
			}
		})
	}
}
