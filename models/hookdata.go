package models

import (
	"errors"
	"strings"

	"github.com/widuu/gojson"
)

// 错误信息
var(
	ErrWebHookDataResolveFailure = errors.New("web hook data resolve failure")
	ErrRepositoryNameNoExist = errors.New("repository name does not exist")
	ErrBranchNameNoExist = errors.New("branch name does not exist")
	ErrPushUserNoExist = errors.New("push user does not exist")
	ErrNoData = errors.New("data does no exist")
	ErrInvalidParameter = errors.New("Invalid parameter")
)

// HookData Git请求时发送的信息
type HookData struct {
	data string
}

// Json 获取Json结构
func (m *HookData) Json () *gojson.Js {
	return gojson.Json(m.data);
}

// ResolveHookRequest 解析GitWebHook请求时发送的内容
func ResolveHookRequest (data string) (HookData,error) {

	hookData := HookData{}

	result := gojson.Json(data);

	if !result.IsValid()  {
		return hookData, ErrWebHookDataResolveFailure
	}

	hookData.data = data

	return hookData,nil
}

// RepositoryName 解析仓库名称
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

// BranchName 获取分支名称
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

	index := strings.IndexAny(branchName,"/")

	if index > 0 {
		branchName = branchName[index+1:]
	}
	return branchName,nil
}
// HookType 获取Git类型
func (m *HookData) HookType() (string,error){
	//github的data格式
	if uid := m.Json().Get("pusher").Get("name");uid.IsValid() {
		return "github",nil
	}

	//gitlab 格式
	if uid := m.Json().Get("user_name");uid.IsValid() {
		return "gitlab",nil;
	}

	//gogs 格式
	if uid := m.Json().Get("pusher").Get("username");uid.IsValid() {
		return "gogs",nil;
	}

	//gitosc的data格式
	if uid := m.Json().Get("push_data").Get("user").Get("name");uid.IsValid() {
		return  "gitosc",nil
	}
	return "",ErrPushUserNoExist
}
// PushUser 获取推送者
func (m *HookData) PushUser () (string,error){

	//github的data格式
	if uid := m.Json().Get("pusher").Get("name");uid.IsValid() {
		return uid.Tostring(),nil
	}

	//gitlab 格式
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

// PushEmail 获取推送者邮箱
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
// PushSha 获取当前推送的SHA值
func (m *HookData) PushSha() (string,error){
	if value := m.Json().Get("after"); value.IsValid() {
		return value.Tostring(),nil
	}
	return "",ErrNoData
}






















