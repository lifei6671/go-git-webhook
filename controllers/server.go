package controllers

import (
	"strings"
	"go-git-webhook/models"
	"go-git-webhook/modules/pager"
	"strconv"
	"fmt"
)

type ServerController struct {
	BaseController
}

func (c *ServerController) Index() {
	c.Prepare()

	c.Layout = ""
	c.TplName = "server/index.html"


	pageIndex, _ := c.GetInt("page", 1)

	var servers []models.Server

	pageOptions := pager.PageOptions{
		TableName:  models.NewServer().TableName(),
		EnableFirstLastLink : true,
		CurrentPage : pageIndex,
		ParamName : "page",
		Conditions : " AND create_at = " + strconv.Itoa(c.Member.MemberId) + " order by server_id desc",
	}



	//返回分页信息,
	//第一个:为返回的当前页面数据集合,ResultSet类型
	//第二个:生成的分页链接
	//第三个:返回总记录数
	//第四个:返回总页数
	totalItem, totalCount, rs, pageHtml := pager.GetPagerLinks(&pageOptions, c.Ctx)

	_,err := rs.QueryRows(&servers)      //把当前页面的数据序列化进一个切片内

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(servers)

	c.Data["lists"] = servers
	c.Data["html"] = pageHtml
	c.Data["totalItem"] = totalItem
	c.Data["totalCount"] = totalCount

}


func (c *ServerController) Edit()  {
	c.Prepare()
	c.Layout = ""
	c.TplName = "server/edit.html"

	if c.Ctx.Input.IsPost() {
		id,_ := c.GetInt("id",0)
		account := c.GetString("account", "")
		serverName := c.GetString("name", "");
		ipAddress := c.GetString("ip", "")
		port, err := c.GetInt("port", 22)
		serverType := c.GetString("type", "ssh")
		status,_ := c.GetInt("status",0)

		if status !=0 && status != 1 {
			status = 0
		}

		if err != nil {
			c.JsonResult(500, "端口号错误");
		}
		tag := c.GetString("tag", "")

		key := c.GetString("key", "")

		if serverName == "" {
			c.JsonResult(500, "Server Name is require.")
		}
		if ipAddress == "" {
			c.JsonResult(500, "Server Ip is require.")
		}
		if port <= 0 {
			c.JsonResult(500, "Port is require.")
		}
		if tag != "" {

		}
		if key == "" {
			c.JsonResult(500, "SSH Private Key or Account Password is require.")
		}

		if !strings.EqualFold(serverType, "ssh") && !strings.EqualFold(serverType, "client") {
			c.JsonResult(500, "Server Type error.")
		}
		server := models.NewServer()

		if id > 0{
			if err := server.Find(id);err != nil {
				c.JsonResult(500,err.Error())
			}
			//如果不是本人创建则返回403
			if server.CreateAt != c.Member.MemberId {
				c.Abort("403")
			}
		}

		server.Account = account
		server.CreateAt = c.Member.MemberId
		server.IpAddress = ipAddress
		server.Name = serverName
		server.Port = port
		server.Tag = tag
		server.PrivateKey = key
		server.Type = serverType
		server.Status = status

		if err := server.Save(); err != nil {
			c.JsonResult(500, "Save failed:" + err.Error())
		} else {
			data := make(map[string]interface{},5)

			if id <= 0 {
				c.TplName = "server/index_list.html"
				c.Data["ServerId"] = server.ServerId
				c.Data["Name"] = server.Name
				c.Data["IpAddress"] = server.IpAddress
				c.Data["Port"] = server.Port
				c.Data["Tag"] = server.Tag
				c.Data["Status"] = server.Status
				c.Data["Type"] = server.Type
				c.Data["CreateTime"] = server.CreateTime

				view, err := c.RenderString()

				if err != nil {
					fmt.Println(err)
				}
				data["view"] = view
			}



			data["errcode"] = 0
			data["message"] = "ok"

			data["data"] = server

			c.Data["json"] = data
			c.ServeJSON(true)
			c.StopRun()

		}
	}

	id,err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Abort("404")
	}
	server := models.NewServer()

	if err := server.Find(id); err != nil {
		c.Abort("404")
	}
	//如果不是本人创建则返回403
	if server.CreateAt != c.Member.MemberId {
		c.Abort("403")
	}
	if c.Ctx.Input.IsAjax() {

		c.JsonResult(0,"ok",*server)
	}
	c.Data["Server"] = server



}

//删除一个Server
func (c *ServerController) Delete() {
	id,_ := c.GetInt("id",0)

	if id <= 0 {
		c.JsonResult(500,"Server ID is require.")
	}
	server := models.NewServer()

	if err := server.Find(id);err != nil {
		c.JsonResult(500,err.Error())
	}
	if server.CreateAt != c.Member.MemberId {
		c.JsonResult(403,"Permission denied")
	}
	if err := server.Delete();err != nil {
		c.JsonResult(500,err.Error())
	}

	webHook := models.NewWebHook()

	webHook.DeleteForServerId(id)

	c.JsonResult(0,"ok")
}

