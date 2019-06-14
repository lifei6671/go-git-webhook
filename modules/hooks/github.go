package hooks

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	// ErrNotFound .
	ErrNotFound = errors.New("Value not found.")
)

// GitHubWebHook struct
type GitHubWebHook struct {
	data    string
	dataMap map[string]interface{}
}

func NewGitHubWebHook(value string) (*GitHubWebHook, error) {
	hook := &GitHubWebHook{}

	var dataMap map[string]interface{}

	if err := json.Unmarshal([]byte(value), &dataMap); err != nil {
		return hook, err
	}
	hook.data = value
	hook.dataMap = dataMap
	return hook, nil
}

// ServiceName 当前推送的服务名称.
func (p *GitHubWebHook) ServiceName() string {
	return "GitHub"
}

// BeforeValue 获取推送钱的Hash值.
func (p *GitHubWebHook) BeforeValue() (string, error) {

	return p.XPath("/before")
}

// AfterValue 获取当前的Hash值.
func (p *GitHubWebHook) AfterValue() (string, error) {

	return p.XPath("/after")
}

// RepositoryName 获取仓库名称.
func (p *GitHubWebHook) RepositoryName() (string, error) {
	return p.XPath("/repository/name")
}

// BranchName 获取分支名称.
func (p *GitHubWebHook) BranchName() (string, error) {

	value, err := p.XPath("/ref")

	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(value, "refs/heads/"), nil
}

// UserName 获取用户名称.
func (p *GitHubWebHook) UserName() (string, error) {
	return p.XPath("/pusher/name")
}

// UserEmail 获取用户邮箱.
func (p *GitHubWebHook) UserEmail() (string, error) {
	return p.XPath("/pusher/email")
}

// DefaultBranch 默认分支.
func (p *GitHubWebHook) DefaultBranch() (string, error) {
	return p.XPath("/repository/default_branch")
}

// XPath 读取指定路径下的值： /project/git_ssh_url 则表示读取从根目录开始project下的git_ssh_url的值
func (p *GitHubWebHook) XPath(xpath string) (string, error) {

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
