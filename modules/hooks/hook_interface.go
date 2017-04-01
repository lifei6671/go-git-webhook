// Package hooks 解析 Git WebHook 的内容信息
package hooks

// WebHookRequestInterface Git WebHook 请求内容解析接口.
type WebHookRequestInterface interface {

	// ServiceName 当前推送的服务名称.
	ServiceName()(string)

	// BeforeValue 获取推送钱的Hash值.
	BeforeValue()(string,error)

	// AfterValue 获取当前的Hash值.
	AfterValue()(string,error)

	// RepositoryName 获取仓库名称.
	RepositoryName()(string,error)

	// BranchName 获取分支名称.
	BranchName()(string,error)

	// UserName 获取用户名称.
	UserName()(string,error)

	// UserEmail 获取用户邮箱.
	UserEmail()(string,error)

	// DefaultBranch 默认分支.
	DefaultBranch()(string,error)

	// XPath 读取指定路径下的值： /project/git_ssh_url 则表示读取从根目录开始project下的git_ssh_url的值
	XPath(p string)(string,error)
}


