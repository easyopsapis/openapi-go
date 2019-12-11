package openapi_test

import (
	"bytes"
	"errors"
	"github.com/easyopsapis/openapi-go/v1"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/easyops-cn/giraffe-micro"
	"github.com/easyops-cn/giraffe-micro/plugins/restv2"
	giraffeproto "github.com/easyops-cn/go-proto-giraffe"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
)

type modifyMessage struct {
	ID                string   `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Name              string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Desc              string   `protobuf:"bytes,4,opt,name=desc,proto3" json:"desc,omitempty"`
	XXX_RestFieldMask []string `protobuf:"bytes,6,rep,name=XXX_RestFieldMask,json=XXXRestFieldMask,proto3" json:"XXX_RestFieldMask,omitempty"`
}

func (m *modifyMessage) Reset()         { *m = modifyMessage{} }
func (m *modifyMessage) String() string { return proto.CompactTextString(m) }
func (*modifyMessage) ProtoMessage()    {}

func Test_middleware_NewRequest(t *testing.T) {
	type fields struct {
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
			name: "Test_HappyPath",
			fields: fields{
				Middleware: restv2.DefaultMiddleware,
			},
			args: args{
				rule: &giraffeproto.HttpRule{
					Pattern: &giraffeproto.HttpRule_Get{Get: "/test"},
					Body:    "",
				},
				in: &modifyMessage{
					ID:                "123",
					Name:              "abc",
					XXX_RestFieldMask: []string{"name"},
				},
			},
			want: func() *http.Request {
				r, _ := http.NewRequest("GET", "/test?name=abc", nil)
				return r
			}(),
			wantErr: false,
		},
		{
			name: "Test_WithoutRule",
			fields: fields{
				Middleware: restv2.DefaultMiddleware,
			},
			args: args{
				in: &modifyMessage{
					ID:                "123",
					Name:              "abc",
					XXX_RestFieldMask: []string{"name"},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &openapi.Middleware{
				Middleware: tt.fields.Middleware,
			}
			got, err := m.NewRequest(tt.args.rule, tt.args.in)
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

func newResponse(code int, body []byte) *http.Response {
	rec := httptest.NewRecorder()
	rec.Code = code
	rec.Body = bytes.NewBuffer(body)
	return rec.Result()
}

type errReadCloser struct{}

func (i errReadCloser) Read(p []byte) (n int, err error) { return 0, errors.New("always failed") }
func (i errReadCloser) Close() error                     { return errors.New("always failed") }

func Test_middleware_ParseResponse(t *testing.T) {
	type fields struct {
		Middleware restv2.Middleware
	}
	type args struct {
		rule giraffe.HttpRule
		resp *http.Response
		out  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//TODO 增加对 out 的断言
		wantErr bool
	}{
		{
			name:   "Test_HappyPath",
			fields: fields{Middleware: restv2.DefaultMiddleware},
			args: args{
				rule: &giraffeproto.HttpRule{
					Pattern: &giraffeproto.HttpRule_Get{Get: "/test"},
					Body:    "data",
				},
				resp: newResponse(http.StatusOK, []byte("{\"code\":0,\"codeExplain\":\"ok\",\"data\":{}}")),
				out:  &types.Struct{},
			},
			wantErr: false,
		},
		{
			name:   "Test_WithError",
			fields: fields{Middleware: restv2.DefaultMiddleware},
			args: args{
				rule: &giraffeproto.HttpRule{
					Pattern: &giraffeproto.HttpRule_Get{Get: "/test"},
					Body:    "data",
				},
				resp: newResponse(http.StatusInternalServerError, []byte("{\"code\":100004,\"codeExplain\":\"ok\",\"data\":{}}")),
				out:  &types.Struct{},
			},
			wantErr: true,
		},
		{
			name:   "Test_ErrorReadCloser",
			fields: fields{Middleware: restv2.DefaultMiddleware},
			args: args{
				rule: &giraffeproto.HttpRule{
					Pattern: &giraffeproto.HttpRule_Get{Get: "/test"},
					Body:    "data",
				},
				resp: &http.Response{StatusCode: http.StatusInternalServerError, Body: &errReadCloser{}},
				out:  &types.Struct{},
			},
			wantErr: true,
		},
		{
			name:   "Test_ErrorReadCloser",
			fields: fields{Middleware: restv2.DefaultMiddleware},
			args: args{
				rule: &giraffeproto.HttpRule{
					Pattern: &giraffeproto.HttpRule_Get{Get: "/test"},
					Body:    "data",
				},
				resp: newResponse(http.StatusInternalServerError, []byte("")),
				out:  &types.Struct{},
			},
			wantErr: true,
		},
		{
			name:   "Test_HasNotErrorCode",
			fields: fields{Middleware: restv2.DefaultMiddleware},
			args: args{
				rule: &giraffeproto.HttpRule{
					Pattern: &giraffeproto.HttpRule_Get{Get: "/test"},
				},
				resp: newResponse(http.StatusInternalServerError, []byte("{}")),
				out:  &types.Struct{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &openapi.Middleware{
				Middleware: tt.fields.Middleware,
			}
			if err := m.ParseResponse(tt.args.rule, tt.args.resp, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("ParseResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
