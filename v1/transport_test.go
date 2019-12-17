package openapi

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type testRoundTripper struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (t *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.roundTrip(req)
}

type testSigner struct {
	sign func(expires time.Time, request Request) error
}

func (t testSigner) Sign(expires time.Time, request Request) error {
	return t.sign(expires, request)
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
				rt: http.DefaultTransport,
				sig: ApiKey{
					AccessKey: "3fc93fed595063856df3ee1a",
					SecretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
				},
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

func Test_transport_RoundTrip(t1 *testing.T) {
	type fields struct {
		rt  http.RoundTripper
		sig Signer
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
				rt: &testRoundTripper{roundTrip: func(h *http.Request) (response *http.Response, err error) {
					return nil, nil
				}},
				sig: &testSigner{sign: func(expires time.Time, request Request) error {
					return nil
				}},
			},
			args: args{req: &http.Request{}},
		},
		{
			fields: fields{
				rt: &testRoundTripper{roundTrip: func(h *http.Request) (response *http.Response, err error) {
					return nil, nil
				}},
				sig: &testSigner{sign: func(expires time.Time, request Request) error {
					return errors.New("unknown error")
				}},
			},
			args:    args{req: &http.Request{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &transport{
				rt:  tt.fields.rt,
				sig: tt.fields.sig,
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
