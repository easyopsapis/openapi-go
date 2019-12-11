package openapi

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_fieldMask(t *testing.T) {
	type exampleMessage struct {
		XXX_RestFieldMask []string `protobuf:"bytes,6,rep,name=XXX_RestFieldMask,json=XXXRestFieldMask,proto3" json:"XXX_RestFieldMask,omitempty"`
	}

	type messageWithWrongTypeFieldMask struct {
		XXX_RestFieldMask []int `protobuf:"bytes,6,rep,name=XXX_RestFieldMask,json=XXXRestFieldMask,proto3" json:"XXX_RestFieldMask,omitempty"`
	}

	type messageWithWrongTypeFieldMask2 struct {
		XXX_RestFieldMask string `protobuf:"bytes,6,rep,name=XXX_RestFieldMask,json=XXXRestFieldMask,proto3" json:"XXX_RestFieldMask,omitempty"`
	}

	type messageWithoutFieldMask struct {
	}
	type args struct {
		in interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test_HappyPath",
			args: args{
				in: &exampleMessage{
					XXX_RestFieldMask: []string{"a.b.c", "a.b", "a"},
				},
			},
			want: []string{"a.b.c", "a.b", "a"},
		},
		{
			name: "Test_HappyPath_WithWrongTypeFieldMask",
			args: args{
				in: &messageWithWrongTypeFieldMask{
					XXX_RestFieldMask: []int{1, 2, 3},
				},
			},
			want: nil,
		},
		{
			name: "Test_HappyPath_WithWrongTypeFieldMask2",
			args: args{
				in: &messageWithWrongTypeFieldMask2{
					XXX_RestFieldMask: "a.b.c",
				},
			},
			want: nil,
		},
		{
			name: "Test_HappyPath_EmptyFieldMask",
			args: args{
				in: &exampleMessage{
					XXX_RestFieldMask: []string{},
				},
			},
			want: []string{},
		},
		{
			name: "Test_HappyPath_NilFieldMask",
			args: args{
				in: &exampleMessage{
					XXX_RestFieldMask: nil,
				},
			},
			want: []string{},
		},
		{
			name: "Test_HappyPath_NilFieldMask",
			args: args{
				in: &messageWithoutFieldMask{},
			},
			want: nil,
		},
		{
			name: "Test_WithNilMessage",
			args: args{
				in: (*exampleMessage)(nil),
			},
			want: nil,
		},
		{
			name: "Test_WithNil",
			args: args{
				in: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fieldMask(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFieldMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withFieldMask(t *testing.T) {
	type args struct {
		req       *http.Request
		fieldMask []string
	}
	tests := []struct {
		name string
		args args
		want *http.Request
	}{
		{
			name: "Test_HappyPath",
			args: args{
				req: func() *http.Request {
					r, _ := http.NewRequest("GET", "/object?q1=1&q2=2", nil)
					return r
				}(),
				fieldMask: []string{"a.b", "q2"},
			},
			want: func() *http.Request {
				r, _ := http.NewRequest("GET", "/object?q2=2", nil)
				return r
			}(),
		},
		{
			name: "field mask empty",
			args: args{
				req: func() *http.Request {
					r, _ := http.NewRequest("GET", "/object?q1=1&q2=2", nil)
					return r
				}(),
				fieldMask: []string{},
			},
			want: func() *http.Request {
				r, _ := http.NewRequest("GET", "/object?q1=1&q2=2", nil)
				return r
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withFieldMask(tt.args.req, tt.args.fieldMask)
		})
	}
}
