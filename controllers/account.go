package controllers

import (
	"fmt"
	"time"
	"go-git-webhook/modules/passwords"
	"github.com/astaxie/beego"
	"go-git-webhook/modules/gob"
	"go-git-webhook/conf"
	"go-git-webhook/models"
)

type AccountController struct {
	BaseController
}

//用户登录
func (c *AccountController) Login()  {
	c.Prepare()

	var remember struct { MemberId int ; Account string; Time time.Time}

	//如果Cookie中存在登录信息
	if cookie,ok := c.GetSecureCookie(conf.GetAppKey(),"login");ok{

		if err := gob.Decode(cookie,&remember); err == nil {
			member := models.NewMember()
			if err := models.NewMember().Find(remember.MemberId); err == nil {
				c.SetMember(*member)

				c.Redirect(beego.URLFor("HomeController.Index"), 302)
				c.StopRun()
			}
		}
	}

	fmt.Println(passwords.PasswordHash("123456"))

	if c.Ctx.Input.IsPost() {
		account := c.GetString("account")
		password := c.GetString("passwd")

		member,err := models.NewMember().Login(account,password)

		//如果没有数据
		if err == nil {
			c.SetMember(*member)
			c.JsonResult(0,"ok")
			c.StopRun()
		}else{
			fmt.Println(err)
			c.JsonResult(500,"账号或密码错误",nil)
		}

		return
	}else{

		c.Layout = ""
		c.TplName = "account/login.html"
	}
}

//退出登录
func (c *AccountController) Logout(){
	c.SetMember(models.Member{});

	c.Redirect(beego.URLFor("AccountController.Login"),302)
}

func (c *AccountController) Lists (){
	c.Prepare()


}