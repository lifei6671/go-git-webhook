package controllers

import (
	"strings"
	"go-git-webhook/models"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Index() {
	c.Prepare()

	c.Layout = ""
	c.TplName = "home/index.html"


}

func (c *HomeController) Edit() {
	c.Prepare()
	c.Layout = ""
	c.TplName = "home/edit.html"

	if c.Ctx.Input.IsPost() {
		id,_ 	:= c.GetInt("id",0)
		name 	:= strings.TrimSpace(c.GetString("name",""))
		branch 	:= strings.TrimSpace(c.GetString("branch",""))
		tag 	:= strings.TrimSpace(c.GetString("tag",""))
		shell 	:= strings.TrimSpace(c.GetString("shell",""))

		if name == "" {
			c.JsonResult(500,"Repository Name is require.")
		}
		if branch == "" {
			branch = "master"
		}
		if tag == "" {
			c.JsonResult(500,"Server Tag is require.")
		}
		if shell == "" {
			c.JsonResult(500,"Callback Shell Script is require.")
		}

		webHook := models.NewWebHook()

		if id > 0 {
			if err := webHook.Find(id);err != nil {
				c.JsonResult(500,err.Error())
			}
			if webHook.CreateAt != c.Member.MemberId {
				c.JsonResult(403,"Permission denied")
			}
		}

		webHook.RepositoryName = name
		webHook.BranchName = branch
		webHook.Tag = tag
		webHook.Shell = shell
	}
}

func (c *HomeController) Delete()  {

}