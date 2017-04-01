package hooks

import (
	"strings"
	"encoding/json"
)

// GogsWebHook struct .
type GogsWebHook struct {
	dataMap map[string]interface{}
	data string
}


// NewGogsWebHook 创建一个对象.
func NewGogsWebHook(value string) (*GogsWebHook,error) {
	hook := &GogsWebHook{}

	var dataMap map[string]interface{}

	if err := json.Unmarshal([]byte(value),&dataMap);err != nil {
		return hook,err
	}
	hook.data = value
	hook.dataMap = dataMap
	return hook,nil
}
// ServiceName 当前推送的服务名称.
func (p *GogsWebHook) ServiceName() (string)  {
	return "GitLab"
}

// BeforeValue 获取推送钱的Hash值.
func (p *GogsWebHook) BeforeValue() (string,error)  {

	return p.XPath("/before")
}

// AfterValue 获取当前的Hash值.
func (p *GogsWebHook) AfterValue() (string,error)  {

	return p.XPath("/after")
}

// RepositoryName 获取仓库名称.
func (p *GogsWebHook) RepositoryName()(string,error){
	return p.XPath("/repository/name")
}

// BranchName 获取分支名称.
func (p *GogsWebHook) BranchName()(string,error) {

	value,err := p.XPath("/ref")

	if err != nil {
		return "",err
	}

	return strings.TrimLeft(value,"refs/heads/"),nil
}

// UserName 获取用户名称.
func (p *GogsWebHook) UserName()(string,error){
	return p.XPath("/pusher/username")
}

// UserEmail 获取用户邮箱.
func (p *GogsWebHook) UserEmail() (string,error){
	return p.XPath("/pusher/email")
}

// DefaultBranch 默认分支.
func (p *GogsWebHook)DefaultBranch()(string,error){
	return p.XPath("/repository/default_branch")
}

// XPath 读取指定路径下的值： /project/git_ssh_url 则表示读取从根目录开始project下的git_ssh_url的值
func (p *GogsWebHook) XPath(xpath string)(string,error) {

	if strings.Trim(xpath," ") == "" {
		return "",ErrNotFound
	}
	paths := strings.Split(xpath,"/")

	if len(paths) <= 0 {
		return "",ErrNotFound
	}

	dataMap := p.dataMap

	for _,key := range paths  {
		if key == ""{
			continue
		}
		if data,ok := dataMap[key];ok {
			if data1,ok := data.(map[string]interface{}); ok {
				dataMap = data1
			}
		}
	}

	lastKey := paths[len(paths)-1]


	if v,ok := dataMap[lastKey];ok {

		if value,ok := v.(string);ok {
			return value,nil
		}
	}
	return "",ErrNotFound
}