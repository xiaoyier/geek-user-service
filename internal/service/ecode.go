package service

import (
	"github.com/go-kratos/kratos/v2/errors"
	"strconv"
)

type ResponseCode int32

const (
	CodeSucc ResponseCode = iota
	ErrCodeUserExisted ResponseCode = iota + 1000
	ErrCodeUserNotFound
	ErrCodePasswordWrong
	ErrCodeInternalError
)

var CodeMessageMap = map[ResponseCode]string{
	CodeSucc: "success",
	ErrCodeUserExisted: "user already existed",
	ErrCodeUserNotFound: "user not found",
	ErrCodePasswordWrong: "wrong password",
	ErrCodeInternalError: "internal error",
}

func (r ResponseCode) Message() string {
	if s, ok := CodeMessageMap[r]; ok {
		return s
	}
	return "server busy"
}

func (r ResponseCode) Reason() string {
	return strconv.Itoa(int(r))
}

func (r ResponseCode) ReasonAndMessage() (string,string) {
	return r.Reason(), r.Message()
}

func InternalServerWithCode(code ResponseCode) error {
	reason := code.Reason()
	message := code.Message()
	return errors.InternalServer(reason, message)
}

func NotFoundWithCode(code ResponseCode) error {
	reason := code.Reason()
	message := code.Message()
	return errors.NotFound(reason, message)
}