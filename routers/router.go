package routers

import (
	"github.com/astaxie/beego"
	"go-git-webhook/controllers"
)

func init()  {
	beego.Router("/login", &controllers.AccountController{},"*:Login");
	beego.Router("/logout", &controllers.AccountController{},"*:Logout");
}
