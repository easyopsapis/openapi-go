package openapi

import (
	"github.com/easyops-cn/giraffe-micro/plugins/restv2"
	"net/http"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		address   string
		accessKey string
		secretKey string
		options   []ClientOption
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			args: args{
				address:   "192.168.100.162:8080",
				accessKey: "3fc93fed595063856df3ee1a",
				secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
				options:   []ClientOption{WithClient(&http.Client{})},
			},
			want: &Client{
				Client: &restv2.Client{
					Client: &http.Client{
						Transport: &transport{
							accessKey: "3fc93fed595063856df3ee1a",
							secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
							rt:        http.DefaultTransport,
						},
					},
					Middleware:  restv2.DefaultMiddleware,
					NameService: restv2.StaticAddress("192.168.100.162:8080"),
				},
				transportOptions: []TransportOption{RoundTripper(nil)},
			},
		},
		{
			args: args{
				address:   "192.168.100.162:8080",
				accessKey: "3fc93fed595063856df3ee1a",
				secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
				options:   []ClientOption{WithClient(nil)},
			},
			want: &Client{
				Client: &restv2.Client{
					Client: &http.Client{
						Transport: &transport{
							accessKey: "3fc93fed595063856df3ee1a",
							secretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
							rt:        http.DefaultTransport,
						},
					},
					Middleware:  restv2.DefaultMiddleware,
					NameService: restv2.StaticAddress("192.168.100.162:8080"),
				},
				transportOptions: []TransportOption{RoundTripper(nil)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.address, tt.args.accessKey, tt.args.secretKey, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && !reflect.DeepEqual(got.Client, tt.want.Client) {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}
