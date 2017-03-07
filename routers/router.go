package routers

import (
	"github.com/astaxie/beego"
	"go-git-webhook/controllers"
)

func init()  {
	beego.Router("/", &controllers.HomeController{},"*:Index")

	beego.Router("/hook/edit/?:id", &controllers.HomeController{},"*:Edit")
	beego.Router("/hook/delete", &controllers.HomeController{},"post:Delete")

	beego.Router("/server", &controllers.ServerController{},"*:Index")
	beego.Router("/server/edit/?:id", &controllers.ServerController{},"*:Edit")
	beego.Router("/server/delete",&controllers.ServerController{},"post:Delete")

	beego.Router("/member",&controllers.MemberController{},"*:Index")

	beego.Router("/my", &controllers.MemberController{},"*:My")

	beego.Router("/login", &controllers.AccountController{},"*:Login");
	beego.Router("/logout", &controllers.AccountController{},"*:Logout");
}
