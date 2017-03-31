package controllers

import (
	"strings"
	"github.com/lifei6671/go-git-webhook/models"
	"github.com/lifei6671/go-git-webhook/modules/pager"
	"strconv"
	"fmt"
)

//首页WebHook控制器
type HomeController struct {
	BaseController
}

//首页
func (c *HomeController) Index() {
	c.Prepare()

	c.Layout = ""
	c.TplName = "home/index.html"

	pageIndex, _ := c.GetInt("page", 1)

	var hooks []models.WebHook

	pageOptions := pager.PageOptions{
		TableName:  models.NewWebHook().TableName(),
		EnableFirstLastLink : true,
		CurrentPage : pageIndex,
		ParamName : "page",
		PageSize : 15,
		Conditions : " AND create_at = " + strconv.Itoa(c.Member.MemberId) + " order by web_hook_id desc",
	}



	//返回分页信息,
	//第一个:为返回的当前页面数据集合,ResultSet类型
	//第二个:生成的分页链接
	//第三个:返回总记录数
	//第四个:返回总页数
	totalItem, totalCount, rs, pageHtml := pager.GetPagerLinks(&pageOptions, c.Ctx)

	_,err := rs.QueryRows(&hooks)      //把当前页面的数据序列化进一个切片内

	if err != nil {
		fmt.Println(err)
	}

	c.Data["lists"] = hooks
	c.Data["html"] = pageHtml
	c.Data["totalItem"] = totalItem
	c.Data["totalCount"] = totalCount
	c.Data["BaseUrl"] = c.BaseUrl()
	c.Data["WebHook"] = true
}

//编辑
func (c *HomeController) Edit() {
	c.Prepare()
	c.Layout = ""
	c.TplName = "home/edit.html"

	if c.Ctx.Input.IsPost() {
		id, _ := c.GetInt("id", 0)
		name := strings.TrimSpace(c.GetString("name", ""))
		branch := strings.TrimSpace(c.GetString("branch", ""))
		tag := strings.TrimSpace(c.GetString("tag", ""))
		shell := strings.TrimSpace(c.GetString("shell", ""))
		status, _ := c.GetInt("status", 0)
		hookType := c.GetString("hook_type","")

		if name == "" {
			c.JsonResult(500, "Repository Name is require.")
		}
		if branch == "" {
			branch = "master"
		}

		if shell == "" {
			c.JsonResult(500, "Callback Shell Script is require.")
		}

		webHook := models.NewWebHook()

		if id > 0 {
			webHook.WebHookId = id
			if err := webHook.Find(); err != nil {
				c.JsonResult(500, err.Error())
			}
			if webHook.CreateAt != c.Member.MemberId {
				c.JsonResult(403, "Permission denied")
			}
		}

		webHook.RepositoryName = name
		webHook.BranchName = branch
		webHook.Tag = tag
		webHook.Shell = shell
		webHook.Status = status
		webHook.CreateAt = c.Member.MemberId
		webHook.HookType = hookType

		if err := webHook.Save(); err != nil {
			c.JsonResult(500, err.Error())
		}
		data := make(map[string]interface{},5)

		if id <= 0 {

			view ,_:= c.ExecuteViewPathTemplate("home/index_list.html",webHook)

			data["view"] = view
		}



		data["errcode"] = 0
		data["message"] = "ok"

		data["data"] = webHook

		c.Data["json"] = data
		c.ServeJSON(true)
		c.StopRun()

	}

	id,_ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if id <= 0 {
		c.Abort("404")
	}

	webHook := models.NewWebHook()
	webHook.WebHookId = id

	if err := webHook.Find();err != nil {
		c.TplName = "errors/500.html"
		c.Data["Message"] = err.Error()
	}else{
		c.Data["Model"] = webHook

		c.Data["HookUrl"] = c.BaseUrl() + c.URLFor("PayloadController.Index",":key",webHook.Key)
	}
	c.Data["WebHook"] = true
}

//删除
func (c *HomeController) Delete()  {
	webHookId,_ := c.GetInt("id",0)
	if webHookId <= 0 {
		c.JsonResult(500,"Server ID is require.")
	}

	webHook := models.NewWebHook()
	webHook.WebHookId = webHookId

	if err := webHook.Find();err != nil {
		c.JsonResult(500,"Git WebHook does not exist")
	}
	if webHook.CreateAt != c.Member.MemberId {
		c.JsonResult(403,"Permission denied")
	}

	if err := webHook.Delete(); err != nil {
		c.JsonResult(500,"failed to delete")
	}

	models.NewRelation().DeleteByWhere(" AND web_hook_id = ?",webHookId)

	models.NewScheduler().DeleteByWhere(" AND web_hook_id = ?",webHookId)

	c.JsonResult(0,"ok")
}

