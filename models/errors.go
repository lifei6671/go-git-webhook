package models

import "errors"

var(
	ErrMemberNoExist = errors.New("member deos not exist")
	ErrorMemberPasswordError = errors.New("member password error")
	ErrServerAlreadyExist = errors.New("Server already exists")
)
