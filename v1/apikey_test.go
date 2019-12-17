package openapi

import (
	"errors"
	"testing"
	"time"
)

type testRequest struct {
	setSignature func(accessKey, signature string, expires time.Time)
	stringToSign func() (string, error)
}

func (t testRequest) SetSignature(accessKey, signature string, expires time.Time) {
	t.setSignature(accessKey, signature, expires)
}

func (t testRequest) StringToSign() (string, error) {
	return t.stringToSign()
}

func TestApiKey_Sign(t *testing.T) {
	type fields struct {
		AccessKey string
		SecretKey string
	}
	type args struct {
		expires time.Time
		r       Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				AccessKey: "3fc93fed595063856df3ee1a",
				SecretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
			},
			args: args{
				expires: time.Unix(1460314842, 0),
				r: &testRequest{
					setSignature: func(accessKey, signature string, expires time.Time) {
						want := "e01d1fee0425994caa85a9ff46e6ba1630cea4b7"
						if signature != want {
							t.Errorf("want %v, got %v", want, signature)
						}
					},
					stringToSign: func() (s string, err error) {
						return "GET\n/cmdb/object/list\npage1pageSize30\n\n", nil
					},
				},
			},
		},
		{
			fields: fields{
				AccessKey: "3fc93fed595063856df3ee1a",
				SecretKey: "1e338744a33426b3394e0ae9cd45af9c4e0d5fee5aad497e969cd21c65963d36",
			},
			args: args{
				expires: time.Unix(1460314842, 0),
				r: &testRequest{
					setSignature: func(accessKey, signature string, expires time.Time) {},
					stringToSign: func() (s string, err error) {
						return "", errors.New("unknown error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ApiKey{
				AccessKey: tt.fields.AccessKey,
				SecretKey: tt.fields.SecretKey,
			}
			if err := a.Sign(tt.args.expires, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
