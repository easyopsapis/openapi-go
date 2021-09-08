package gerr

import (
	"github.com/easyops-cn/giraffe-micro/codes"
)

type StatusError interface {
	Error() string
	Status() Status
}

type statusError status

func (e *statusError) Error() string {
	m := (*status)(e).s
	switch {
	case m.CodeExplain != "":
		return m.CodeExplain
	case m.Error != "":
		return m.Error
	case m.Message != "":
		return m.Message
	case m.Msg != "":
		return m.Msg
	default:
		return codes.Code_name[int32(m.Code)]
	}
}

func (e *statusError) Status() Status {
	return (*status)(e)
}
