package controllers

import (
	"bytes"

	"github.com/lifei6671/go-git-webhook/models"
	"github.com/lifei6671/go-git-webhook/conf"

	"github.com/astaxie/beego"
)

// BaseController 基础控制器
type BaseController struct {
	beego.Controller
	Member *models.Member
	Scheme string
}

// Prepare 预处理.
func (c *BaseController) Prepare (){
	c.Data["SiteName"] = "Git WebHook"
	c.Data["Member"] = models.Member{}

	if member,ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0{
		c.Member = &member
		c.Data["Member"] = c.Member
	}

	scheme := "http"

	if c.Ctx.Request.TLS != nil {
		scheme = "https"
	}
	c.Scheme = scheme
}

// SetMember 获取或设置当前登录用户信息,如果 MemberId 小于 0 则标识删除 Session
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

// JsonResult 响应 json 结果
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

// UrlFor .
func (c *BaseController) UrlFor (endpoint string, values ...interface{}) string {

	return c.BaseUrl() + beego.URLFor(endpoint,values...)
}

// BaseUrl .
func (c *BaseController) BaseUrl() string {

	if baseUrl := beego.AppConfig.String("base_url"); baseUrl != "" {
		return baseUrl
	}

	return c.Ctx.Input.Scheme() + "://" + c.Ctx.Request.Host
}

// NotFound .
func (c *BaseController) NotFound(message interface{})  {
	c.TplName = "errors/404.html"
	c.Layout = ""
	c.Data["Model"] = map[string]interface{}{"Message":message}

	html,_ := c.RenderString()

	c.Abort(html)
}

// Forbidden .
func (c *BaseController) Forbidden(message interface{}) {
	c.TplName = "errors/403.html"
	c.Layout = ""
	c.Data["Model"] = map[string]interface{}{"Message":message}

	html,_ := c.RenderString()

	c.Abort(html)
}

// ServerError .
func (c *BaseController) ServerError (message interface{}) {
	c.TplName = "errors/500.html"
	c.Layout = ""
	c.Data["Model"] = map[string]interface{}{"Message":message}


	html,_ := c.RenderString()

	c.Abort(html)
}

// ExecuteViewPathTemplate 执行指定的模板并返回执行结果.
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