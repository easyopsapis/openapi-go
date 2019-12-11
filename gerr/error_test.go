package gerr

import (
	"reflect"
	"testing"

	"github.com/easyops-cn/giraffe-micro/codes"
)

func Test_statusError_Error(t *testing.T) {
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
			e := &statusError{
				s: tt.fields.s,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("statusError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_statusError_Status(t *testing.T) {
	type fields struct {
		s Message
	}
	tests := []struct {
		name   string
		fields fields
		want   Status
	}{
		{
			name: "Test_HappyPath",
			fields: fields{
				s: Message{
					Code:        codes.Code_INTERNAL,
					CodeExplain: "codeExplain: don't worry be happy",
					Message:     "message: don't worry be happy",
					Msg:         "msg: don't worry be happy",
					Error:       "error: don't worry be happy",
				},
			},
			want: &status{
				s: Message{
					Code:        codes.Code_INTERNAL,
					CodeExplain: "codeExplain: don't worry be happy",
					Message:     "message: don't worry be happy",
					Msg:         "msg: don't worry be happy",
					Error:       "error: don't worry be happy",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &statusError{
				s: tt.fields.s,
			}
			if got := e.Status(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("statusError.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
