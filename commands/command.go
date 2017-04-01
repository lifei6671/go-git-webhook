// package commands 为命令行定义.
package commands

import (
	"net/url"
	"time"
	"fmt"

	"github.com/lifei6671/go-git-webhook/tasks"
	"github.com/lifei6671/go-git-webhook/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

// RegisterDataBase 注册数据库
func RegisterDataBase()  {
	host := beego.AppConfig.String("db_host")
	database := beego.AppConfig.String("db_database")
	username := beego.AppConfig.String("db_username")
	password := beego.AppConfig.String("db_password")
	timezone := beego.AppConfig.String("timezone")

	port := beego.AppConfig.String("db_port")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s",username,password,host,port,database,url.QueryEscape(timezone))


	orm.RegisterDataBase("default", "mysql", dataSource)

	orm.DefaultTimeLoc, _ = time.LoadLocation(timezone)
}

// RegisterModel 注册Model
func RegisterModel()  {
	orm.RegisterModel(new(models.Member))
	orm.RegisterModel(new(models.Server))
	orm.RegisterModel(new(models.WebHook))
	orm.RegisterModel(new(models.Scheduler))
	orm.RegisterModel(new(models.Relation))
}
// RegisterLogger 注册日志
func RegisterLogger()  {

	logs.SetLogger("console")
	logs.SetLogger("file",`{"filename":"logs/log.log"}`)
	logs.EnableFuncCallDepth(true)
	logs.Async()
}
// RegisterTaskQueue 注册队列
func RegisterTaskQueue()  {

	schedulerList,err := models.NewScheduler().QuerySchedulerByState("wait");
	if err == nil {
		for _,scheduler := range schedulerList {
			tasks.Add(tasks.Task{SchedulerId: scheduler.SchedulerId,WebHookId:scheduler.WebHookId,ServerId:scheduler.ServerId})
		}
	}else{
		fmt.Println(err)
	}

}
// RunCommand 注册orm命令行工具
func RunCommand()  {
	orm.RunCommand()
	Install()
	Version()
}

// Run 启动Web监听
func Run()  {
	beego.Run()
}
