package gerr

import (
	"errors"
	"reflect"
	"testing"

	"github.com/easyops-cn/giraffe-micro/codes"
)

func Test_status_Code(t *testing.T) {
	type fields struct {
		s Message
	}
	tests := []struct {
		name   string
		fields fields
		want   codes.Code
	}{
		{
			name: "Test_HappyPath",
			fields: fields{
				s: Message{
					Code: codes.Code_ALREADY_EXISTS,
				},
			},
			want: codes.Code_ALREADY_EXISTS,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &status{
				s: tt.fields.s,
			}
			if got := s.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("status.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_status_CodeExplain(t *testing.T) {
	type fields struct {
		s Message
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test_HappyPath_CodeExplain",
			fields: fields{
				s: Message{
					Code:        codes.Code_INTERNAL,
					CodeExplain: "codeExplain: don't worry be happy",
					Message:     "message: don't worry be happy",
					Msg:         "msg: don't worry be happy",
					Error:       "error: don't worry be happy",
				},
			},
			want: "codeExplain: don't worry be happy",
		},
		{
			name: "Test_HappyPath_Error",
			fields: fields{
				s: Message{
					Code:    codes.Code_INTERNAL,
					Message: "message: don't worry be happy",
					Msg:     "msg: don't worry be happy",
					Error:   "error: don't worry be happy",
				},
			},
			want: "error: don't worry be happy",
		},
		{
			name: "Test_HappyPath_Message",
			fields: fields{
				s: Message{
					Code:    codes.Code_INTERNAL,
					Message: "message: don't worry be happy",
					Msg:     "msg: don't worry be happy",
				},
			},
			want: "message: don't worry be happy",
		},
		{
			name: "Test_HappyPath_Msg",
			fields: fields{
				s: Message{
					Code: codes.Code_INTERNAL,
					Msg:  "msg: don't worry be happy",
				},
			},
			want: "msg: don't worry be happy",
		},
		{
			name: "Test_HappyPath_CodeName",
			fields: fields{
				s: Message{
					Code: codes.Code_INTERNAL,
				},
			},
			want: codes.Code_name[int32(codes.Code_INTERNAL)],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &status{
				s: tt.fields.s,
			}
			if got := s.CodeExplain(); got != tt.want {
				t.Errorf("status.CodeExplain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_status_Message(t *testing.T) {
	type fields struct {
		s Message
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test_HappyPath_Message",
			fields: fields{
				s: Message{
					Code:        codes.Code_INTERNAL,
					CodeExplain: "codeExplain: don't worry be happy",
					Message:     "message: don't worry be happy",
					Msg:         "msg: don't worry be happy",
					Error:       "error: don't worry be happy",
				},
			},
			want: "message: don't worry be happy",
		},
		{
			name: "Test_HappyPath_Msg",
			fields: fields{
				s: Message{
					Code:        codes.Code_INTERNAL,
					CodeExplain: "codeExplain: don't worry be happy",
					Msg:         "msg: don't worry be happy",
					Error:       "error: don't worry be happy",
				},
			},
			want: "msg: don't worry be happy",
		},
		{
			name: "Test_HappyPath_Error",
			fields: fields{
				s: Message{
					Code:        codes.Code_INTERNAL,
					CodeExplain: "codeExplain: don't worry be happy",
					Error:       "error: don't worry be happy",
				},
			},
			want: "error: don't worry be happy",
		},
		{
			name: "Test_HappyPath_CodeExplain",
			fields: fields{
				s: Message{
					Code:        codes.Code_INTERNAL,
					CodeExplain: "codeExplain: don't worry be happy",
				},
			},
			want: "codeExplain: don't worry be happy",
		},
		{
			name: "Test_HappyPath_CodeName",
			fields: fields{
				s: Message{
					Code: codes.Code_INTERNAL,
				},
			},
			want: codes.Code_name[int32(codes.Code_INTERNAL)],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &status{
				s: tt.fields.s,
			}
			if got := s.Message(); got != tt.want {
				t.Errorf("status.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_status_Err(t *testing.T) {
	type fields struct {
		s Message
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test_HappyPath_OK",
			fields: fields{
				s: Message{
					Code: codes.Code_OK,
				},
			},
			wantErr: false,
		},
		{
			name: "Test_HappyPath_NotFound",
			fields: fields{
				s: Message{
					Code: codes.Code_NOT_FOUND,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &status{
				s: tt.fields.s,
			}
			if err := s.Err(); (err != nil) != tt.wantErr {
				t.Errorf("status.Err() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFromError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Status
	}{
		{
			name: "Test_HappyPath",
			args: args{
				err: (*statusError)(&status{}),
			},
			want: &status{},
		},
		{
			name: "Test_WithNil",
			args: args{
				err: nil,
			},
			want: &status{
				s: Message{
					Code: codes.Code_OK,
				},
			},
		},
		{
			name: "Test_NormalError",
			args: args{
				err: errors.New("normal error"),
			},
			want: &status{
				s: Message{
					Code:        codes.Code_UNKNOWN,
					CodeExplain: "normal error",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorProto(t *testing.T) {
	type args struct {
		s *Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_HappyPath",
			args: args{
				s: &Message{
					Code:        0,
					CodeExplain: "OK",
				},
			},
			wantErr: false,
		},
		{
			name: "Test_Error",
			args: args{
				s: &Message{
					Code: codes.Code_ABORTED,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ErrorProto(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ErrorProto() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
