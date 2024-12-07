package ecode

import (
	"strconv"
)

type ErrorX struct {
	*Errno
	Ext        error  `json:"-"`
	ExtMessage string `json:"ext_msg,omitempty"`
}

func (e *ErrorX) Error() string {
	if e.Ext == nil {
		return "Err - code: " + strconv.Itoa(e.Code) + ", msg: " + e.Msg
	}
	return "Err - code: " + strconv.Itoa(e.Code) + ", msg: " + e.Msg + ", ext: " + e.Ext.Error()
}

func Wrap(errno *Errno, err error) *ErrorX {
	return &ErrorX{
		Errno: errno,
		Ext:   err,
	}
}
