package models

import "errors"

var(
	//用户不存在
	ErrMemberNoExist = errors.New("用户不存在")
	//密码错误
	ErrorMemberPasswordError = errors.New("用户密码错误")
	//指定的服务已存在
	ErrServerAlreadyExist = errors.New("服务已存在")
)
