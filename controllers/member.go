package controllers

import (
	"go-git-webhook/modules/pager"
	"go-git-webhook/modules/passwords"
	"github.com/astaxie/beego/logs"
	"go-git-webhook/models"
	"fmt"
)

type MemberController struct {
	BaseController
}

func (c *MemberController) Index() {
	c.Prepare()

	if c.Member.Role != 0 {
		c.Abort("403")
	}

	c.Layout = ""
	c.TplName = "member/list.html"
	c.Data["MemberSelected"] = true

	pageIndex, _ := c.GetInt("page", 1)

	var members []models.Member

	pageOptions := pager.PageOptions{
		TableName:  models.NewMember().TableName(),
		EnableFirstLastLink : true,
		CurrentPage : pageIndex,
		ParamName : "page",
		Conditions : " order by member_id desc",
	}



	//返回分页信息,
	//第一个:为返回的当前页面数据集合,ResultSet类型
	//第二个:生成的分页链接
	//第三个:返回总记录数
	//第四个:返回总页数
	totalItem, totalCount, rs, pageHtml := pager.GetPagerLinks(&pageOptions, c.Ctx)

	_,err := rs.QueryRows(&members)      //把当前页面的数据序列化进一个切片内

	if err != nil {
		logs.Error("",err.Error())
	}

	c.Data["lists"] = members
	c.Data["html"] = pageHtml
	c.Data["totalItem"] = totalItem
	c.Data["totalCount"] = totalCount
}


func (c *MemberController) My(){
	c.Layout = ""
	c.TplName = "member/my.html"
	c.Data["MemberSelected"] = true
}

func (c *MemberController) Edit() {
	c.TplName = "member/edit.html"

	c.Prepare()

	member_id ,_ := c.GetInt(":id")

	member := models.NewMember()

	if member_id > 0 {
		member.MemberId = member_id
		if err := member.Find(); err != nil {
			c.ServerError("Data query error:" + err.Error())
		}
	}

	if c.Ctx.Input.IsPost() {
		password := c.GetString("password")
		account := c.GetString("account")

		if member.MemberId > 0 {
			if password != "" {
				pass, _ := passwords.PasswordHash(password)
				member.Password = pass
			}
		} else {
			if account == ""{
				c.JsonResult(500,"Account is require.")
			}
			if password == "" {
				c.JsonResult(500,"Password is require.")
			}
			member.Role = 1
			member.Account = account
			member.Password = password
		}

		member.Email = c.GetString("email")
		member.Phone = c.GetString("phone")
		member.Avatar = c.GetString("avatar")

		if member.Avatar == "" {
			member.Avatar = "/static/images/headimgurl.jpg"
		}

		var result error

		if member.MemberId > 0 {

			result = member.Update()
		} else {
			result = member.Add()
		}

		if result != nil {
			c.JsonResult(500, result.Error())
		}
		fmt.Printf("%+v",*member)
		view, err := c.ExecuteViewPathTemplate("member/list_item.html", *member)
		if err != nil {
			logs.Error("", err.Error())
		}

		data := map[string]interface{}{
			"view" : view,
		}
		c.JsonResult(0, "ok", data)

	}


	c.Data["Model"] = member
	c.Data["IsSelf"] = false
	if member.MemberId == c.Member.MemberId {
		c.Data["IsSelf"] = true
	}
}

func (c *MemberController) Delete() {
	c.Prepare()

	if c.Member.Role != 0 {
		c.Abort("403")
	}

	member_id ,err := c.GetInt(":id")

	if err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"Parameter error.")
	}

	member := models.NewMember()
	member.MemberId = member_id

	if err := member.Find();err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"Data query error.")
	}

	if member.Role == 0 {
		c.JsonResult(500,"不能删除管理员用户")
	}
	if err := member.Delete();err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"删除失败")
	}
	c.JsonResult(0,"ok")
}