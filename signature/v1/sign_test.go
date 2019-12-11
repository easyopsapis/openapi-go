package signature

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type testReader struct {
	read func(p []byte) (n int, err error)
}

func (t testReader) Read(p []byte) (n int, err error) {
	return t.read(p)
}

func TestSignRequest(t *testing.T) {
	type args struct {
		accessKey string
		secretKey string
		t         time.Time
		req       *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				accessKey: "3fc93fed595063856df3ee1a",
				secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
				t:         time.Unix(1460314842, 0),
				req: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1/cmdb/object/list?page=1&pageSize=30", nil)
					return req
				}(),
			},
			want: "e01d1fee0425994caa85a9ff46e6ba1630cea4b7",
		},
		{
			args: args{
				accessKey: "3fc93fed595063856df3ee1a",
				secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
				t:         time.Unix(1460314842, 0),
				req: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1/cmdb/object/list", bytes.NewReader([]byte("{\"id\":123}")))
					req.Header.Add("Content-Type", "application/json")
					return req
				}(),
			},
			want: "0c24f5b2d88b056fca35e650047baab2047d7989",
		},
		{
			args: args{
				accessKey: "3fc93fed595063856df3ee1a",
				secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
				t:         time.Unix(1460314842, 0),
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
			args: args{
				accessKey: "3fc93fed595063856df3ee1a",
				secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
				t:         time.Unix(1460314842, 0),
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
		t.Run(tt.name, func(t *testing.T) {
			got, err := SignRequest(tt.args.accessKey, tt.args.secretKey, tt.args.t, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SignRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
