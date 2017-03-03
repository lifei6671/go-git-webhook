package models

import (
	"errors"
	"github.com/widuu/gojson"
	"strings"
)

var(
	ErrWebHookDataResolveFailure = errors.New("web hook data resolve failure")
	ErrRepositoryNameNoExist = errors.New("repository name does not exist")
	ErrBranchNameNoExist = errors.New("branch name does not exist")
	ErrPushUserNoExist = errors.New("push user does not exist")
	ErrNoData = errors.New("data does no exist")
)

type HookData struct {
	data string
}

func (m *HookData) Json () *gojson.Js {
	return gojson.Json(m.data);
}

func ResolveHookRequest (data string) (HookData,error) {

	hookData := HookData{}

	result := gojson.Json(data);

	if !result.IsValid()  {
		return hookData, ErrWebHookDataResolveFailure
	}

	hookData.data = data

	return hookData,nil
}

//解析仓库名称
func (m *HookData) RepositoryName() (string,error){

	result := m.Json().Get("repository").Get("name")

	if name:= result.Tostring() ; name != "" && result.IsValid() {
		return  name,nil
	}

	result = m.Json().Get("push_data").Get("name")

	if name:= result.Tostring() ; name != "" && result.IsValid() {
		return  name,nil
	}

	return "",ErrRepositoryNameNoExist
}

// 获取分支名称
func (m *HookData) BranchName () (string,error){

	//github, gitlib
	branch  := m.Json().Get("ref")

	if !branch.IsValid() {
		branch := m.Json().Get("push_data").Get("ref")

		if !branch.IsValid() {
			return "",ErrBranchNameNoExist
		}
	}

	branchName := branch.Tostring();

	index := strings.Index(branchName,"/")

	if index > 0 {
		branchName = string([]byte(branchName)[index+1:])
	}
	return branchName,nil
}

// 获取推送者
func (m *HookData) PushUser () (string,error){

	//github的data格式
	if uid := m.Json().Get("pusher").Get("name");uid.IsValid() {
		return uid.Tostring(),nil
	}

	//gitlib 格式
	if uid := m.Json().Get("user_name");uid.IsValid() {
		return uid.Tostring(),nil;
	}

	//gogs 格式
	if uid := m.Json().Get("pusher").Get("username");uid.IsValid() {
		return uid.Tostring(),nil;
	}

	//gitosc的data格式
	if uid := m.Json().Get("push_data").Get("user").Get("name");uid.IsValid() {
		return  uid.Tostring(),nil
	}
	return "",ErrPushUserNoExist
}

// 获取推送者邮箱
func (m *HookData) PushEmail() (string,error) {

	//github的data格式
	if email := m.Json().Get("pusher").Get("email");email.IsValid() {
		return email.Tostring(),nil
	}

	//gitlib 格式
	if email := m.Json().Get("user_email");email.IsValid() {
		return email.Tostring(),nil
	}

	//gitosc的data格式
	if email := m.Json().Get("push_data");email.IsValid() {
		return email.Tostring(),nil
	}

	return "",ErrNoData
}























