package hooks

import (
	"encoding/json"
	"strings"
)

// GitLabWebHook struct.
type GitLabWebHook struct {
	dataMap map[string]interface{}
	data    string
}

// NewGitLabWebHook 创建一个对象.
func NewGitLabWebHook(value string) (*GitLabWebHook, error) {
	hook := &GitLabWebHook{}

	var dataMap map[string]interface{}

	if err := json.Unmarshal([]byte(value), &dataMap); err != nil {
		return hook, err
	}
	hook.data = value
	hook.dataMap = dataMap
	return hook, nil
}

// ServiceName 当前推送的服务名称.
func (p *GitLabWebHook) ServiceName() string {
	return "GitLab"
}

// BeforeValue 获取推送钱的Hash值.
func (p *GitLabWebHook) BeforeValue() (string, error) {

	return p.XPath("/before")
}

// AfterValue 获取当前的Hash值.
func (p *GitLabWebHook) AfterValue() (string, error) {

	return p.XPath("/after")
}

// RepositoryName 获取仓库名称.
func (p *GitLabWebHook) RepositoryName() (string, error) {
	return p.XPath("/project/name")
}

// BranchName 获取分支名称.
func (p *GitLabWebHook) BranchName() (string, error) {

	value, err := p.XPath("/ref")

	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(value, "refs/heads/"), nil
}

// UserName 获取用户名称.
func (p *GitLabWebHook) UserName() (string, error) {
	return p.XPath("/user_name")
}

// UserEmail 获取用户邮箱.
func (p *GitLabWebHook) UserEmail() (string, error) {
	return p.XPath("/user_email")
}

// DefaultBranch 默认分支.
func (p *GitLabWebHook) DefaultBranch() (string, error) {
	return p.XPath("/project/default_branch")
}

// XPath 读取指定路径下的值： /project/git_ssh_url 则表示读取从根目录开始project下的git_ssh_url的值
func (p *GitLabWebHook) XPath(xpath string) (string, error) {

	if strings.Trim(xpath, " ") == "" {
		return "", ErrNotFound
	}
	paths := strings.Split(xpath, "/")

	if len(paths) <= 0 {
		return "", ErrNotFound
	}

	dataMap := p.dataMap

	for _, key := range paths {
		if key == "" {
			continue
		}
		if data, ok := dataMap[key]; ok {
			if data1, ok := data.(map[string]interface{}); ok {
				dataMap = data1
			}
		}
	}

	lastKey := paths[len(paths)-1]

	if v, ok := dataMap[lastKey]; ok {

		if value, ok := v.(string); ok {
			return value, nil
		}
	}
	return "", ErrNotFound
}
