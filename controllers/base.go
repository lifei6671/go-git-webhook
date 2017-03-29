package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lifei6671/go-git-webhook/models"
	"github.com/lifei6671/go-git-webhook/conf"
	"bytes"
)

type BaseController struct {
	beego.Controller
	Member *models.Member
	Scheme string
}

func (c *BaseController) Prepare (){
	c.Data["SiteName"] = "Git WebHook"
	c.Data["Member"] = models.Member{}

	if member,ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0{
		c.Member = &member
		c.Data["Member"] = c.Member
	}else{
		//member := models.NewMember()
		//member.MemberId = 1
		//member.Find()
		//c.Member = member
		//c.Data["Member"] = *c.Member
	}
	scheme := "http"

	if c.Ctx.Request.TLS != nil {
		scheme = "https"
	}
	c.Scheme = scheme
}

//获取或设置当前登录用户信息,如果 MemberId 小于 0 则标识删除 Session
func (c *BaseController) SetMember(member models.Member) {

	if member.MemberId <= 0 {
		c.DelSession(conf.LoginSessionName)
		c.DelSession("uid")
		c.DestroySession()
	} else {
		c.SetSession(conf.LoginSessionName, member)
		c.SetSession("uid", member.MemberId)
	}
}

//响应 json 结果
func (c *BaseController) JsonResult(errCode int,errMsg string,data ...interface{}){
	json := make(map[string]interface{},3)

	json["errcode"] = errCode
	json["message"] = errMsg

	if len(data) > 0 && data[0] != nil{
		json["data"] = data[0]
	}

	c.Data["json"] = json
	c.ServeJSON(true)
	c.StopRun()
}

func (c *BaseController) UrlFor (endpoint string, values ...interface{}) string {

	return c.BaseUrl() + beego.URLFor(endpoint,values...)
}

func (c *BaseController) BaseUrl() string {
	scheme := "http://"

	if c.Ctx.Request.TLS != nil {
		scheme = "https://"
	}
	return scheme + c.Ctx.Request.Host
}

func (c *BaseController) NotFound(message interface{})  {
	c.TplName = "errors/404.html"
	c.Layout = ""
	c.Data["Model"] = map[string]interface{}{"Message":message}

	html,_ := c.RenderString()

	c.Abort(html)
}

func (c *BaseController) Forbidden(message interface{}) {
	c.TplName = "errors/403.html"
	c.Layout = ""
	c.Data["Model"] = map[string]interface{}{"Message":message}

	html,_ := c.RenderString()

	c.Abort(html)
}

func (c *BaseController) ServerError (message interface{}) {
	c.TplName = "errors/500.html"
	c.Layout = ""
	c.Data["Model"] = map[string]interface{}{"Message":message}


	html,_ := c.RenderString()

	c.Abort(html)
}

func (c *BaseController) ExecuteViewPathTemplate(tplName string,data interface{}) (string,error){
	var buf bytes.Buffer

	viewPath := c.ViewPath

	if c.ViewPath == "" {
		viewPath = beego.BConfig.WebConfig.ViewsPath

	}

	if err := beego.ExecuteViewPathTemplate(&buf,tplName,viewPath,data); err != nil {
		return "",err
	}
	return buf.String(),nil
}