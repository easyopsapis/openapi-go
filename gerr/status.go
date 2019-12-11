package gerr

import (
	"github.com/easyops-cn/giraffe-micro/codes"
)

type Status interface {
	Code() codes.Code
	CodeExplain() string
	Message() string
	Err() error
}

type status struct {
	s Message
}

func (s *status) Code() codes.Code {
	return s.s.Code
}

func (s *status) CodeExplain() string {
	switch {
	case s.s.CodeExplain != "":
		return s.s.CodeExplain
	case s.s.Error != "":
		return s.s.Error
	case s.s.Message != "":
		return s.s.Message
	case s.s.Msg != "":
		return s.s.Msg
	default:
		return codes.Code_name[int32(s.s.Code)]
	}
}

func (s *status) Message() string {
	switch {
	case s.s.Message != "":
		return s.s.Message
	case s.s.Msg != "":
		return s.s.Msg
	case s.s.Error != "":
		return s.s.Error
	case s.s.CodeExplain != "":
		return s.s.CodeExplain
	default:
		return codes.Code_name[int32(s.s.Code)]
	}
}

func (s *status) Err() error {
	if s.s.Code == codes.Code_OK {
		return nil
	}
	return (*statusError)(s)
}

func FromError(err error) Status {
	if err == nil {
		return NewStatus(codes.Code_OK, "", "")
	}
	switch v := err.(type) {
	case interface{ Status() Status }:
		return v.Status()
	default:
		return NewStatus(codes.Code_UNKNOWN, err.Error(), "")
	}
}

func ErrorProto(s *Message) error {
	return FromProto(s).Err()
}

func FromProto(s *Message) Status {
	return &status{s: *s}
}

func NewStatus(code codes.Code, codeExplain string, message string) Status {
	return &status{
		s: Message{
			Code:        code,
			CodeExplain: codeExplain,
			Message:     message,
		},
	}
}
