package controllers

import (
	"github.com/lifei6671/go-git-webhook/models"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/go-git-webhook/tasks"
	"github.com/lifei6671/go-git-webhook/modules/hooks"
	"strings"
	"errors"
)

// PayloadController Git回调页面控制器.
type PayloadController struct {
	BaseController
}
// Index 回调首页 .
func (c *PayloadController) Index() {
	c.Prepare()
	c.TplName = ""

	key := c.Ctx.Input.Param(":key")
	if key == ""{
		c.Abort("404")
	}

	webHook := models.NewWebHook()

	if err := webHook.FindByKey(key);err != nil {
		c.Ctx.WriteString("Git WebHook no found.")
		c.StopRun()
	}

	body := string(c.Ctx.Input.RequestBody[0:])


	if body == ""{
		c.Ctx.WriteString("Request body error.")
		c.StopRun()
	}

	hookData ,err := createWebHook(webHook.HookType,body)

	if err != nil {
		c.Ctx.WriteString(err.Error())
		logs.Error("",err.Error())
		c.StopRun()
	}

	branch,err := hookData.BranchName()

	if err != nil {
		c.Ctx.WriteString("Request is not valid Git webhook: " + err.Error())
		logs.Error("Get Branch Name => ",err)
		c.StopRun()
	}

	repo ,err := hookData.RepositoryName()

	if err != nil {
		c.Ctx.WriteString("Request is not valid Git webhook: " + err.Error())
		logs.Error("Get Repository Name => ",err.Error())
		c.StopRun()
	}

	if webHook.RepositoryName != repo || (webHook.BranchName != branch && "heads/"+ webHook.BranchName != branch){
		c.Ctx.WriteString( "Not match the Repo and Branch.")
		logs.Error("Not match the Repo and Branch.")
		c.StopRun()
	}

	scheduler := models.NewScheduler()

	scheduler.WebHookId = webHook.WebHookId
	scheduler.Data = body

	if push,err := hookData.UserName() ;err == nil {
		scheduler.PushUser = push
	}
	if value,err := hookData.AfterValue(); err == nil{
		scheduler.ShaValue = value
	}else{
		logs.Error("",err.Error())
		scheduler.ShaValue = "无"
	}

	scheduler.Status = "wait"

	relationDetailed,err := models.FindRelationDetailedByWhere("AND relation.web_hook_id = ?",webHook.WebHookId)

	if err != nil {
		logs.Error(5001,err.Error())

		c.Ctx.WriteString("Data error")
		c.StopRun()
	}
	if len(relationDetailed) < 0 {
		c.Ctx.WriteString("Data is empty")
		c.StopRun()
	}


	schedulerList := make([]models.Scheduler,len(relationDetailed))

	for i,relation := range relationDetailed {

		scheduler.SchedulerId = 0
		scheduler.ServerId = relation.ServerId
		scheduler.WebHookId = relation.WebHookId
		scheduler.RelationId = relation.RelationId

		schedulerList[i] = *scheduler
	}


	nums,err := models.NewScheduler().InsertMulti(schedulerList)

	if nums <=0 || err != nil {
		if err != nil {
			logs.Error(0,err.Error())

		}
		c.Ctx.WriteString("Data error")
	}
	for _, scheduler := range schedulerList {
		tasks.Add(tasks.Task{ SchedulerId : scheduler.SchedulerId ,ServerId:scheduler.ServerId,WebHookId:scheduler.WebHookId})
	}

	c.Ctx.WriteString("Work put into Queue.")
	c.StopRun()
}

func createWebHook(t string,data string) (hooks.WebHookRequestInterface,error) {
	if strings.EqualFold(t,"github") {
		return hooks.NewGitHubWebHook(data)
	}else if strings.EqualFold(t,"gitlab") {
		return hooks.NewGitLabWebHook(data)
	}else if strings.EqualFold(t,"gogs") {
		return hooks.NewGogsWebHook(data)
	}else if strings.EqualFold(t,"gitosc") {
		return hooks.NewGitOSCWebHook(data)
	}
	return nil,errors.New("The type does not support.")
}

