package routers

import (
	"github.com/astaxie/beego"
	"github.com/lifei6671/go-git-webhook/controllers"
)

func init()  {
	beego.Router("/", &controllers.HomeController{},"*:Index")

	beego.Router("/hook/edit/?:id", &controllers.HomeController{},"*:Edit")
	beego.Router("/hook/delete", &controllers.HomeController{},"post:Delete")

	beego.Router("/hook/server_list/:id", &controllers.RelationController{},"*:Index")
	beego.Router("/hook/server_list/add", &controllers.RelationController{},"*:AddServer")
	beego.Router("/hook/server_list/delete/:id", &controllers.RelationController{},"*:DeleteServer")

	beego.Router("/hook/scheduler/:id", &controllers.SchedulerController{},"*:Index")
	beego.Router("/hook/scheduler/console/:scheduler_id", &controllers.SchedulerController{},"*:Console")
	beego.Router("/hook/scheduler/resume/:scheduler_id", &controllers.SchedulerController{},"*:Resume")
	beego.Router("/hook/scheduler/cancel/:scheduler_id", &controllers.SchedulerController{},"*:Cancel")
	beego.Router("/hook/scheduler/status/:scheduler_id", &controllers.SchedulerController{},"*:Status")

	//用于GitHub Gogs Gitlab 等通知使用
	beego.Router("/payload/:key",&controllers.PayloadController{},"*:Index")

	beego.Router("/server", &controllers.ServerController{},"*:Index")
	beego.Router("/server/edit/?:id", &controllers.ServerController{},"*:Edit")
	beego.Router("/server/delete",&controllers.ServerController{},"post:Delete")

	beego.Router("/member",&controllers.MemberController{},"*:Index")
	beego.Router("/member/edit/?:id", &controllers.MemberController{},"*:Edit")
	beego.Router("/member/delete/?:id", &controllers.MemberController{},"*:Delete")
	beego.Router("/member/upload", &controllers.MemberController{},"post:Upload")

	beego.Router("/my", &controllers.MemberController{},"*:My")

	beego.Router("/login", &controllers.AccountController{},"*:Login");
	beego.Router("/logout", &controllers.AccountController{},"*:Logout");
}
