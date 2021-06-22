package biz

import (
	"errors"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase, NewUserUsecase)


var (
	ErrRecordNotFound = errors.New("record not fount")
	ErrRecordAlreadyExist = errors.New("record already exists")
	ErrUserIDGen = errors.New("generate user id failed")
	ErrPasswordWrong = errors.New("wrong password")
)