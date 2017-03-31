// package conf 为配置相关.
package conf

import "github.com/astaxie/beego"

// 登录用户的Session名
var LoginSessionName = "LoginSessionName"

// app_key
func GetAppKey()  (string) {
	return beego.AppConfig.DefaultString("app_key","go-git-webhook")
}
// 队列长度
func QueueSize() int {
	queueSize := beego.AppConfig.DefaultInt("queue_size",100)

	return queueSize
}