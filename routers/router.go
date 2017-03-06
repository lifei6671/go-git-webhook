package routers

import (
	"github.com/astaxie/beego"
	"go-git-webhook/controllers"
)

func init()  {
	beego.Router("/", &controllers.HomeController{},"*:Index")

	beego.Router("/server", &controllers.ServerController{},"*:Index")

	beego.Router("/member",&controllers.MemberController{},"*:Index")

	beego.Router("/my", &controllers.MemberController{},"*:My")

	beego.Router("/login", &controllers.AccountController{},"*:Login");
	beego.Router("/logout", &controllers.AccountController{},"*:Logout");
}
